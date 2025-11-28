package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MediaStorage defines the capability required by handlers/services to store media assets.
type MediaStorage interface {
	Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
}

type minioStorage struct {
	client  *minio.Client
	bucket  string
	baseURL string
}

// NewMinioStorageFromEnv initializes a MinIO-backed MediaStorage using environment variables.
func NewMinioStorageFromEnv() (MediaStorage, error) {
	endpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET_NAME")
	baseURL := strings.TrimRight(os.Getenv("ASSET_BASE_URL"), "/")
	useSSL := strings.EqualFold(os.Getenv("MINIO_USE_SSL"), "true")

	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return nil, errors.New("minio storage is not configured (check MINIO_* envs)")
	}

	parsed, err := parseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	// Prefer scheme from URL unless explicitly overridden via MINIO_USE_SSL.
	if !strings.EqualFold(os.Getenv("MINIO_USE_SSL"), "") {
		parsed.secure = useSSL
	}

	client, err := minio.New(parsed.host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: parsed.secure,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &minioStorage{
		client:  client,
		bucket:  bucket,
		baseURL: buildBaseURL(baseURL, parsed, bucket),
	}, nil
}

func (s *minioStorage) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.baseURL, strings.TrimLeft(objectName, "/")), nil
}

type endpointInfo struct {
	host   string
	scheme string
	secure bool
}

func (e endpointInfo) baseURL() string {
	scheme := e.scheme
	if scheme == "" {
		scheme = "http"
		if e.secure {
			scheme = "https"
		}
	}
	return fmt.Sprintf("%s://%s", scheme, e.host)
}

func parseEndpoint(raw string) (endpointInfo, error) {
	raw = strings.TrimSpace(raw)
	info := endpointInfo{}

	if raw == "" {
		return info, errors.New("empty MINIO_ENDPOINT")
	}

	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		u, err := url.Parse(raw)
		if err != nil {
			return info, err
		}
		info.host = u.Host
		info.scheme = u.Scheme
		info.secure = u.Scheme == "https"
		return info, nil
	}

	info.host = raw
	info.secure = false
	return info, nil
}

func buildBaseURL(custom string, info endpointInfo, bucket string) string {
	if custom != "" {
		return strings.TrimRight(custom, "/")
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(info.baseURL(), "/"), strings.TrimLeft(bucket, "/"))
}
