package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//Token create/verify

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type JWTManager struct {
	secret []byte
	ttl time.Duration
	issuer string
}

func NewJWTManager(secret, issuer string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		ttl: ttl,
		issuer: issuer,
	}
}

func (m *JWTManager) GenerateToken(userID uint64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: m.issuer,
			IssuedAt: jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (m *JWTManager) ParseAndValidate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, ErrExpiredToken
			}
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(m.issuer),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrExpiredToken
	}
	if claims.Issuer != m.issuer {
		return nil, ErrInvalidToken
	}
	if claims.UserID == 0 {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
