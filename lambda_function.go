package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bondhansarker/exif_metadata"
	"github.com/bondhansarker/exif_metadata/file_template"
)

type MediaObject struct {
	Type       string                   `json:"type"`
	Extension  string                   `json:"extension"`
	Size       string                   `json:"size"`
	Resolution exif_metadata.Resolution `json:"resolution"`
	Location   exif_metadata.Location   `json:"location"`
	Timestamp  int64                    `json:"timestamp"`
	Errors     []string                 `json:"errors"`
}

func handler(ctx context.Context, s3Event events.S3Event) (*MediaObject, error) {
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
	metaDataObject, err := exif_metadata.FetchMetaData(fileObject)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	mediaObject := &MediaObject{
		Type:      metaDataObject.ContentInfo.Type,
		Size:      metaDataObject.ContentInfo.Size,
		Extension: metaDataObject.ContentInfo.Extension,
		Errors:    make([]string, 0),
	}

	resolution, err := metaDataObject.Resolution()
	if err != nil {
		fmt.Println(err)
		mediaObject.Errors = append(mediaObject.Errors, err.Error())
	} else {
		mediaObject.Resolution = *resolution
	}

	location, err := metaDataObject.Location()
	if err != nil {
		fmt.Println(err)
		mediaObject.Errors = append(mediaObject.Errors, err.Error())
	} else {
		mediaObject.Location = *location
	}

	timestamp, err := metaDataObject.DateTime()
	if err != nil {
		fmt.Println(err)
		mediaObject.Errors = append(mediaObject.Errors, err.Error())
	} else {
		mediaObject.Timestamp = timestamp.Timestamp.Unix()
	}
	return mediaObject, nil
}

func main() {
	// fmt.Println(handler(context.TODO(), TestData()))
	lambda.Start(handler)
}
