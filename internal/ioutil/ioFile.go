package ioutil

import (
	"io"
)

type IoFile interface {
	io.Reader
	io.Writer
	io.Closer
	Name() string // Name returns the name of the file as presented to Open.
}


