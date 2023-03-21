package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

func TestData() events.S3Event {
	jsonFile, err := os.Open(fmt.Sprintf("test.json"))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	var s3Event events.S3Event
	err = json.Unmarshal(byteValue, &s3Event)
	if err != nil {
		fmt.Println(err)
	}
	return s3Event
}
