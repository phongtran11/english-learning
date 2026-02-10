package domain

import "github.com/golang-jwt/jwt/v5"

type AuthClaims struct {
	UserId uint
	Email  string
	jwt.RegisteredClaims
}

type RegisterRequest struct {
	Email    string
	Password string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type LoginRequest struct {
	Email    string
	Password string
}

type RefreshTokenRequest struct {
	RefreshToken string
}
