package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"github.com/Temisaputra/warOnk/internal/infrastructure/config"
	"github.com/Temisaputra/warOnk/internal/repository"
	"github.com/Temisaputra/warOnk/pkg/helper"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("invalid token")
	ErrUnauthorized = errors.New("unauthorized")
	ErrMissingToken = errors.New("missing token")
)

// JwtService interface untuk semua kebutuhan JWT
type JwtService interface {
	GenerateToken(user *entity.Users) (string, error)
	ValidateCurrentUserGRPC(ctx context.Context) (*entity.Users, error)
}

type jwtService struct {
	cfg      config.Config
	log      zap.Logger
	userRepo repository.UserRepository
}

// context key type biar aman tidak tabrakan
type contextKey string

const (
	UserContextKey contextKey = "USER_CONTEXT_KEY"
	JWTContextKey  contextKey = "JWT_CONTEXT_KEY"
)

// Constructor
func NewJwtService(cfg config.Config, log zap.Logger, userRepo repository.UserRepository) JwtService {
	return &jwtService{
		cfg:      cfg,
		log:      log,
		userRepo: userRepo,
	}
}

// Claims custom untuk JWT
type Claims struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

//////////////////////////////////////////////////////
// GENERATE TOKEN
//////////////////////////////////////////////////////

func (s *jwtService) GenerateToken(user *entity.Users) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 1) // expired 1 jam
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "warOnk-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

//////////////////////////////////////////////////////
// VALIDASI JWT UNTUK GRPC
//////////////////////////////////////////////////////

func (s *jwtService) ValidateCurrentUserGRPC(ctx context.Context) (*entity.Users, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrMissingToken
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, ErrMissingToken
	}

	return s.validateToken(authHeader[0])
}

//////////////////////////////////////////////////////
// FUNGSI GENERIK VALIDASI TOKEN
//////////////////////////////////////////////////////

func (s *jwtService) validateToken(authHeader string) (*entity.Users, error) {
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, helper.NewErrUnauthorized("invalid token format")
	}

	tokenString := parts[1]
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, helper.NewErrUnauthorized("token expired")
		}
		return nil, helper.NewErrUnauthorized("invalid token")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, helper.NewErrUnauthorized("token expired")
	}

	// Buat object user dari claims
	user := &entity.Users{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}

	return user, nil
}

//////////////////////////////////////////////////////
// CONTEXT HANDLER GRPC
//////////////////////////////////////////////////////

// SetUserContextGRPC masukkan user ke dalam context gRPC
func SetUserContextGRPC(ctx context.Context, user *entity.Users, token string) context.Context {
	ctx = context.WithValue(ctx, UserContextKey, user)
	ctx = context.WithValue(ctx, JWTContextKey, token)
	return ctx
}

func GetUserContext(ctx context.Context) *entity.Users {
	user, ok := ctx.Value(UserContextKey).(*entity.Users)
	if !ok {
		return nil
	}
	return user
}

func GetJWTContext(ctx context.Context) string {
	token, _ := ctx.Value(JWTContextKey).(string)
	return token
}
