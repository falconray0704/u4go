package u4go

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTempFile_NewFile(t *testing.T) {
	fileLocation := "./tmp/"
	var tmpFile *TmpFile = nil

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
	path, clean, err := tmpFile.NewFile(fileLocation,"appCfgs-tmp", contents)
	defer clean()

	assert.NotEqual(t, "", path, "Default temp file creation expect not empty path.")
	assert.NotNil(t,clean, "Default temp file creation expect non-nil clean() function.")
	assert.Nil(t, err, "Default temp file creation expect nil err.")
}


