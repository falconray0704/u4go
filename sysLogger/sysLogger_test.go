package sysLogger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestInit_dev(t *testing.T) {
	var (
		err error
		logger *zap.Logger
		closeDev, closeRel func() error
	)

	logger, closeDev, err = Init(true)
	assert.Nil(t, err, "Init dev mode logger expert nil error")
	assert.NotNil(t, logger, "Init dev mode expert non-nil logger")
	Debug("constructed a logger")
	Info("constructed a logger")
	Warn("constructed a logger")
	Error("constructed a logger")
	Sync()
	closeDev()

	logger, closeRel, err = Init(false)
	assert.Nil(t, err, "Init rel mode logger expert nil error")
	assert.NotNil(t, logger, "Init rel mode expert non-nil logger")
	Debug("constructed a logger")
	Info("constructed a logger")
	Warn("constructed a logger")
	Error("constructed a logger")
	Sync()
	closeRel()

}

func TestSysLoggerDefaultFuncs(t *testing.T) {
	var err error
	defaultLogFunc("default log function.")
	err = defaultSyncFunc()
	assert.Nil(t, err, "Default sysLogger Sync() expect nil error")
	err = defaultCloseFunc()
	assert.Nil(t, err, "Default sysLogger Close() expect nil error")
}
