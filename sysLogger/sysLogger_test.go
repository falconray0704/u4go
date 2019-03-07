package sysLogger

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
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

	var (
		st *StubbedTeeCoreBuilderError
	)
	st = WithStubTeeCoreBuilderError(Init, true)
	assert.Nil(t, st.retLogger, "WithStubbedTeeCoreBuilderError expect nil logger.")
	assert.Nil(t, st.retClose, "WithStubbedTeeCoreBuilderError expect nil close()")
	assert.NotNil(t, st.retErr, "WithStubbedTeeCoreBuilderError expect non-nil error")

}



type StubbedTeeCoreBuilderError struct {
	retLogger *zap.Logger
	retClose func() error
	retErr error

	pre		TeeCoreBuilder
}

func WithStubTeeCoreBuilderError(init InitFunc, isDevMode bool) *StubbedTeeCoreBuilderError {
	st := StubTeeCoreBuilderError()
	defer st.UnstubTeeCoreBuilderError()
	st.retLogger, st.retClose, st.retErr = init(isDevMode)
	return st
}

func StubTeeCoreBuilderError() *StubbedTeeCoreBuilderError {
	s := &StubbedTeeCoreBuilderError{pre:teeCoreBuilder}
	teeCoreBuilder = s.teeCoreBuilder
	return s
}

func (st *StubbedTeeCoreBuilderError) UnstubTeeCoreBuilderError() {
	teeCoreBuilder = st.pre
}

func (st *StubbedTeeCoreBuilderError) teeCoreBuilder(cfg *SysLogConfig) (zap.AtomicLevel, zapcore.Core, func(), error) {
	return NewSysLogLevel(cfg.LogLevel), nil, nil, errors.New("StubbedTeeCoreBuilder error")
}


func TestSysLogConfig_NewSysLogTeeCore(t *testing.T) {
	var (
		st *StubbedCoreBuilderError
	)

	st = WithStubCoreBuilderErrorSTDERR(NewSysLogTeeCore)
	assert.Nil(t, st.retCore, "Mocking stderr fail expect nil zapcore.Core return.")
	assert.Nil(t, st.retClose, "Mocking stderr fail expect nil close() return.")
	assert.NotNil(t, st.retErr, "Mocking stderr fail expect non-nil err return.")

	st = WithStubCoreBuilderErrorSTDOUT(NewSysLogTeeCore)
	assert.Nil(t, st.retCore, "Mocking stdout fail expect nil zapcore.Core return.")
	assert.Nil(t, st.retClose, "Mocking stdout fail expect nil close() return.")
	assert.NotNil(t, st.retErr, "Mocking stdout fail expect non-nil err return.")

	st = WithStubCoreBuilderErrorUnsupportedOutput(NewSysLogTeeCore)
	assert.Nil(t, st.retCore, "Mocking unsupported output file fail expect nil zapcore.Core return.")
	assert.Nil(t, st.retClose, "Mocking unsupported output file fail expect nil close() return.")
	assert.NotNil(t, st.retErr, "Mocking unsupported output file fail expect non-nil err return.")

	st = WithStubCoreBuilderErrorUnsupportedJson(NewSysLogTeeCore)
	assert.Nil(t, st.retCore, "Mocking unsupported output Json fail expect nil zapcore.Core return.")
	assert.Nil(t, st.retClose, "Mocking unsupported output Json fail expect nil close() return.")
	assert.NotNil(t, st.retErr, "Mocking unsupported output Json fail expect non-nil err return.")
}

type StubbedCoreBuilderError struct {
	retCore zapcore.Core
	retClose func()
	retErr error

	pre		CoreBuilder
}

func WithStubCoreBuilderErrorUnsupportedOutput(buildTeeCore TeeCoreBuilder) *StubbedCoreBuilderError {
	st := StubCoreBuilderErrorUnsupportedOutput()
	defer st.UnstubCoreBuilderError()
	cfg := NewDevSysLogConfigDefault()
	cfg.LogFilePrefix = "ftpftp://ftp"
	_, st.retCore, st.retClose, st.retErr = buildTeeCore(cfg)
	return st
}

func StubCoreBuilderErrorUnsupportedOutput() *StubbedCoreBuilderError {
	s := &StubbedCoreBuilderError{pre: coreBuilder}
	coreBuilder = newSysLogCore_unsupported_output_err
	return s
}

func WithStubCoreBuilderErrorUnsupportedJson(buildTeeCore TeeCoreBuilder) *StubbedCoreBuilderError {
	st := StubCoreBuilderErrorUnsupportedJson()
	defer st.UnstubCoreBuilderError()
	cfg := NewDevSysLogConfigDefault()
	_, st.retCore, st.retClose, st.retErr = buildTeeCore(cfg)
	return st
}

func StubCoreBuilderErrorUnsupportedJson() *StubbedCoreBuilderError {
	s := &StubbedCoreBuilderError{pre: coreBuilder}
	coreBuilder = newSysLogCore_unsupported_Json_err
	return s
}

func WithStubCoreBuilderErrorSTDOUT(buildTeeCore TeeCoreBuilder) *StubbedCoreBuilderError {
	st := StubCoreBuilderErrorSTDOUT()
	defer st.UnstubCoreBuilderError()
	cfg := NewDevSysLogConfigDefault()
	cfg.ConsoleOutput = STDOUT
	_, st.retCore, st.retClose, st.retErr = buildTeeCore(cfg)
	return st
}

func StubCoreBuilderErrorSTDOUT() *StubbedCoreBuilderError {
	s := &StubbedCoreBuilderError{pre: coreBuilder}
	coreBuilder = newSysLogCore_STDOUT_err
	return s
}

func WithStubCoreBuilderErrorSTDERR(buildTeeCore TeeCoreBuilder) *StubbedCoreBuilderError {
	st := StubCoreBuilderErrorSTDERR()
	defer st.UnstubCoreBuilderError()
	cfg := NewDevSysLogConfigDefault()
	_, st.retCore, st.retClose, st.retErr = buildTeeCore(cfg)
	return st
}

func StubCoreBuilderErrorSTDERR() *StubbedCoreBuilderError {
	s := &StubbedCoreBuilderError{pre: coreBuilder}
	coreBuilder = newSysLogCore_STDERR_err
	return s
}

func (st *StubbedCoreBuilderError) UnstubCoreBuilderError() {
	coreBuilder = st.pre
}

func newSysLogCore_unsupported_Json_err(isDevMode, isJsonEncoder bool, logLevel zap.AtomicLevel, logFilePath string) (zapcore.Core, func(), error) {
	if strings.Contains(logFilePath, "Json.log") {
		return nil, nil, errors.New("mocking unsupported Json file error from newSysLogCoreErr()")
	}
	return NewSysLogCore(isDevMode, isJsonEncoder, logLevel, logFilePath)
}

func newSysLogCore_unsupported_output_err(isDevMode, isJsonEncoder bool, logLevel zap.AtomicLevel, logFilePath string) (zapcore.Core, func(), error) {
	if strings.Contains(logFilePath, "ftpftp") {
		return nil, nil, errors.New("mocking unsupported output error from newSysLogCoreErr()")
	}
	return NewSysLogCore(isDevMode, isJsonEncoder, logLevel, logFilePath)
}

func newSysLogCore_STDOUT_err(isDevMode, isJsonEncoder bool, logLevel zap.AtomicLevel, logFilePath string) (zapcore.Core, func(), error) {
	if logFilePath == STDOUT {
		return nil, nil, errors.New("mocking logFilePath == STDOUT error from newSysLogCoreErr()")
	}
	return NewSysLogCore(isDevMode, isJsonEncoder, logLevel, logFilePath)
}

func newSysLogCore_STDERR_err(isDevMode, isJsonEncoder bool, logLevel zap.AtomicLevel, logFilePath string) (zapcore.Core, func(), error) {
	if logFilePath == STDERR {
		return nil, nil, errors.New("mocking logFilePath == STDERR error from newSysLogCoreErr()")
	}
	return NewSysLogCore(isDevMode, isJsonEncoder, logLevel, logFilePath)
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
