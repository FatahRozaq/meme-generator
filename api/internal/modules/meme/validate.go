package meme

import (
	"strings"
)

// validateBaseImage memvalidasi ekstensi dan ukuran gambar dasar, lalu
// mengembalikan Content-Type kanonik yang ikut ditandatangani di presigned
// URL.
func validateBaseImage(ext string, size int64) (string, error) {
	contentType, ok := allowedType(ext)
	if !ok {
		return "", ErrInvalidFile
	}
	if size <= 0 || size > maxBaseImageBytes {
		return "", ErrFileTooLarge
	}
	return contentType, nil
}

// allowedType memetakan ekstensi yang diizinkan ke Content-Type kanonik.
// Ekstensi di luar daftar ditolak, jadi bucket tidak bisa dipakai menyimpan
// file .exe/.zip/.html atau tipe lain yang bukan gambar.
func allowedType(ext string) (string, bool) {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg", true
	case ".png":
		return "image/png", true
	case ".gif":
		return "image/gif", true
	case ".webp":
		return "image/webp", true
	}
	return "", false
}
