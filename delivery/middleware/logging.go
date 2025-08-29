package middleware

import (
	"context"
	"time"

	"github.com/Temisaputra/warOnk/internal/infrastructure/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GRPCLoggingInterceptor(logger *zap.Logger, cfg *config.Config) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		requestID := uuid.New().String()

		// Simpan requestID di context
		ctx = context.WithValue(ctx, "request_id", requestID)

		// Eksekusi handler
		resp, err := handler(ctx, req)

		latency := time.Since(start)
		st, _ := status.FromError(err)

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", info.FullMethod),
			zap.String("service", cfg.AppName),
			zap.String("env", cfg.Env),
			zap.Duration("latency", latency),
			zap.String("status", st.Code().String()),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		switch {
		case st.Code().String() == "Internal":
			logger.Error("gRPC server error", fields...)
		case st.Code().String() == "NotFound":
			logger.Warn("gRPC client error", fields...)
		default:
			logger.Info("gRPC request handled", fields...)
		}

		return resp, err
	}
}
