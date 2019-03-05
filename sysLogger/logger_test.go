package sysLogger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfigLogger_fail_on_consoleFileOut(t *testing.T) {

	logsLocation := DefaultLogsLocation
	consoleFileOut := "ftpftp://log.root"
	consoleFileErr := logsLocation + "errConsole.logs"
	jsonFileOut := logsLocation + "outJson.logs"
	jsonFileErr := logsLocation + "errJson.logs"

	logger, _, err := NewConfigLogger(consoleFileOut, consoleFileErr, jsonFileOut, jsonFileErr)
	assert.NotNilf(t, err, "Open an no permission file:%s for consoleFileOut expect non-nil error", consoleFileOut)
	assert.Nilf(t, logger, "Open an no permission file:%s for consoleFileOut expect nil logger", consoleFileOut)

}

func TestNewConfigLogger_fail_on_consoleFileErr(t *testing.T) {

	logsLocation := DefaultLogsLocation
	consoleFileOut := logsLocation + "outConsole.logs"
	consoleFileErr := "ftpftp://log.root"
	jsonFileOut := logsLocation + "outJson.logs"
	jsonFileErr := logsLocation + "errJson.logs"

	logger, _, err := NewConfigLogger(consoleFileOut, consoleFileErr, jsonFileOut, jsonFileErr)
	assert.NotNilf(t, err, "Open an no permission file:%s for consoleFileErr expect non-nil error", consoleFileErr)
	assert.Nilf(t, logger, "Open an no permission file:%s for consoleFileErr expect nil logger", consoleFileErr)

}

func TestNewConfigLogger_fail_on_jsonFileOut(t *testing.T) {

	logsLocation := DefaultLogsLocation
	consoleFileOut := logsLocation + "outConsole.logs"
	consoleFileErr := logsLocation + "errConsole.logs"
	jsonFileOut := "ftpftp://log.root"
	jsonFileErr := logsLocation + "errJson.logs"

	logger, _, err := NewConfigLogger(consoleFileOut, consoleFileErr, jsonFileOut, jsonFileErr)
	assert.NotNilf(t, err, "Open an no permission file:%s for jsonFileOut expect non-nil error", jsonFileOut)
	assert.Nilf(t, logger, "Open an no permission file:%s for jsonFileOut expect nil logger", jsonFileOut)

}

func TestNewConfigLogger_fail_on_jsonFileErr(t *testing.T) {

	logsLocation := DefaultLogsLocation
	consoleFileOut := logsLocation + "outConsole.logs"
	consoleFileErr := logsLocation + "errConsole.logs"
	jsonFileOut := logsLocation + "outJson.logs"
	jsonFileErr := "ftpftp://log.root"

	logger, _, err := NewConfigLogger(consoleFileOut, consoleFileErr, jsonFileOut, jsonFileErr)
	assert.NotNilf(t, err, "Open an no permission file:%s for jsonFileErr expect non-nil error", jsonFileErr)
	assert.Nilf(t, logger, "Open an no permission file:%s for jsonFileErr expect nil logger", jsonFileErr)

}

func TestNewConfigLogger(t *testing.T) {

	logsLocation := DefaultLogsLocation
	consoleFileOut := logsLocation + "outConsole.logs"
	consoleFileErr := logsLocation + "errConsole.logs"
	jsonFileOut := logsLocation + "outJson.logs"
	jsonFileErr := logsLocation + "errJson.logs"

	logger, _, err := NewConfigLogger(consoleFileOut, consoleFileErr, jsonFileOut, jsonFileErr)
	assert.Nilf(t, err, "Testing the full success path expect nil error")
	assert.NotNilf(t, logger, "Testing the full success path expect non-nil logger")

	logger.Info("constructed a logger")
	logger.Info("constructed a logger 2")
	logger.Warn("constructed a logger")
	logger.Warn("constructed a logger 2")
	logger.Error("constructed a logger")
	logger.Error("constructed a logger 2")

}
