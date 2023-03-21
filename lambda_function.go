package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bondhansarker/exif_metadata"
	"github.com/bondhansarker/exif_metadata/file_template"
)

func handler(ctx context.Context, s3Event events.S3Event) (*exif_metadata.StructuredFileMetadata, error) {
	record := s3Event.Records[0]
	s3Key := record.S3.Object.Key
	s3Bucket := record.S3.Bucket.Name
	region := record.AWSRegion
	tempFileName := "/tmp/temp"

	awsConfig := &file_template.Config{
		S3: &file_template.S3Config{
			BucketName: s3Bucket,
			Region:     region,
		},
	}

	awsClient := file_template.NewBucketImplementation(awsConfig)

	OneMBChunkRange := "bytes=0-1048576"
	fileObject, err := awsClient.DownloadObjectAsFile(ctx, s3Key, OneMBChunkRange, tempFileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	_, structuredMetadata, err := exif_metadata.FetchMetaData(fileObject)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return structuredMetadata, nil
}

func main() {
	// fmt.Println(handler(context.TODO(), TestData()))
	lambda.Start(handler)
}
