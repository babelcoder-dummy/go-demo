package port

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, password string) error
}
