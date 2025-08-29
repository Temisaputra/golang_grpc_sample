package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"github.com/Temisaputra/warOnk/internal/infrastructure/config"
	"github.com/Temisaputra/warOnk/internal/repository"
	"github.com/Temisaputra/warOnk/pkg/helper"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("invalid token")
	ErrUnauthorized = errors.New("unauthorized")
	ErrMissingToken = errors.New("missing token")
)

type JwtService interface {
	GenerateToken(user *entity.Users) (string, error)
	ValidateCurrentUser(r *http.Request) (*entity.Users, error)
}

type jwtService struct {
	cfg      config.Config
	log      zap.Logger
	userRepo repository.UserRepository
}

// context key type biar aman tidak tabrakan
type UserContextKey string
type JWTContextKey string

const (
	UserContext UserContextKey = "USER_CONTEXT_KEY"
	JWTContext  JWTContextKey  = "JWT_CONTEXT_KEY"
)

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

// GenerateToken bikin token baru untuk user
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

// ValidateCurrentUser validasi JWT dari header Authorization
func (s *jwtService) ValidateCurrentUser(r *http.Request) (*entity.Users, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, helper.NewErrUnauthorized("missing token")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, helper.NewErrUnauthorized("invalid token format")
	}

	tokenString := parts[1]
	claims := &Claims{}

	// Parse token sesuai JWT v5
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	// Tangani error expired secara spesifik
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, helper.NewErrUnauthorized("token expired")
		}
		return nil, helper.NewErrUnauthorized("invalid token")
	}

	// Optional check (claims.ExpiresAt biasanya dicek oleh JWT lib)
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

// SetUserContext masukkan user & jwt token ke dalam context
func SetUserContext(r *http.Request, user *entity.Users) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, UserContext, user)
	ctx = context.WithValue(ctx, JWTContext, r.Header.Get("Authorization"))
	return r.WithContext(ctx)
}

func GetUserContext(ctx context.Context) *entity.Users {
	user, ok := ctx.Value(UserContext).(*entity.Users)
	if !ok {
		return nil
	}
	return user
}

func GetJWTContext(ctx context.Context) string {
	token, _ := ctx.Value(JWTContext).(string)
	return token
}
