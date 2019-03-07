package main

import (
	"fmt"
	slog "github.com/falconray0704/u4go/sysLogger"
	"os"
)

/*
func demoPresets(sysCfg *Config) {
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

func demoBasicConfiguratoins(sysCfg *Config) {
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
	  "outputPaths": ["stdout"],
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
	cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, sysCfg.ZapLog.LogsPath + "logs")
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
}

func demoGetCommandLineArgs(sysCfg *cfg.Config) {
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

	//fmt.Println("word:", *wordPtr)
	//fmt.Println("numb:", *numbPtr)
	//fmt.Println("fork:", *boolPtr)
	//fmt.Println("svar:", svar)
	//fmt.Println("tail:", flag.Args())

	logger, closer, err := NewConfigLogger(sysCfg)
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


func demoAdvancedConfigurations(logger *zap.Logger) {

	ymlPtr := flag.String("c", "./sysDatas/cfgs/appCfgs.yaml", "yaml file to read config from")
	flag.Parse()
	sysCfg, err := cfg.NewConfig(*ymlPtr)

	if err != nil {
		log.Fatalln(err)
	}

	logsLocation := sysCfg.ZapLog.LogsPath
	consoleFileOut := logsLocation + "outConsole.logs"
	consoleFileErr := logsLocation + "errConsole.logs"
	jsonFileOut := logsLocation + "outJson.logs"
	jsonFileErr := logsLocation + "errJson.logs"
	logger, closers, err := slog.NewConfigLogger(consoleFileOut, consoleFileErr, jsonFileOut, jsonFileErr)

	if err != nil {
		fmt.Errorf("Error:%v", err)
	}
	defer func() {
		for _, cf := range closers {
			cf()
		}
	}()

	slog.Info("constructed a logger")
	slog.Info("constructed a logger 2")
	slog.Warn("constructed a logger")
	slog.Warn("constructed a logger 2")
	slog.Error("constructed a logger")
	slog.Error("constructed a logger 2")
}
*/

func main() {
	/*
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
	*/
	//demoGetCommandLineArgs(sysCfg)

	_, _, err := slog.Init(true)
	if err != nil {
		fmt.Printf("Init system logger fail: %s.\n", err.Error())
		os.Exit(1)
	}
	defer func() {
		slog.Sync()
		slog.Close()
	}()

	slog.Debug("constructed a logger")
	slog.Info("constructed a logger")
	slog.Warn("constructed a logger")
	slog.Error("constructed a logger")

}
