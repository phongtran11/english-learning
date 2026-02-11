package service

import "github.com/golang-jwt/jwt/v5"

// authClaims is used for JWT token generation and parsing.
// This is an implementation detail of the auth service, not a domain entity.
type authClaims struct {
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
