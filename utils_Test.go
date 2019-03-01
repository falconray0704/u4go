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

type GetIoFileFunc func(dir, pattern string) (uIoUtil.IoFile, error)

// TempFile persists contents and returns the path and a clean func
func newTempFile(getIoFileFunc GetIoFileFunc, fileLocation, fileNamePrefix string, contents []byte) (path string, clean func(), err error) {

	var errOne, multiErr error
	var ioFile uIoUtil.IoFile

	if ioFile, errOne = getIoFileFunc(fileLocation, fileNamePrefix); errOne != nil {
		err = multierr.Append(multiErr, errOne)
		return "",nil, err
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

func getIoutilTempfile(fileLocation, fileNamePrefix string) (uIoUtil.IoFile, error) {
	return ioutil.TempFile(fileLocation, fileNamePrefix)
}

func TempFile(fileLocation, fileNamePrefix string, contents []byte) (path string, clean func(), err error) {
	return newTempFile(getIoutilTempfile, fileLocation, fileNamePrefix, contents)
}

