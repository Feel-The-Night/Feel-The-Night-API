package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

const (
	AccessTokenDuration  = 15 * time.Minute   //15 minutos
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 dias

	issuer          = "auth-api"
	audienceAccess  = "access"
	audienceRefresh = "refresh"
)

// common auth errors
var (
	ErrExpiredToken  = errors.New("Token expired")
	ErrInvalidToken  = errors.New("Token invalid")
	ErrMissingSecret = errors.New("Missing JWT_SECRET configuration")
)

func getSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, ErrMissingSecret
	}

	return []byte(secret), nil
}

func newTokenID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GenerateAccessToken(userID int64, email string, name string) (string, error) {
	secret, err := getSecret()
	if err != nil {
		return "", err
	}

	now := time.Now()

	claims := Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			// Emissor
			Issuer: issuer,
			// Tipo de Token:
			Audience: jwt.ClaimStrings{audienceAccess},
			// Para quem
			Subject: fmt.Sprintf("%d", userID),
			// Quando foi emitido
			IssuedAt: jwt.NewNumericDate(now),
			// Quando expira
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenDuration)),
			// Nao usar antes da data
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("erro ao assinar access token: %w", err)
	}

	return tokenString, nil

}

func GenerateRefreshToken(userID int64) (string, error) {
	secret, err := getSecret()
	if err != nil {
		return "", err
	}

	jti, err := newTokenID()
	if err != nil {
		return "", fmt.Errorf("erro ao gerar id do token: %w", err)
	}

	now := time.Now()

	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    issuer,
			Subject:   strconv.FormatInt(userID, 10),
			Audience:  jwt.ClaimStrings{audienceRefresh},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("erro ao assinar refresh token: %w", err)
	}

	return tokenString, nil

}

func parseToken(tokenString string, claims jwt.Claims, audience string) (*jwt.Token, error) {
	secret, err := getSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
		// só aceita HS256
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil

}

func validateAccessToken(tokenString string) (*Claims, error) {
	token, err := parseToken(tokenString, &Claims{}, audienceAccess)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil

}

func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := parseToken(tokenString, &RefreshClaims{}, audienceRefresh)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
