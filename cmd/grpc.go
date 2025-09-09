package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcDelivery "github.com/Temisaputra/warOnk/internal/delivery/grpc"
	"github.com/Temisaputra/warOnk/internal/delivery/middleware"
	repository "github.com/Temisaputra/warOnk/internal/infrastructure/db"
	"github.com/Temisaputra/warOnk/internal/usecase"
	"github.com/Temisaputra/warOnk/pb/authpb"
	"github.com/Temisaputra/warOnk/pb/userpb"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		deps := InitDependencies() // 🔑 ambil dari bootstrap.go
		defer func() {
			sqlDB, _ := deps.DB.DB()
			sqlDB.Close()
			deps.Logger.Sync()
		}()

		// Repositories
		userRepo := repository.NewUserRepo(deps.DB)
		authRepo := repository.NewAuthRepo(deps.DB)
		transactionRepo := repository.NewTransactionRepo(deps.DB)

		// Usecases
		userUC := usecase.NewUserUsecase(userRepo, transactionRepo)
		authUC := usecase.NewAuthUsecase(authRepo, userRepo, transactionRepo, deps.JwtService)

		// Delivery layer
		userService := grpcDelivery.NewUserServiceServer(*userUC) // alias delivery/grpc
		authService := grpcDelivery.NewAuthServiceServer(*authUC)

		// Start gRPC server
		addr := fmt.Sprintf(":%s", deps.Cfg.GRPCPort)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			deps.Logger.Fatal("failed to listen", zap.Error(err))
		}

		jwtMiddleware := middleware.NewAuthMiddleware(deps.JwtService)
		// Register gRPC services
		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				jwtMiddleware.UnaryInterceptor, // <-- Middleware JWT
				middleware.GRPCLoggingInterceptor(deps.Logger, deps.Cfg),
			),
		)

		userpb.RegisterUserServiceServer(s, userService)
		authpb.RegisterAuthServiceServer(s, authService)

		// Enable reflection
		reflection.Register(s)

		// Run server async
		go func() {
			deps.Logger.Info("Starting gRPC server", zap.String("port", deps.Cfg.GRPCPort))
			if err := s.Serve(lis); err != nil {
				deps.Logger.Fatal("failed to serve", zap.Error(err))
			}
		}()

		// Graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		deps.Logger.Info("Shutting down gRPC server...")
		s.GracefulStop()
		deps.Logger.Info("gRPC server exited properly")
	},
}
