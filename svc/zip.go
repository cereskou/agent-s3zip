package svc

import (
	"log"
	"os"

	"ditto.co.jp/agent-s3zip/cx"
)

//Zip -
func (s *Service) Zip(pkg *cx.Package, filename string) error {
	cxzip := cx.NewGzip(filename)
	for _, fi := range pkg.Data {
		log.Printf("zip %v", fi.Local)
		err := cxzip.AddFile(fi)
		if err != nil {
			return err
		}
		//free space
		err = os.Remove(fi.Local)
		if err != nil {
			log.Printf("Remove error. %v", err)
		}
	}
	cxzip.Close()

	return nil
}
