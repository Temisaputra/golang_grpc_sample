package middleware

import (
	"context"
	"strings"

	"github.com/Temisaputra/warOnk/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthMiddleware struct {
	jwtSvc auth.JwtService
}

// Constructor middleware
func NewAuthMiddleware(jwtSvc auth.JwtService) *AuthMiddleware {
	return &AuthMiddleware{jwtSvc: jwtSvc}
}

// UnaryInterceptor untuk validasi JWT di setiap request gRPC
func (m *AuthMiddleware) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	// ✅ Skip validasi untuk public endpoint (misal Login/Register)
	if strings.Contains(info.FullMethod, "Login") || strings.Contains(info.FullMethod, "Register") {
		return handler(ctx, req)
	}

	// ✅ Validasi token menggunakan service
	user, err := m.jwtSvc.ValidateCurrentUserGRPC(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized: %v", err)
	}

	// ✅ Simpan user ke context supaya bisa dipakai di service
	ctx = auth.SetUserContextGRPC(ctx, user, "")

	// ✅ Lanjutkan ke handler berikutnya
	return handler(ctx, req)
}
