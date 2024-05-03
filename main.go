package main

import (
	"context"
	"database/sql"
	"dbstore/dumper"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alexmullins/zip"
	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	tmpDir string = "tmp"
)

type Config struct {
	DatabaseHost          string
	DatabaseName          string
	DatabaseUser          string
	DatabasePassword      string
	BucketEndpoint        string
	BucketAccessKeyID     string
	BucketSecretAccessKey string
	BucketName            string
	Password              string
	Interval              time.Duration
	Attempts              int
}

func createDump(cfg Config) (string, string, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabaseName))
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	name := time.Now().Format(time.RFC3339) + cfg.DatabaseName
	dumper, err := dumper.Register(db, tmpDir, name)
	if err != nil {
		return "", "", err
	}
	defer dumper.Close()

	resultFilename, err := dumper.Dump()
	if err != nil {
		return "", "", err
	}

	return resultFilename, name, nil
}

func zipFile(path string, name string, password string) (string, error) {
	zipFilename := tmpDir + "/" + name + ".zip"
	zipFile, err := os.Create(zipFilename)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipw := zip.NewWriter(zipFile)
	defer zipw.Close()

	w, err := zipw.Encrypt(filepath.Base(path), password)
	if err != nil {
		return "", err
	}

	fileToZip, err := os.Open(path)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(w, fileToZip)
	if err != nil {
		return "", err
	}

	return zipFilename, nil
}

func uploadFile(ctx context.Context, minioClient *minio.Client, bucketName, name string, path string) error {
	contentType := "application/zip"

	_, err := minioClient.FPutObject(ctx, bucketName, name, path, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("error uploading dump: %v", err)
	}

	return nil
}

func main() {
	log.Println("Starting application...")

	cfg := Config{
		DatabaseHost:          os.Getenv("DATABASE_HOST"),
		DatabaseName:          os.Getenv("DATABASE_NAME"),
		DatabaseUser:          os.Getenv("DATABASE_USER"),
		DatabasePassword:      os.Getenv("DATABASE_PASSWORD"),
		BucketEndpoint:        os.Getenv("BUCKET_ENDPOINT"),
		BucketAccessKeyID:     os.Getenv("BUCKET_ACCESS_KEY_ID"),
		BucketSecretAccessKey: os.Getenv("BUCKET_SECRET_ACCESS_KEY"),
		BucketName:            os.Getenv("BUCKET_NAME"),
		Password:              os.Getenv("PASSWORD"),
	}

	interval, _ := strconv.Atoi(os.Getenv("INTERVAL"))
	attempts, _ := strconv.Atoi(os.Getenv("ATTEMPTS"))
	cfg.Interval = time.Duration(interval) * time.Second
	cfg.Attempts = attempts

	minioClient, err := minio.New(cfg.BucketEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.BucketAccessKeyID, cfg.BucketSecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("Error creating MinIO client: %v", err)
	}

	ctx := context.Background()

	for {
		var dump string
		var name string
		var err error

		if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
			log.Fatalf("Error creating tmp directory: %v", err)
		}

		for i := 0; i < cfg.Attempts; i++ {
			dump, name, err = createDump(cfg)

			if err == nil {
				break
			}

			log.Printf("Failed to create dump in attempt %d: %v\n", i+1, err)
			time.Sleep(time.Second * 5)
		}

		if err != nil {
			log.Println("Failed to create dump after multiple attempts. Exiting...")
			return
		}

		zippedDump, err := zipFile(dump, name, cfg.Password)
		if err != nil {
			log.Println("Failed to create dump after multiple attempts. Exiting...")
			return
		}

		err = uploadFile(ctx, minioClient, cfg.BucketName, name+".zip", zippedDump)
		if err != nil {
			log.Printf("Error uploading dump: %v\n", err)
			return
		}

		log.Printf("Completed run, %s", dump)
		os.RemoveAll(tmpDir)
		time.Sleep(cfg.Interval)
	}
}
