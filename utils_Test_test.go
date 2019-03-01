package u4go

import (
	"errors"
	uIoUtil "github.com/falconray0704/u4go/internal/ioutil"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_TempFile(t *testing.T) {
	fileLocation := "./tmp/"

	var contents = []byte(`
zap_log:
  logsPath: ./tmp/
  consoleFileOut: outConsole.logs
  consoleFileErr: errConsole.logs
  jsonFileOut: outJson.logs
  jsonFileErr: errJson.logs

dbs_infos:
  mysql:
    user_name: admin
    password: root
    url:	mysql.doryhub.com
    port:	2000
  redis:
    user_name: admin
    password: root
    url:	redis.doryhub.com
    port:	3000

`)
	path, clean, err := TempFile(fileLocation,"appCfgs-tmp", contents)
	defer clean()

	assert.NotEqual(t, "", path, "Default temp file creation expect not empty path.")
	assert.NotNil(t,clean, "Default temp file creation expect non-nil clean() function.")
	assert.Nil(t, err, "Default temp file creation expect nil err.")
}

func Test_newTempFile_Mock_getIoFileFunc_Err(t *testing.T) {
	fileLocation := "./tmp/"
	contents := []byte{}

	mock_getIoFile_Err := func(fileLocation, fileNamePrefix string) (uIoUtil.IoFile, error) {
		return nil, errors.New("Mocking get IoFile fail.")
	}

	path, clean, err := newTempFile(mock_getIoFile_Err,fileLocation,"appCfgs-tmp", contents)

	assert.Equal(t, "", path, "Get IoFile fail expect empty path.")
	assert.Nil(t, clean, "Get IoFile fail expect nil clean().")
	assert.Error(t, err, "Get IoFile fail expect non-nil  err.")
}

func Test_newTempFile_Mock_IoFile_Write(t *testing.T) {
	var (
		fileLocation = "./tmp/"
		contents = []byte("temp file contents")

		path string
		clean func()
		err error
		mock_getIoFile_Write GetIoFileFunc
	)



	mock_getIoFile_Write = func(fileLocation, fileNamePrefix string) (uIoUtil.IoFile, error) {
		ioFile := uIoUtil.GetIoFileMockAllSuccess()
		ioFile.FileName = fileLocation + fileNamePrefix
		return ioFile, nil
	}
	path, clean, err = newTempFile(mock_getIoFile_Write,fileLocation,"appCfgs-tmp", contents)
	assert.NotEqual(t, "", path, "Get IoFile success expect not empty path.")
	assert.NotNil(t, clean, "Get IoFile success expect non-nil clean().")
	assert.NoError(t, err, "Get IoFile success expect no err.")


	mock_getIoFile_Write = func(fileLocation, fileNamePrefix string) (uIoUtil.IoFile, error) {
		ioFile := uIoUtil.GetIoFileMockAllSuccess()
		ioFile.FileName = fileLocation + fileNamePrefix
		ioFile.WriteFunc = uIoUtil.WriteFuncMockErr
		return ioFile, nil
	}
	path, clean, err = newTempFile(mock_getIoFile_Write,fileLocation,"appCfgs-tmp", contents)
	assert.NotEqual(t, "", path, "Get IoFile success, but write fail expect not empty path.")
	assert.NotNil(t, clean, "Get IoFile success, but write fail expect non-nil clean().")
	assert.Error(t, err, "Get IoFile success, but write fail  expect non-nil err.")

	mock_getIoFile_Write = func(fileLocation, fileNamePrefix string) (uIoUtil.IoFile, error) {
		ioFile := uIoUtil.GetIoFileMockAllSuccess()
		ioFile.FileName = fileLocation + fileNamePrefix
		ioFile.WriteFunc = uIoUtil.WriteFuncMockLenNotEnough
		return ioFile, nil
	}
	path, clean, err = newTempFile(mock_getIoFile_Write,fileLocation,"appCfgs-tmp", contents)
	assert.NotEqual(t, "", path, "Get IoFile success, but write incorrect len expect not empty path.")
	assert.NotNil(t, clean, "Get IoFile success, but write incorrect len expect non-nil clean().")
	assert.Error(t, err, "Get IoFile success, but write incorrect len  expect non-nil err.")

	mock_getIoFile_Write = func(fileLocation, fileNamePrefix string) (uIoUtil.IoFile, error) {
		ioFile := uIoUtil.GetIoFileMockAllSuccess()
		ioFile.FileName = fileLocation + fileNamePrefix
		ioFile.WriteFunc = uIoUtil.WriteFuncMockSuccess
		ioFile.CloseFunc = uIoUtil.CloseFuncMockErr
		return ioFile, nil
	}
	path, clean, err = newTempFile(mock_getIoFile_Write,fileLocation,"appCfgs-tmp", contents)
	assert.NotEqual(t, "", path, "Get IoFile success, but close fail expect not empty path.")
	assert.NotNil(t, clean, "Get IoFile success, but close fail expect non-nil clean().")
	assert.Error(t, err, "Get IoFile success, but close fail expect non-nil err.")
}



