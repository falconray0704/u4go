package cfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/falconray0704/u4go"
)

func TestNewConfig(t *testing.T) {
	fileLocation := "./tmp/"
	contents := []byte(`
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
	path, clean, err := u4go.TempFile(fileLocation,"appCfgs-tmp", contents)
	assert.NoError(t, err, "Config file should be existed.")
	assert.NotNil(t, clean, "Loading config file should be success.")
	defer clean()

	config, err := NewConfig(path)
	assert.NoErrorf(t, err, "Parsing config file should be success with correct contents: %s", contents)

	assert.Equal(t, "./tmp/",config.ZapLog.LogsPath )
	assert.Equal(t, "outConsole.logs", config.ZapLog.ConsoleFileOut)
	assert.Equal(t, "errConsole.logs", config.ZapLog.ConsoleFileErr)
	assert.Equal(t, "outJson.logs", config.ZapLog.JsonFileOut)
	assert.Equal(t, "errJson.logs", config.ZapLog.JsonFileErr)

	dbmysql, ok := config.DBsInfos["mysql"]
	assert.True(t, ok)
	assert.Equal(t, DBInfo{UserName: "admin", Password: "root", Url: "mysql.doryhub.com", Port: 2000}, dbmysql)

	dbredis, ok := config.DBsInfos["redis"]
	assert.True(t, ok)
	assert.Equal(t, DBInfo{UserName: "admin", Password: "root", Url: "redis.doryhub.com", Port: 3000}, dbredis)
}

func TestNewConfigReadDataError(t *testing.T) {
	path :="./tmp/noAppCfg.yaml"

	config, err := NewConfig(path)
	assert.NotNilf(t, err, "Load the non-exsit config file:%s expect non-nil error.", path)
	assert.Nil(t, config, "Load the non-exsit config file:%s expect nil config .", path)

}

func TestNewConfigParseDataError(t *testing.T) {
	fileLocation := "./tmp/"
	contents := []byte(
`
	incorrect datas should not be parsed.
`)
	path, clean, err := u4go.TempFile(fileLocation,"appCfgs-tmp", contents)
	assert.NoError(t, err, "Config file should be existed.")
	assert.NotNil(t, clean, "Loading config file should be success.")
	defer clean()

	config, err := NewConfig(path)
	assert.Error(t, err, "Parsing config file should be fail with incorrect contents.")

	assert.Nil(t, config, "Load the incorrect config file:%s expect nil config.", contents)
	assert.NotNil(t, err, "Load the incorrect config file:%s expect non-nil error.", contents)

}
