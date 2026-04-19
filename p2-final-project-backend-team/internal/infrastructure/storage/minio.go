package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

// StorageClient interface defines methods for object storage operations
type StorageClient interface {
	UploadFile(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error
	DownloadFile(ctx context.Context, objectKey string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, objectKey string) error
	GetPresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error)
	GetPublicURL(objectKey string) string
	EnsureBucketExists(ctx context.Context) error
	GetBucketName() string
}

type minioClient struct {
	client         *minio.Client
	bucketName     string
	region         string
	publicEndpoint string
}

// NewMinioClient creates a new MinIO client
func NewMinioClient(config *configs.MinioConfig) (StorageClient, error) {
	// Initialize MinIO client
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	mc := &minioClient{
		client:         client,
		bucketName:     config.BucketName,
		region:         config.Region,
		publicEndpoint: resolvePublicEndpoint(config),
	}

	// Ensure bucket exists
	if err := mc.EnsureBucketExists(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	log.Info().
		Str("endpoint", config.Endpoint).
		Str("bucket", config.BucketName).
		Msg("MinIO client initialized successfully")

	return mc, nil
}

// UploadFile uploads a file to MinIO
func (m *minioClient) UploadFile(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, m.bucketName, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	log.Info().
		Str("bucket", m.bucketName).
		Str("object_key", objectKey).
		Int64("size", size).
		Msg("File uploaded successfully to MinIO")

	return nil
}

// DownloadFile downloads a file from MinIO
func (m *minioClient) DownloadFile(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	object, err := m.client.GetObject(ctx, m.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from MinIO: %w", err)
	}

	log.Info().
		Str("bucket", m.bucketName).
		Str("object_key", objectKey).
		Msg("File downloaded successfully from MinIO")

	return object, nil
}

// DeleteFile deletes a file from MinIO
func (m *minioClient) DeleteFile(ctx context.Context, objectKey string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	log.Info().
		Str("bucket", m.bucketName).
		Str("object_key", objectKey).
		Msg("File deleted successfully from MinIO")

	return nil
}

// GetPresignedURL generates a presigned URL for temporary access to a file
func (m *minioClient) GetPresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectKey, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	log.Debug().
		Str("bucket", m.bucketName).
		Str("object_key", objectKey).
		Dur("expiry", expiry).
		Msg("Presigned URL generated successfully")

	return presignedURL.String(), nil
}

func resolvePublicEndpoint(config *configs.MinioConfig) string {
	if config.PublicEndpoint != "" {
		return strings.TrimRight(config.PublicEndpoint, "/")
	}

	scheme := "http"
	if config.UseSSL {
		scheme = "https"
	}
	return scheme + "://" + strings.TrimRight(config.Endpoint, "/")
}

// GetPublicURL returns a direct public URL for objects in a public bucket.
func (m *minioClient) GetPublicURL(objectKey string) string {
	return fmt.Sprintf("%s/%s/%s", m.publicEndpoint, m.bucketName, strings.TrimLeft(objectKey, "/"))
}

// EnsureBucketExists checks if the bucket exists, if not, creates it
func (m *minioClient) EnsureBucketExists(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = m.client.MakeBucket(ctx, m.bucketName, minio.MakeBucketOptions{
			Region: m.region,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		log.Info().
			Str("bucket", m.bucketName).
			Str("region", m.region).
			Msg("Bucket created successfully")
	} else {
		log.Debug().
			Str("bucket", m.bucketName).
			Msg("Bucket already exists")
	}

	return nil
}

// GetBucketName returns the bucket name
func (m *minioClient) GetBucketName() string {
	return m.bucketName
}
