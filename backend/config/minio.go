package config

import (
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinio() {
	var err error

	MinioClient, err = minio.New(
		fmt.Sprintf("%s:%s", ENV.MinioHost, ENV.MinioPort),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				ENV.MinioUser,
				ENV.MinioPassword,
				"",
			),
		},
	)
	if err != nil {
		log.Fatal("Cannot connect to Minio:", err)
	}

	log.Println("Connected to Minio at", ENV.MinioHost, ":", ENV.MinioPort)
}
