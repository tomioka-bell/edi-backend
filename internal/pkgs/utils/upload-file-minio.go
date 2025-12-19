package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

func UploadSignedUrl(
	documentType string,
	fileName string,
	file io.Reader,
) (string, error) {
	const useSSL = true

	minioClient, err := minio.New(viper.GetString("minio.endpoint"), &minio.Options{
		Creds: credentials.NewStaticV4(
			viper.GetString("minio.access_key"),
			viper.GetString("minio.secret_key"),
			"",
		),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println("MinIO client error:", err)
		return "", fmt.Errorf("failed to create MinIO client: %w", err)
	}

	uniqueFileName := fileName

	currentDate := time.Now()
	datePath := fmt.Sprintf("%d/%02d/%02d", currentDate.Year(), currentDate.Month(), currentDate.Day())

	filePath := fmt.Sprintf("%s/%s/%s", documentType, datePath, uniqueFileName)

	expires := time.Minute * 10
	presignedURL, err := minioClient.PresignedPutObject(
		context.Background(),
		viper.GetString("minio.bucket_name"),
		filePath,
		expires,
	)
	if err != nil {
		fmt.Println("Presigned URL error:", err)
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	buf := new(bytes.Buffer)
	header := make([]byte, 512)
	n, err := io.ReadFull(file, header)
	if err != nil && err != io.ErrUnexpectedEOF {
		fmt.Println("Failed to read file header:", err)
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	contentType := http.DetectContentType(header[:n])
	fmt.Println("Detected Content-Type:", contentType)

	_, err = buf.Write(header[:n])
	if err != nil {
		return "", fmt.Errorf("failed to write header data to buffer: %w", err)
	}
	_, err = io.Copy(buf, file)
	if err != nil {
		fmt.Println("Failed to copy file data to buffer:", err)
		return "", fmt.Errorf("failed to copy file data to buffer: %w", err)
	}

	req, err := http.NewRequest("PUT", presignedURL.String(), buf)
	if err != nil {
		fmt.Println("HTTP request creation error:", err)
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("File upload error:", err)
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		fmt.Println("File upload failed, status:", res.Status, "body:", string(body))
		return "", fmt.Errorf("failed to upload file, status: %s", res.Status)
	}

	fmt.Println("File uploaded successfully, status:", res.Status)

	return filePath, nil
}
