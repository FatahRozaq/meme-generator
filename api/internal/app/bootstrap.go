package app

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nyxsr/meme-portal/api/internal/config"
	"github.com/nyxsr/meme-portal/api/internal/platform/cache"
	"github.com/nyxsr/meme-portal/api/internal/platform/database"
	"github.com/nyxsr/meme-portal/api/internal/platform/storage"
	"github.com/redis/go-redis/v9"
)

type dependencies struct {
	db    *pgxpool.Pool
	redis *redis.Client
	s3    *s3.Client
}

func bootstrap(ctx context.Context, cfg config.Config) (deps *dependencies, err error) {
	var cleanups []func()
	defer func() {
		if err != nil {
			for _, cleanup := range slices.Backward(cleanups) {
				cleanup()
			}
		}
	}()

	db, err := database.NewPostgres(ctx, cfg.DB)
	if err != nil {
		return nil, err
	}
	cleanups = append(cleanups, db.Close)

	redisClient, err := cache.NewRedis(ctx, cfg.Redis)
	if err != nil {
		return nil, err
	}
	cleanups = append(cleanups, func() { _ = redisClient.Close() })

	s3Client, err := storage.NewStorage(ctx, cfg.Storage)
	if err != nil {
		return nil, err
	}

	// Browser mengunggah langsung ke S3 memakai presigned PUT, jadi bucket
	// harus mengizinkan origin web. Kegagalan di sini tidak fatal: S3 lain
	// mungkin mengatur CORS di luar aplikasi.
	if err := storage.EnsureCORS(ctx, s3Client, cfg.Storage.Bucket, withLoopbackOrigins(splitOrigins(cfg.AllowedOrigins))); err != nil {
		slog.Warn("bootstrap: gagal memasang CORS bucket", "err", err)
	}

	return &dependencies{db: db, redis: redisClient, s3: s3Client}, nil
}

// splitOrigins memecah daftar origin (dipisah koma) dari konfigurasi.
func splitOrigins(raw string) []string {
	var out []string
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			out = append(out, o)
		}
	}
	return out
}

// withLoopbackOrigins menambahkan origin loopback standar agar unggah
// langsung dari browser tetap jalan walau pengguna membuka app lewat
// 127.0.0.1 alih-alih localhost.
func withLoopbackOrigins(origins []string) []string {
	seen := make(map[string]bool, len(origins))
	for _, o := range origins {
		seen[o] = true
	}
	for _, o := range []string{"http://localhost:3000", "http://127.0.0.1:3000"} {
		if !seen[o] {
			origins = append(origins, o)
			seen[o] = true
		}
	}
	return origins
}
