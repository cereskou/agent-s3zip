package svc

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//Upload -
func (s *Service) Upload(reader *io.PipeReader, bucket, key string) (string, error) {
	//
	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	resp, err := s.Uploader().Upload(input)
	if err != nil {
		return "", err
	}

	return resp.Location, nil
}
