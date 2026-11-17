package auth

import "errors"

var (
	ErrInvalidCredential = errors.New("email atau password salah")
	ErrEmailTaken        = errors.New("email sudah terdaftar")
	ErrUserNotFound      = errors.New("user tidak ditemukan")
	ErrInvalidToken      = errors.New("token tidak valid")
)
