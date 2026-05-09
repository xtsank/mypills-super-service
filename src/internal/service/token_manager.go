package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/config"
)

type TokenManager interface {
	GenerateToken(userID uuid.UUID, isAdmin bool) (string, error)
	VerifyToken(tokenStr string) (uuid.UUID, bool, error)
}

type JWTManager struct {
	secretKey string
	duration  time.Duration
}

func NewJWTManager(i do.Injector) (TokenManager, error) {
	cfg := do.MustInvoke[*config.Config](i)

	secret, duration := cfg.GetJWTConfig()

	return &JWTManager{
		secretKey: secret,
		duration:  duration,
	}, nil
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID  uuid.UUID `json:"user_id"`
	IsAdmin bool      `json:"is_admin"`
}

func (m *JWTManager) GenerateToken(userID uuid.UUID, isAdmin bool) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:  userID,
		IsAdmin: isAdmin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) VerifyToken(tokenStr string) (uuid.UUID, bool, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return uuid.Nil, false, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims.UserID, claims.IsAdmin, nil
	}

	return uuid.Nil, false, fmt.Errorf("invalid token")
}
