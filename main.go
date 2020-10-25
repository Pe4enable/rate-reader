package rate_reader

import (
	"context"
	"flag"
	"fmt"
	"github.com/rate-reader/config"
	logger "github.com/rate-reader/logger"
	repository "github.com/rate-reader/repositories"
	services "github.com/rate-reader/services"
	"os"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "./config/config.yml", "A path to config file")
	isDebugMode := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	log, err := initLogger(*isDebugMode)
	if err != nil {
		fmt.Println("Can't initialize logger", err)
		return
	}

	if err := config.Load(configPath); err != nil {
		log.Fatal("Loading of configuration failed with error:", err)
	}
	log.Infof("Rate reader will be started with configuration %s", config.Cfg.String())
	ctx := context.Background()
	ctx = logger.ToContext(ctx, log)

	if err := repository.New(ctx, config.Cfg.Db); err != nil {
		log.Fatal("Can't create repository: ", err)
	}

	if err := services.NewReader(ctx, config.Cfg, repository.GetRepository()); err != nil {
		log.Fatal("Can't create rate reader: ", err)
	}

	rateReader := services.GetRateReader()

	err = rateReader.Start(ctx)
	if err != nil {
		log.Fatal("Can't start rate reader: ", err)
	}
	defer rateReader.Stop(ctx)

}

// initLogger initializes logger: create logger, set logger format: json or text.
// text is used if application was started with flag '-debug'
// set log level according to environment variable LOG_LEVEL,
// if LOG_LEVEL was not set it uses INFO by default,
// if application was started with flag '-debug' set DEBUG level
func initLogger(isDebug bool) (logger.ILogger, error) {
	var level = config.LogLevelDefault
	if isDebug {
		level = logger.DEBUG
	}
	if l, ok := os.LookupEnv(config.LogLevelEnvKey); ok {
		level = logger.Level(l)
	}
	return logger.Init(!isDebug, level)
}
