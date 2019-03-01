package ioutil

import (
	"errors"
	"io"
)

type IoFile interface {
	io.Reader
	io.Writer
	io.Closer
	Name() string // Name returns the name of the file as presented to Open.
}

type ReadIoFile func (p []byte) (n int, err error)
type WriteIoFile func (p []byte) (n int, err error)
type CloseIoFile func () error

func ReadFuncMockErr(p []byte) (n int, err error) {
	return 0, errors.New("mocking io.Reader err")
}

func ReadFuncMockLenNotEnough(p []byte) (n int, err error) {
	return len(p) - 1, nil
}

func ReadFuncMockSuccess(p []byte) (n int, err error) {
	return len(p), nil
}

func WriteFuncMockErr(p []byte) (n int, err error) {
	return 0, errors.New("mocking io.Writer err")
}

func WriteFuncMockLenNotEnough(p []byte) (n int, err error) {
	return len(p) - 1, nil
}

func WriteFuncMockSuccess(p []byte) (n int, err error) {
	return len(p), nil
}

func CloseFuncMockErr() error {
	return errors.New("mocking io.Closer err")
}

func CloseFuncMockSuccess() error {
	return nil
}

type IoFileMock struct {
	ReadFunc ReadIoFile
	WriteFunc WriteIoFile
	CloseFunc CloseIoFile
	FileName string //the name of the file as presented to Open.
}

func (file * IoFileMock) Read(p []byte) (n int, err error) {
	return file.ReadFunc(p)
}

func (file * IoFileMock) Write(p []byte) (n int, err error) {
	return file.WriteFunc(p)
}

func (file * IoFileMock) Close() error {
	return file.CloseFunc()
}

func (file *IoFileMock) Name() string {
	return file.FileName
}

func GetIoFileMockAllSuccess() *IoFileMock {
	return &IoFileMock{
		ReadFunc: ReadFuncMockSuccess,
		WriteFunc:WriteFuncMockSuccess,
		CloseFunc:CloseFuncMockSuccess,
	}
}




