package u4go

import (
	"errors"
	"fmt"
	uIoUtil "github.com/falconray0704/u4go/internal/ioutil"
	"go.uber.org/multierr"
	"io/ioutil"
	"os"
)

type TmpFile struct {
	GetFileFunc	func(dir, pattern string) (uIoUtil.IoFile, error)
	name string // name of the file as presented to Open.
}

// Name returns the name of the file as presented to Open.
func (this *TmpFile) Name() string {
	return this.name
}

// TempFile persists contents and returns the path and a clean func
func (this *TmpFile) NewFile(fileLocation, fileNamePrefix string, contents []byte) (path string, clean func(), err error) {

	var errOne, multiErr error
	var ioFile uIoUtil.IoFile
	//var filePath string

	if this == nil || this.GetFileFunc == nil {
		if ioFile, errOne = ioutil.TempFile(fileLocation, fileNamePrefix); errOne != nil {
			err = multierr.Append(multiErr, errOne)
			return "",nil, err
		}
	}

	defer func() {
		clean = func() {
			os.Remove(ioFile.Name())
		}
		if errOne = ioFile.Close(); errOne != nil {
			err = multierr.Append(multiErr, errOne)
		}
	}()

	if l, errOne := ioFile.Write(contents); errOne != nil {
		err = multierr.Append(multiErr, errOne)
		return ioFile.Name(), clean, err
	} else if l != len(contents) {
		err = multierr.Append(multiErr,
			errors.New(fmt.Sprintf("Write contents error, total len:%d, write len:%d", len(contents), l)))
		return ioFile.Name(), clean, err
	}

	return ioFile.Name(), clean, nil
}

func TempFile(fileLocation, fileNamePrefix string, contents []byte) (path string, clean func(), err error) {
	var tmpFile *TmpFile
	return tmpFile.NewFile(fileLocation, fileNamePrefix, contents)
}

