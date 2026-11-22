// Package storage menyediakan klien S3 (Garage) untuk penyimpanan objek.
package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	cfgpkg "github.com/nyxsr/meme-portal/api/internal/config"
)

func NewStorage(ctx context.Context, cfg cfgpkg.StorageConfig) (*s3.Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("storage: gagal memuat konfigurasi aws: %w", err)
	}

	scheme := "http"
	if cfg.UseSSL {
		scheme = "https"
	}

	return s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(scheme + "://" + cfg.Endpoint)
		o.UsePathStyle = true
		o.HTTPClient = &http.Client{}
	}), nil
}

// EnsureCORS memasang aturan CORS pada bucket agar browser boleh mengunggah
// langsung ke S3 (presigned PUT) dari origin web. Tanpa ini preflight OPTIONS
// ditolak Garage dan fetch gagal dengan pesan "CORS failed to fetch".
func EnsureCORS(ctx context.Context, client *s3.Client, bucket string, origins []string) error {
	if len(origins) == 0 {
		return nil
	}
	_, err := client.PutBucketCors(ctx, &s3.PutBucketCorsInput{
		Bucket: aws.String(bucket),
		CORSConfiguration: &types.CORSConfiguration{
			CORSRules: []types.CORSRule{{
				AllowedOrigins: origins,
				AllowedMethods: []string{
					http.MethodGet, http.MethodHead, http.MethodPut,
				},
				AllowedHeaders: []string{"*"},
				ExposeHeaders:  []string{"ETag"},
				MaxAgeSeconds:  aws.Int32(3000),
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("storage: set CORS bucket gagal: %w", err)
	}
	return nil
}
