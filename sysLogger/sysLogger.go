package sysLogger

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	STDOUT = "stdout"
	STDERR = "stderr"

	DefaultLogsLocation = "./logDatas/"
)

type InitFunc func(isDevMode bool) (logger *zap.Logger, closeLogger func() error, err error)
type LogFieldsFunc func(msg string, fields ...zap.Field)
type TeeCoreBuilder func(cfg *SysLogConfig) (zap.AtomicLevel, zapcore.Core, func(), error)
type CoreBuilder func(isDevMode, isJsonEncoder bool, logLevel zap.AtomicLevel, logFilePath string) (zapcore.Core, func(), error)

var (
	coreBuilder CoreBuilder = NewSysLogCore
	teeCoreBuilder TeeCoreBuilder = NewSysLogTeeCore
	Debug	LogFieldsFunc = defaultLogFunc
	Info	LogFieldsFunc = defaultLogFunc
	Warn	LogFieldsFunc = defaultLogFunc
	Error	LogFieldsFunc = defaultLogFunc
	DPanic	LogFieldsFunc = defaultLogFunc
	Panic	LogFieldsFunc = defaultLogFunc
	Fatal	LogFieldsFunc = defaultLogFunc
	Sync	= defaultSyncFunc
	Close	= defaultCloseFunc

	sysLogger *SysLogger
	Log *zap.Logger

)

func defaultLogFunc(msg string, fields ...zap.Field) {
	fmt.Printf("Default logger, msg:%s, fields:%+v\n", msg, fields)
}

func defaultSyncFunc() error {
	fmt.Println("Default logger sync() do nothing.")
	return nil
}

func defaultCloseFunc() error {
	fmt.Println("Default logger Close() do nothing.")
	return nil
}

func init() {
	sysLogger = new(SysLogger)
}

type SysLogger struct {
	SysLogCfg *SysLogConfig
	CurrentLogLevel zap.AtomicLevel
	ZapLogger *zap.Logger
	//Closer func()
}

func Init(isDevMode bool) (logger *zap.Logger, closeLogger func() error, err error) {
	var (
		sysLogCfg *SysLogConfig
		sysLogLevel zap.AtomicLevel
		sysLogTeeCore zapcore.Core
		sysLogCloser func()

		errOnce error
	)

	if !isDevMode {
		sysLogCfg = NewRelSysLogConfigDefault()
	} else {
		sysLogCfg = NewDevSysLogConfigDefault()
	}

	if sysLogLevel, sysLogTeeCore, sysLogCloser, errOnce = sysLogCfg.buildTeeCore(sysLogCfg); errOnce != nil  {
		return nil, nil, errOnce
	}

	sysLogger.SysLogCfg = sysLogCfg
	sysLogger.CurrentLogLevel = sysLogLevel
	sysLogger.ZapLogger = zap.New(sysLogTeeCore)

	Log = sysLogger.ZapLogger

	Debug = sysLogger.ZapLogger.Debug
	Info = sysLogger.ZapLogger.Info
	Warn = sysLogger.ZapLogger.Warn
	Error = sysLogger.ZapLogger.Error
	DPanic = sysLogger.ZapLogger.DPanic
	Panic = sysLogger.ZapLogger.Panic
	Fatal = sysLogger.ZapLogger.Fatal
	Sync = sysLogger.ZapLogger.Sync

	Close = func() error {
		sysLogCloser()
		return nil
	}

	return Log, Close, nil
}

type SysLogConfig struct {
	IsDevMode bool 			`yaml:"isDevMode"`
	LogLevel string			`yaml:"logLevel"`

	EnableConsole bool		`yaml:"enableConsole"`
	EnableConsoleFile bool	`yaml:"enableConsoleFile"`
	EnalbeJsonFile bool		`yaml:"enableJsonFile"`

	LogsLocation string		`yaml:"logsLocation"` // logs storage location , must end with "\" or "/" which depend on OS
	LogFilePrefix string 	`yaml:"logFilePrefix"` // prefix of log files
	ConsoleOutput string	`yaml:"consoleOutput"` // only support "stdout" or "stderr"

	buildCore	CoreBuilder
	buildTeeCore TeeCoreBuilder
}

func NewRelSysLogConfigDefault() *SysLogConfig  {
	return &SysLogConfig{
		IsDevMode: false,
		LogLevel: "info",
		EnableConsole: false,
		EnableConsoleFile: true,
		EnalbeJsonFile: true,
		LogsLocation: DefaultLogsLocation,
		LogFilePrefix: "rel",
		ConsoleOutput: STDERR,
		buildCore: coreBuilder,
		buildTeeCore: teeCoreBuilder}
}

func NewDevSysLogConfigDefault() *SysLogConfig  {
	return &SysLogConfig{
		IsDevMode: true,
		LogLevel: "debug",
		EnableConsole: true,
		EnableConsoleFile: true,
		EnalbeJsonFile: true,
		LogsLocation: DefaultLogsLocation,
		LogFilePrefix: "dev",
		ConsoleOutput: STDERR,
		buildCore: coreBuilder,
		buildTeeCore: teeCoreBuilder}
}

func NewSysLogTeeCore(cfg *SysLogConfig) (zap.AtomicLevel, zapcore.Core, func(), error) {
	var (
		errOnce, multiErr error

		coreConsole  	zapcore.Core
		coreConsoleFile zapcore.Core
		coreJsonFile	zapcore.Core

		coreClose		func()
		closers			= []func(){}
		cores			[]zapcore.Core

		logLevel		zap.AtomicLevel
	)

	closerHub := func() {
		for _, cls := range closers {
			cls()
		}
	}

	defer func() {
		if multiErr != nil {
			closerHub()
		}
	}()

	logLevel = NewSysLogLevel(cfg.LogLevel)

	if cfg.EnableConsole {
		coreConsole, coreClose, errOnce = cfg.buildCore(cfg.IsDevMode, false, logLevel, cfg.ConsoleOutput)
		if errOnce != nil {
			multiErr = multierr.Append(multiErr, errOnce)
			return logLevel, nil, nil, multiErr
		} else {
			cores = append(cores, coreConsole)
		}
	}

	if cfg.EnableConsoleFile {
		coreConsoleFile, coreClose, errOnce = cfg.buildCore(cfg.IsDevMode, false, logLevel, cfg.LogsLocation + cfg.LogFilePrefix + "Console.log")
		if errOnce != nil {
			multiErr = multierr.Append(multiErr, errOnce)
			return logLevel, nil, nil, multiErr
		} else {
			closers = append(closers, coreClose)
			cores = append(cores, coreConsoleFile)
		}
	}

	if cfg.EnalbeJsonFile {
		coreJsonFile, coreClose, errOnce = cfg.buildCore(cfg.IsDevMode, true, logLevel, cfg.LogsLocation + cfg.LogFilePrefix + "Json.log")
		if errOnce != nil {
			multiErr = multierr.Append(multiErr, errOnce)
			return logLevel, nil, nil, multiErr
		} else {
			closers = append(closers, coreClose)
			cores = append(cores, coreJsonFile)
		}
	}

	return logLevel,zapcore.NewTee(cores ...), closerHub, nil
}

func NewSysLogCore(isDevMode, isJsonEncoder bool, logLevel zap.AtomicLevel, logFilePath string) (zapcore.Core, func(), error) {

	var (
		sinker zapcore.WriteSyncer
		closer func()
		err error
	)

	if logFilePath == STDERR {
		sinker = zapcore.Lock(os.Stderr)
	} else if logFilePath == STDOUT {
		sinker = zapcore.Lock(os.Stderr)
	} else {
		if sinker, closer, err = NewFileSinker(logFilePath); err != nil {
			return nil, nil, err
		}
	}

	logEncoder := NewSysLogEncoder(isDevMode, isJsonEncoder)

	return zapcore.NewCore(logEncoder, sinker, logLevel), closer, nil
}

// Construct log file sinker
func NewFileSinker(logFilePath string) (sink zapcore.WriteSyncer, close func(), err error) {
	var (
		sinker zapcore.WriteSyncer
		closer func()
		errOnce, multiErr error
	)

	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(error); ok {
				multiErr = multierr.Append(multiErr, e)
				sink = nil
				close = nil
				err = multiErr
			}
		}
	}()

	if logFilePath == "" {
		multiErr = multierr.Append(multiErr, errors.New("file path of log should not be empty"))
		return nil, nil, multiErr
	}

	sinker, closer, errOnce = zap.Open(logFilePath)
	if errOnce != nil {
		panic(errOnce)
	}

	return sinker, closer, nil
}

// Construct encoder
func NewSysLogEncoder(isDevMode bool, isJsonEncoder bool) zapcore.Encoder {

	var (
		encoderConfig 	zapcore.EncoderConfig
		encoder 		zapcore.Encoder
	)

	if isDevMode {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	if isJsonEncoder {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return encoder
}

// Always get supported level. If fail, fallback to "debug"
func NewSysLogLevel(logLevelStr string) zap.AtomicLevel {

	var (
		errOnce error
	)

	logLevel := zap.NewAtomicLevel()
	if errOnce = logLevel.UnmarshalText([]byte(logLevelStr)); errOnce != nil {
		logLevel.UnmarshalText([]byte("debug"))
	}

	return logLevel
}

func SetCurrentLevel(levelStr string) error {
	return sysLogger.CurrentLogLevel.UnmarshalText([]byte(levelStr))
}

func GetCurrentLevel() string {
	return sysLogger.CurrentLogLevel.String()
}


