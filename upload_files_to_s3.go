package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	// Define the command-line flags
	awsProfile := flag.String("profile", "", "aws profile")
	awsRegion := flag.String("region", "", "aws region")
	bucketName := flag.String("bucket", "", "s3 bucket name")
	srcFile := flag.String("srcFile", "", "file to upload")
	destFile := flag.String("destFile", "", "s3 file location")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -profile=<aws profile> -region=<aws region> -bucket=<s3 bucket> -srcFile=<file name> -tgtDest=<s3 target destination>\n", path.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "Arguments:")
		flag.PrintDefaults()
	}

	// Parse the command-line arguments
	flag.Parse()

	// Check for missing mandatory arguments
	if *awsProfile == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *awsRegion == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *bucketName == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *srcFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("\nAWS Profile: %s\n", *awsProfile)
	log.Printf("AWS Region: %s\n", *awsRegion)
	log.Printf("Bucket Name: %s\n", *bucketName)
	log.Printf("Source File Name: %s\n", *srcFile)
	log.Printf("Target Destination: %s\n", *destFile)

	// Create an AWS session with the shared config profile
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(*awsRegion)},
		Profile:           *awsProfile,
	})
	if err != nil {
		fmt.Println("Failed to create AWS session:", err)
		return
	}

	// Create an S3 uploader instance
	uploader := s3manager.NewUploader(sess)

	// Specify the bucket and key (object key) for the upload
	bucket := *bucketName

	// Open the file to be uploaded
	file, err := os.Open(*srcFile)
	if err != nil {
		fmt.Println("Failed to open file: ", err)
		return
	}
	defer file.Close()

	// Upload the file to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(*destFile),
		Body:   file,
	})
	if err != nil {
		fmt.Println("Failed to upload file:", err)
		return
	}

	log.Println("File uploaded successfully!")
}
