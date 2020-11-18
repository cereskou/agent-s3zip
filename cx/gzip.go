package cx

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	pgzip "github.com/klauspost/pgzip"
)

//Gzip -
type Gzip struct {
	sync.Mutex
	fw        *os.File
	tw        *tar.Writer
	gw        *pgzip.Writer
	filename  string
	filecount int64
	filesize  int64
	opened    bool
}

//NewGzip -
func NewGzip(fname string) *Gzip {
	return &Gzip{
		opened:    false,
		filecount: 0,
		filesize:  0,
		filename:  fname,
	}
}

//New -
func (z *Gzip) New() (err error) {
	fname := z.filename
	log.Printf("zip name: %v", fname)

	z.fw, err = os.Create(fname)
	if err != nil {
		return err
	}

	//gzip
	z.gw = pgzip.NewWriter(z.fw)
	z.tw = tar.NewWriter(z.gw)

	z.opened = true
	z.filecount = 0
	z.filesize = 0

	return nil
}

//Close -
func (z *Gzip) Close() {
	z.Lock()
	defer z.Unlock()

	z.tw.Close() //tar
	z.gw.Close() //gzip
	z.fw.Close() //file

	z.opened = false
	z.filecount = 0
	z.filesize = 0
}

//Flush -
func (z *Gzip) Flush() error {
	if !z.opened {
		return nil
	}
	return z.tw.Flush()
}

//AddFile -
func (z *Gzip) AddFile(obj *File) error {
	z.Lock()
	defer z.Unlock()
	if !z.opened {
		err := z.New()
		if err != nil {
			return err
		}
	}

	var name string = obj.Relative
	var size int64 = obj.Size
	if obj.Num > 0 && obj.Length > 0 {
		size = obj.Length
		// name = fmt.Sprintf("%v.%v", name, obj.Num)
	}
	f, err := os.Open(obj.Local)
	if err != nil {
		return err
	}
	defer f.Close()
	name = filepath.ToSlash(name)
	header := &tar.Header{
		Name:    name,
		Size:    size,
		Mode:    0644,
		ModTime: obj.ModTime,
	}
	log.Println("***", header.Name)
	err = z.tw.WriteHeader(header)
	if err != nil {
		msg := fmt.Sprintf("Could not write header for file '%s', got error '%s'", header.Name, err.Error())
		return errors.New(msg)
	}
	var rd io.Reader
	rd = f

	_, err = io.Copy(z.tw, rd)
	if err != nil {
		msg := fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", header.Name, err.Error())
		return errors.New(msg)
	}
	//count up
	z.filecount++
	z.filesize += size

	return nil
}
