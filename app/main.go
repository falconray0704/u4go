package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	workingDir string
	fileLogs = "/tmp/logs"
	fileOutConsole = "/tmp/outConsole.logs"
	fileErrConsole = "/tmp/errConsole.logs"
	fileOutJson = "/tmp/outJson.logs"
	fileErrJson = "/tmp/errJson.logs"
)

func demoPresets() {
	// Using zap's preset constructors is the simplest way to get a feel for the
	// package, but they don't allow much customization.
	logger := zap.NewExample() // or NewProduction, or NewDevelopment
	defer logger.Sync()

	const url = "http://example.com"

	// In most circumstances, use the SugaredLogger. It's 4-10x faster than most
	// other structured logging packages and has a familiar, loosely-typed API.
	sugar := logger.Sugar()
	sugar.Infow("Failed to fetch URL.",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	// In the unusual situations where every microsecond matters, use the
	// Logger. It's even faster than the SugaredLogger, but only supports
	// structured logging.
	logger.Info("Failed to fetch URL.",
		// Structured context as strongly typed fields.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}

func demoBasicConfiguratoins() {
	// For some users, the presets offered by the NewProduction, NewDevelopment,
	// and NewExample constructors won't be appropriate. For most of those
	// users, the bundled Config struct offers the right balance of flexibility
	// and convenience. (For more complex needs, see the AdvancedConfiguration
	// example.)
	//
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.

	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/mnt/ld0/gows/src/github.com/falconray0704/u4go/app/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
}

func NewConfigLogger() (logger *zap.Logger, closer []func(), err error) {

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
	fileOutConsoleSink, fileOutConsoleClose, fileOutConsoleErr := zap.Open(fileOutConsole)
	if fileOutConsoleErr != nil {
		panic(fileOutConsoleErr)
	}
	closer = append(closer, fileOutConsoleClose)

	fileErrConsoleSink, fileErrConsoleClose, fileErrConsoleErr :=  zap.Open(fileErrConsole)
	if fileErrConsoleErr != nil {
		panic(fileErrConsoleErr)
	}
	closer = append(closer, fileErrConsoleClose)

	// High-priority output JSON should also go to file "errJson.log", and low-priority
	// output JSON should also go to file "outJson.log".
	fileOutJsonSink, fileOutJsonClose, fileOutJsonErr := zap.Open(fileOutJson)
	if fileOutJsonErr != nil {
		panic(fileOutJsonErr)
	}
	closer = append(closer, fileOutJsonClose)

	fileErrJsonSink, fileErrJsonClose, fileErrJsonErr :=  zap.Open(fileErrJson)
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

func demoAdvancedConfigurations() {

	logger, closer, err := NewConfigLogger()
	if err != nil {
		fmt.Errorf("Error:%v", err)
	}
	defer func() {
		for _, cf := range closer {
			cf()
		}
	}()

	logger.Info("constructed a logger")
	logger.Info("constructed a logger 2")
	logger.Warn("constructed a logger")
	logger.Warn("constructed a logger 2")
	logger.Error("constructed a logger")
	logger.Error("constructed a logger 2")
}

func demoGetCommandLineArgs() {
	// Basic flag declarations are available for string,
	// integer, and boolean options. Here we declare a
	// string flag `word` with a default value `"foo"`
	// and a short description. This `flag.String` function
	// returns a string pointer (not a string value);
	// we'll see how to use this pointer below.
	wordPtr := flag.String("word", "foo", "a string")

	// This declares `numb` and `fork` flags, using a
	// similar approach to the `word` flag.
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	loopPtr := flag.Bool("loop", false, "a bool")

	// It's also possible to declare an option that uses an
	// existing var declared elsewhere in the program.
	// Note that we need to pass in a pointer to the flag
	// declaration function.
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

	// Here we'll just dump out the parsed options and
	// any trailing positional arguments. Note that we
	// need to dereference the pointers with e.g. `*wordPtr`
	// to get the actual option values.
	/*
	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
	*/

	logger, closer, err := NewConfigLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		for _, cf := range closer {
			cf()
		}
	}()

	if *loopPtr {
		for {
			logger.Info("Startup flags:",
				zap.Bool("loop",  *loopPtr),
				zap.String("word:", *wordPtr),
				zap.Int("numb:", *numbPtr),
				zap.Bool("fork:", *boolPtr),
				zap.String("svar", svar),
				zap.Strings("tail:", flag.Args()))

			time.Sleep(time.Second * 1)
		}
	} else {
		logger.Info("Startup flags:",
			zap.Bool("loop",  *loopPtr),
			zap.String("word:", *wordPtr),
			zap.Int("numb:", *numbPtr),
			zap.Bool("fork:", *boolPtr),
			zap.String("svar", svar),
			zap.Strings("tail:", flag.Args()))
	}
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Working directory:", dir)

	workingDir = dir
	fileLogs = workingDir + "/tmp/logs"
	fileOutConsole = workingDir + "/tmp/outConsole.logs"
	fileErrConsole = workingDir + "/tmp/errConsole.logs"
	fileOutJson = workingDir + "/tmp/outJson.logs"
	fileErrJson = workingDir + "/tmp/errJson.logs"

	//demoPresets()
	//demoBasicConfiguratoins()
	//demoAdvancedConfigurations()
	demoGetCommandLineArgs()
}
