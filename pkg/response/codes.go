package response

// Response Codes
const (
	CodeSuccess             = "SUCCESS"
	CodeCreated             = "CREATED"
	CodeBadRequest          = "BAD_REQUEST"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeNotFound            = "NOT_FOUND"
	CodeServerInternalError = "SERVER_INTERNAL_ERROR"
)

// Response Messages
const (
	MsgSuccess             = "Success"
	MsgUserCreated         = "User created"
	MsgUserUpdated         = "User updated"
	MsgUserDeleted         = "User deleted"
	MsgUserRegistered      = "User registered successfully"
	MsgInvalidID           = "Invalid ID"
	MsgUserNotFound        = "User not found"
	MsgLoginSuccess        = "Login success"
	MsgRefreshTokenSuccess = "Refresh token success"
)
