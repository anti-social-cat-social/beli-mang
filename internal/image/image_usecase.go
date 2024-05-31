package image

import (
	localError "belimang/pkg/error"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// createSession creates a new AWS session
func createSession() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	return session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})
}

// UploadFileToS3 uploads a file to S3 with a given prefix
func UploadFileToS3(fileHeader *multipart.FileHeader) (string, *localError.GlobalError) {
	file, errFile := fileHeader.Open()

	if errFile != nil {
		return "", localError.ErrInternalServer("error upload image", errFile)
	}
	defer file.Close()

	// Get file extension
	ext := filepath.Ext(fileHeader.Filename)

	sess, err := createSession()
	if err != nil {
		return "", localError.ErrInternalServer("error upload image", err)
	}

	svc := s3.New(sess)

	bucket := "projectsprint-bucket-public-read"
	key := uuid.NewString() + ext

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(fileHeader.Size),
	})

	if err != nil {
		return "", localError.ErrInternalServer("error upload image", err)
	}

	return fmt.Sprintf("https://awss3.%s", key), nil
}
