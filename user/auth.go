package user

type AuthService interface {
	AuthAccountFromToken(token string) (*Account, error)
}
