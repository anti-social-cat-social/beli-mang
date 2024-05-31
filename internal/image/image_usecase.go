package image

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"github.com/google/uuid"
	"mime/multipart"
	localError "belimang/pkg/error"
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

	sess, err := createSession()
	if err != nil {
	 	return "", localError.ErrInternalServer("error upload image", err)
	}
   
	svc := s3.New(sess)

	bucket := "aws"
	key := uuid.NewString()
   
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(fileHeader.Size),
	})
	if err != nil {
		return "", localError.ErrInternalServer("error upload image", err)
	}
   
	return fmt.Sprintf("https://awss3.%s.jpeg", key), nil
}