package svc

import (
	"fmt"
	"io"
)

//ChunkWriter -
type ChunkWriter struct {
	W         io.WriterAt
	Start     int64
	Size      int64
	Cur       int64
	WithRange string
}

//Write -
func (c *ChunkWriter) Write(p []byte) (n int, err error) {
	if c.Cur >= c.Size && len(c.WithRange) == 0 {
		return 0, io.EOF
	}

	n, err = c.W.WriteAt(p, c.Start+c.Cur)
	c.Cur += int64(n)

	return
}

//ByteRange -
func (c *ChunkWriter) ByteRange() string {
	if len(c.WithRange) != 0 {
		return c.WithRange
	}

	return fmt.Sprintf("bytes=%d-%d", c.Start, c.Start+c.Size-1)
}
