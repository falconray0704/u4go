package sysLogger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func TestNewSysLogCore(t *testing.T) {

	var (
		core zapcore.Core
		close func()
		err error
	)

	l := NewSysLogLevel("debug")

	core, close, err = NewSysLogCore(true, false, l, STDERR)
	assert.NotNil(t, core, "stderr output expect non-nil core.")
	assert.Nil(t, close, "stderr output expect nil close().")
	assert.Nil(t, err, "stderr output expect nil error.")

	core, close, err = NewSysLogCore(true, false, l, STDOUT)
	assert.NotNil(t, core, "stdout output expect non-nil core.")
	assert.Nil(t, close, "stdout output expect nil close().")
	assert.Nil(t, err, "stdout output expect nil error.")

	core, close, err = NewSysLogCore(true, false, l, "ftpftp://ftp")
	assert.Nil(t, core, "unsupported output expect nil core.")
	assert.Nil(t, close, "unsupported output expect nil close().")
	assert.NotNil(t, err, "unsupported output expect non-nil error.")

}

func TestNewFileSinker(t *testing.T) {
	var (
		sinker zapcore.WriteSyncer
		closer func()
		err		error
	)
	// empty file path
	sinker, closer, err = NewFileSinker("")
	assert.Nil(t, sinker, "Empty file path expect nil sinker.")
	assert.Nil(t, closer, "Empty file path expect nil closer.")
	assert.NotNil(t, err, "Empty file path expect non-nil error.")

	// unsupport output path
	sinker, closer, err = NewFileSinker("ftpftp://ftp")
	assert.Nil(t, sinker, "Unsupported file path expect nil sinker.")
	assert.Nil(t, closer, "Unsupported file path expect nil closer.")
	assert.NotNil(t, err, "Unsupported file path expect non-nil error.")

}

func TestNewSysLogLevelFail(t *testing.T) {
	l := NewSysLogLevel("abcd")
	assert.NotNil(t, l, "Unknown level expect fallback to non-nil debug level.")
}

func TestSysLoggerDefaultFuncs(t *testing.T) {
	var err error
	defaultLogFunc("default log function.")
	err = defaultSyncFunc()
	assert.Nil(t, err, "Default sysLogger Sync() expect nil error")
	err = defaultCloseFunc()
	assert.Nil(t, err, "Default sysLogger Close() expect nil error")
}
