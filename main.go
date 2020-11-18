package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"ditto.co.jp/agent-s3zip/cx"
	"ditto.co.jp/agent-s3zip/svc"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

//ResultResponse -
type ResultResponse struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Error  string `json:"error,omitempty"`
}

func readPackage(input string) (*cx.Package, error) {
	var pkg cx.Package
	err := cx.JSON.NewDecoder(strings.NewReader(input)).Decode(&pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

//Response -
func Response(msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:            msg,
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "text/html; charset=utf-8",
		},
	}
}

//ResponseJSON -
func ResponseJSON(msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:            msg,
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

//handler - use s3 event
func handler(ctx context.Context, req events.S3Event) error {
	bucket := req.Records[0].S3.Bucket.Name
	key := req.Records[0].S3.Object.Key

	log.Printf("Bucket: %v", bucket)
	log.Printf("Key   : %v\n", key)

	sv := svc.NewService()
	//read package file
	pkg, err := sv.ReadPackage(bucket, key)
	if err != nil {
		return err
	}
	log.Printf("Pkg: %v", pkg.Count)

	//create a io pipe
	rp, wp := io.Pipe()
	//gzip
	zip := cx.NewGPipe(wp)

	log.Printf("create zip pipe")

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wp.Close()

		//download to io pipe
		err = sv.Download(pkg, zip)
		if err != nil {
			log.Println("error:", err)
		}
		zip.Close()

		wg.Done()
	}()

	gz := strings.ReplaceAll(key, ".jsonx", ".gz")
	location, err := sv.Upload(rp, bucket, gz)
	if err != nil {
		return err
	}

	wg.Wait()
	rp.Close()
	log.Println("uploaded", location)
	//delete .json
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	_, err = sv.Client().DeleteObject(deleteInput)
	if err != nil {
		return err
	}
	log.Printf("delete %v,%v", bucket, key)

	return nil
}

func main() {
	lambda.Start(handler)
}
