package jwt

// UserClaims is the identity returned after validating a JWT (contract shared by HTTP middleware and adapters).
type UserClaims interface {
	GetUserID() string
	GetEmail() string
}
