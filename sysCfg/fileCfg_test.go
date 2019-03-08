package sysCfg

import (
	"github.com/falconray0704/u4go"
	"github.com/stretchr/testify/assert"
	"testing"
)

type loggerCfg struct {
	IsDevMode bool	`yaml:"isDevMode"`
	LogsLocation string `yaml:"logsLocation"`
}

func TestLoadFileCfgs(t *testing.T) {
	var (
		errOnce error
		cfg loggerCfg

	)

	contents := []byte(`

sysLogger:
  isDevMode: true
  logsLocation: "./testDatas"

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

	cfgPath, clean, err := u4go.TempFile("./testDatas/", "appCfgs", contents)
	assert.NoError(t, err, "Create testing data expect always success.")
	assert.NotNil(t, clean, "Create testing data file expect non nil clean().")
	assert.NotEqual(t, "", cfgPath, "Create testing data file expect non empty file path.")
	defer clean()

	errOnce = LoadFileCfgs(cfgPath, "sysLogger", &cfg)
	assert.NoError(t, errOnce, "Load configs from test file expect no error.")

}
