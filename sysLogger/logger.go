package sysLogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
)

func NewConfigLogger(consoleFileOut, consoleFileErr, jsonFileOut, jsonFileErr string) (logger *zap.Logger, closer []func(), err error) {

	defer func() {
		if p := recover(); p !=  nil {
			for _, cf := range closer {
				cf()
			}
			err, _ = p.(error)
		}
	}()
	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// Assume that we have clients for two Kafka topics. The clients implement
	// zapcore.WriteSyncer and are safe for concurrent use. (If they only
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	topicDebugging := zapcore.AddSync(ioutil.Discard)
	topicErrors := zapcore.AddSync(ioutil.Discard)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// High-priority output should also go to file "err.log", and low-priority
	// output should also go to file "out.log".
	fileOutConsoleSink, fileOutConsoleClose, fileOutConsoleErr := zap.Open(consoleFileOut)
	if fileOutConsoleErr != nil {
		panic(fileOutConsoleErr)
	}
	closer = append(closer, fileOutConsoleClose)

	fileErrConsoleSink, fileErrConsoleClose, fileErrConsoleErr :=  zap.Open(consoleFileErr)
	if fileErrConsoleErr != nil {
		panic(fileErrConsoleErr)
	}
	closer = append(closer, fileErrConsoleClose)

	// High-priority output JSON should also go to file "errJson.log", and low-priority
	// output JSON should also go to file "outJson.log".
	fileOutJsonSink, fileOutJsonClose, fileOutJsonErr := zap.Open(jsonFileOut)
	if fileOutJsonErr != nil {
		panic(fileOutJsonErr)
	}
	closer = append(closer, fileOutJsonClose)

	fileErrJsonSink, fileErrJsonClose, fileErrJsonErr :=  zap.Open(jsonFileErr)
	if fileErrJsonErr != nil {
		panic(fileErrJsonErr)
	}
	closer = append(closer, fileErrJsonClose)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, fileErrJsonSink, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, fileErrConsoleSink, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(kafkaEncoder, fileOutJsonSink, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, fileOutConsoleSink, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger = zap.New(core)
	//defer logger.Sync()
	return logger, closer, nil
}


