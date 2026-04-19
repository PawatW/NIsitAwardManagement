package auth

import (
	"fmt"
	"time"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID      `json:"user_id"`
	Email  string         `json:"email"`
	Role   enums.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// Legacy alias for backward compatibility
type Claims = JWTClaims

type RegisterTokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type TokenManager interface {
	GenerateAccessToken(userID uuid.UUID, email string, role enums.UserRole) (string, error)
	GenerateRegisterToken(email string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	ValidateRegisterToken(tokenString string) (*RegisterTokenClaims, error)
}

type jwtManager struct {
	secret      []byte
	expiryHours int
}

func NewJWTManager(cfg *configs.Config) TokenManager {
	return &jwtManager{
		secret:      []byte(cfg.JWT.Secret),
		expiryHours: cfg.JWT.ExpiryHours,
	}
}

// GenerateAccessToken creates a new JWT token
func (j *jwtManager) GenerateAccessToken(userID uuid.UUID, email string, role enums.UserRole) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.expiryHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken validates and parses a JWT token
func (j *jwtManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateRegisterToken creates a temporary JWT token for registration
func (j *jwtManager) GenerateRegisterToken(email string) (string, error) {
	claims := RegisterTokenClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)), // 15 minutes expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign register token: %w", err)
	}

	return signedToken, nil
}

// ValidateRegisterToken validates and parses a registration JWT token
func (j *jwtManager) ValidateRegisterToken(tokenString string) (*RegisterTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RegisterTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse register token: %w", err)
	}

	if claims, ok := token.Claims.(*RegisterTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid register token")
}
