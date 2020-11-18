package cx

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"sync"

	pgzip "github.com/klauspost/pgzip"
)

//GPipe -
type GPipe struct {
	sync.Mutex
	P *io.PipeWriter
	W *tar.Writer
	G *pgzip.Writer
}

//NewGPipe -
func NewGPipe(w *io.PipeWriter) *GPipe {
	z := &GPipe{
		P: w,
	}
	z.G = pgzip.NewWriter(w)
	z.W = tar.NewWriter(z.G)

	return z
}

//AddFile -
func (z *GPipe) AddFile(fi *File, in io.ReadCloser) error {
	z.Lock()
	defer z.Unlock()

	var name string = fi.Relative
	var size int64 = fi.Size
	if fi.Num > 0 && fi.Length > 0 {
		size = fi.Length
		name = fmt.Sprintf("%v.%v", name, fi.Num)
	}
	name = filepath.ToSlash(name)

	header := &tar.Header{
		Name:    name,
		Size:    size,
		Mode:    0644,
		ModTime: fi.ModTime,
	}
	err := z.W.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("Could not write header for file '%s', got error '%s'", header.Name, err.Error())
	}

	_, err = io.Copy(z.W, in)
	if err != nil {
		return fmt.Errorf("Could not copy the file '%s' data to the tarball, got error '%s'", header.Name, err.Error())
	}
	log.Printf("%v, size: %v", header.Name, header.Size)

	return nil
}

//Close -
func (z *GPipe) Close() error {
	z.Lock()
	defer z.Unlock()

	z.W.Close() //tar
	z.G.Close() //gzip
	z.P.Close() //pipe

	return nil
}
