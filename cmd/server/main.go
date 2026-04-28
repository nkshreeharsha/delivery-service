package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/nkshreeharsha/delivery-service/internal/handler"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	/* Instead of creating Access keys and maintaining it we can use IAM roles 
		for service accounts - meaning the pod/microservice is directly assigned an IAM role and 
		it can access the S3 bucket without needing to manage credentials. This is more secure and easier to maintain. 
	*/

	signer, err := newS3Signer(
		os.Getenv("Bucket_name"),
		os.Getenv("Region"),
		os.Getenv("Access_key_ID"),
		os.Getenv("Secret_access_key"),
	)
	if err != nil {
		log.Fatalf("Error creating S3 signer: %v", err)
	}

	h := handler.New(signer,"2026-04-22_1",60*time.Second)

	r := chi.NewRouter()
	r.Use(middleware.Logger) 

	r.Get("/v1/subscriber/{subID}/creativeList", h.GetCreativeList)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Swap for real GCS implementation
type S3Client struct{ 
	presignClient *s3.PresignClient
	bucket string
}

func (s *S3Client) SignURL(objectPath string, expiry time.Duration) (string, error) {
	ctx := context.Background()

	req, err := s.presignClient.PresignGetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(objectPath),
		},
		s3.WithPresignExpires(expiry),
	)
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func newS3Signer(bucket, region, accessKey, secretKey string) (*S3Client, error){
	context := context.Background()

	cfg, err := config.LoadDefaultConfig(context,config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")))

		if err != nil {
			return nil, err
		}

		s3Client := s3.NewFromConfig(cfg)
		presignClient := s3.NewPresignClient(s3Client)

		return  &S3Client{
			presignClient: presignClient,
			bucket: bucket,
		}, nil

}

