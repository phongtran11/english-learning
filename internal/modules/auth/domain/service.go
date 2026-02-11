package domain

// AuthService defines the business logic contract for authentication operations.
type AuthService interface {
	Register(req *RegisterRequest) error
	Login(req *LoginRequest, ip, userAgent string) (*TokenPair, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	Logout(refreshToken string) error
	LogoutAll(userID uint) error
}
