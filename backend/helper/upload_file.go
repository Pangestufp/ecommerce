package helper

import (
	"backend/config"
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func UploadFile(file multipart.File, filename string, client *minio.Client) (string, error) {

	_, err := client.PutObject(
		context.Background(),
		config.ENV.MinioBucket,
		filename,
		file,
		-1,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("/%s/%s", config.ENV.MinioBucket, filename)
	return url, nil
}

func GeneratePresignedURL(filePath string) (string, error) {

	client, err := minio.New(
		fmt.Sprintf("%s:%s", config.ENV.MinioHost, config.ENV.MinioPort),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				config.ENV.MinioUser,
				config.ENV.MinioPassword,
				"",
			),
		},
	)
	if err != nil {
		return "", err
	}

	trimmed := strings.TrimPrefix(filePath, "/")
	parts := strings.SplitN(trimmed, "/", 2)

	if len(parts) != 2 {
		return "", fmt.Errorf("invalid file path")
	}

	bucket := parts[0]
	object := parts[1]

	url, err := client.PresignedGetObject(
		context.Background(),
		bucket,
		object,
		time.Minute*5,
		nil,
	)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
