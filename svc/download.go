package svc

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"ditto.co.jp/agent-s3zip/cx"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cenkalti/backoff"
)

//DownloadPart -
func (s *Service) downloadPart(obj *cx.File, zip *cx.GPipe) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(obj.Bucket),
		Key:    aws.String(obj.Key),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", obj.Offset, obj.Offset+obj.Length-1)),
	}
	if obj.Length == 0 {
		input.Range = nil
	}

	output, err := s._client.GetObject(input)
	if err != nil {
		return err
	}
	defer output.Body.Close()

	//output
	err = zip.AddFile(obj, output.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

//Download -
func (s *Service) Download(pkg *cx.Package, zip *cx.GPipe) error {
	var zipCount int64
	var srcCount int64

	//file channel
	keysChan := make(chan *cx.File, 2048)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for fi := range keysChan {
				log.Printf("%v,%v", fi.Bucket, fi.Key)

				//リトライ
				bo := backoff.NewExponentialBackOff()
				bo.MaxInterval = 15 * time.Second
				bo.MaxElapsedTime = time.Minute

				operation := func() error {
					return s.downloadPart(fi, zip)
				}
				err := backoff.Retry(
					operation,
					backoff.WithMaxRetries(bo, 3))
				if err != nil {
					log.Println("download failed.", fi.Bucket, fi.Key)
				} else {
					atomic.AddInt64(&zipCount, 1)
				}
				//ok
			}
		}()
	}

	log.Printf("download...")
	offset := len(pkg.Dir)
	//scan file
	for _, fi := range pkg.Data {
		relative := fi.Local[offset:]
		if relative[0] != '/' {
			relative = "/" + relative
		}
		fi.Relative = relative

		keysChan <- fi
		srcCount++
	}

	close(keysChan)
	wg.Wait()
	log.Printf("downlaod end. File: %v, Zipped: %v", srcCount, zipCount)

	return nil
}
