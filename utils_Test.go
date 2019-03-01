package u4go

import (
	"errors"
	"fmt"
	uIoUtils "github.com/falconray0704/u4go/internal/ioutils"
	"go.uber.org/multierr"
	"io/ioutil"
	"os"
)

type GetIoFileFunc func(dir, pattern string) (uIoUtils.IoFile, error)

// TempFile persists contents and returns the path and a clean func
func newTempFile(getIoFileFunc GetIoFileFunc, fileLocation, fileNamePrefix string, contents []byte) (path string, clean func(), err error) {

	var errOne, multiErr error
	var ioFile uIoUtils.IoFile

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

func getIoutilTempfile(fileLocation, fileNamePrefix string) (uIoUtils.IoFile, error) {
	return ioutil.TempFile(fileLocation, fileNamePrefix)
}

func TempFile(fileLocation, fileNamePrefix string, contents []byte) (path string, clean func(), err error) {
	return newTempFile(getIoutilTempfile, fileLocation, fileNamePrefix, contents)
}

