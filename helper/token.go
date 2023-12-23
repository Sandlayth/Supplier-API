package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sandlayth/supplier-api/model"
)

// var secretKey []byte = []byte(os.Getenv("TOKEN_SECRET"))
var secretKey []byte = []byte("TOKEN_SECRET")

var accessTokenExpirationTime = time.Minute * 15
var refreshTokenExpirationTime = time.Hour * 12

func GenerateAccessToken(user *model.User) (string, error) {
	claims := model.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := generateToken(&claims, accessTokenExpirationTime)
	if err != nil {
		user.Claims = &claims
	}
	return token, err
}

func GenerateRefreshToken(user *model.User) (string, error) {
	claims := model.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := generateToken(&claims, refreshTokenExpirationTime)
	if err != nil {
		user.Claims = &claims
	}
	return token, err
}

// GenerateToken generates a JWT token for the given user with registered and custom claims.
func generateToken(claims *model.Claims, expirationTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses and validates a JWT token, returning the claims and user.
func VerifyToken(tokenString string) (*model.Claims, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, false, err
	}

	// Check if the token is valid
	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return nil, false, jwt.ErrTokenInvalidClaims
	}

	// Check if the token is expired
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, false, errors.New("token expired")
	}

	// Check if needs to be renewed
	remainingValidity := claims.ExpiresAt.Time.Sub(time.Now())
	threshold := (time.Duration(0.2 * float64(claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time))))
	needsRefresh := remainingValidity < threshold

	return claims, needsRefresh, nil
}
