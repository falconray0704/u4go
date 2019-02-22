package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	path, clean := TempFile(t, "appCfgs-tmp",`

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
	defer clean()

	config, err := NewConfig(path)
	assert.NoError(t, err)

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

