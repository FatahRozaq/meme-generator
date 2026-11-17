package auth

// DTO permintaan/respons modul auth.
type RegisterDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshDTO struct {
	RefreshToken string `json:"refreshToken"`
}

type UserDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type AuthResponse struct {
	User  UserDTO      `json:"user"`
	Token TokenPayload `json:"token"`
}

type TokenPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}
