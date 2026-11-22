package meme

import "time"

const maxPresignTTL = 15 * time.Minute

// maxBaseImageBytes membatasi ukuran gambar dasar yang boleh diunggah.
const maxBaseImageBytes = 8 << 20 // 8 MiB

const (
	defaultLimit = 20
	maxLimit     = 100
)
