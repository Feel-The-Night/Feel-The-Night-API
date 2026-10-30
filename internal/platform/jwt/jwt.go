package jwt

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
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

const (
	AccessTokenDuration  = 15 * time.Minute   // 15 minutes
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 days

	issuer          = "auth-api"
	audienceAccess  = "access"
	audienceRefresh = "refresh"
)

// common auth errors
var (
	ErrExpiredToken  = errors.New("token expired")
	ErrInvalidToken  = errors.New("token invalid")
	ErrMissingSecret = errors.New("missing JWT_SECRET configuration")
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

func GenerateAccessToken(userID int64, email string, Nickname string, role string) (string, error) {
	secret, err := getSecret()
	if err != nil {
		return "", err
	}

	now := time.Now()

	claims := Claims{
		UserID: userID,
		Email:  email,
		Name:   Nickname,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// Issuer
			Issuer: issuer,
			// Token type
			Audience: jwt.ClaimStrings{audienceAccess},
			// Subject (who the token belongs to)
			Subject: fmt.Sprintf("%d", userID),
			// Issued at
			IssuedAt: jwt.NewNumericDate(now),
			// Expires at
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenDuration)),
			// Not valid before
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
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
		return "", fmt.Errorf("failed to generate token id: %w", err)
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
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
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
		// only accept HS256
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

func ValidateAccessToken(tokenString string) (*Claims, error) {
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
