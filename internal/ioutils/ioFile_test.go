package ioutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IoFile_Mock_funcs(t *testing.T)  {

	var (
		cnt int
		err error
		contents = []byte("hello")
		l = len(contents)
	)
	cnt, err = ReadFuncMockErr(contents)
	assert.NotEqualf(t, l, cnt, "ReadFuncMockErr() expect cnt:%d != len(contents):%d", cnt, l)
	assert.NotNil(t, err, "ReadFuncMockErr() expect non-nil err")

	cnt, err = ReadFuncMockLenNotEnough(contents)
	assert.Equalf(t, l - 1, cnt, "ReadFuncMockLenNotEnough() expect cnt:%d == len(contents) - 1:%d", cnt, l)
	assert.Nil(t, err, "ReadFuncMockErr() expect nil err")

	cnt, err = ReadFuncMockSuccess(contents)
	assert.Equalf(t, l, cnt, "ReadFuncMockSuccess() expect cnt:%d == len(contents) :%d", cnt, l)
	assert.Nil(t, err, "ReadFuncMockSuccess() expect nil err")

	cnt, err = WriteFuncMockErr(contents)
	assert.NotEqualf(t, l, cnt, "WriteFuncMockErr() expect cnt:%d != len(contents):%d", cnt, l)
	assert.NotNil(t, err, "WriteFuncMockErr() expect non-nil err")

	cnt, err = WriteFuncMockLenNotEnough(contents)
	assert.Equalf(t, l - 1, cnt, "WriteFuncMockLenNotEnough() expect cnt:%d == len(contents) - 1:%d", cnt, l)
	assert.Nil(t, err, "WriteFuncMockErr() expect nil err")

	cnt, err = WriteFuncMockSuccess(contents)
	assert.Equalf(t, l, cnt, "WriteFuncMockSuccess() expect cnt:%d == len(contents) :%d", cnt, l)
	assert.Nil(t, err, "WriteFuncMockSuccess() expect nil err")

	err = CloseFuncMockErr()
	assert.NotNil(t, err, "CloseFuncMockErr() expect non-nil err")
	err = CloseFuncMockSuccess()
	assert.Nil(t, err, "CloseFuncMockSuccess() expect nil err")

}

func Test_IoFileMock(t *testing.T) {

	var (
		fileName = "HelloWorld"
		cnt int
		err error
		contents = []byte("hello")
		ioFile *IoFileMock
	)


	ioFile = GetIoFileMockAllSuccess()
	ioFile.FileName = fileName
	assert.NotNil(t, ioFile, "GetIoFileMockAllSuccess() should always success")

	cnt, err = ioFile.Read(contents)
	assert.Equal(t, len(contents), cnt, "GetIoFileMockAllSuccess() should always Read() success.")
	assert.Nil(t, err, "GetIoFileMockAllSuccess() should always Read() success with nil err.")
	cnt, err = ioFile.Write(contents)
	assert.Equal(t, len(contents), cnt, "GetIoFileMockAllSuccess() should always Write() success.")
	assert.Nil(t, err, "GetIoFileMockAllSuccess() should always Write() success with nil err.")
	err = ioFile.Close()
	assert.Nil(t, err, "GetIoFileMockAllSuccess() should always Close() success with  nil err.")

	assert.Equalf(t, fileName, ioFile.Name(), "GetIoFileMockAllSuccess() should get Name():%s == expect:%s .",ioFile.Name(), fileName)
}





