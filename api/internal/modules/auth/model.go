package auth

type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
	Role         string
}
