package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	logger "rate-reader/internal/logger"
)

var Cfg *Config

const (
	LogLevelEnvKey  = "LOG_LEVEL"
	LogLevelDefault = logger.INFO
)

type Config struct {
	Db     DB     `mapstructure:"DB"`
	Port   string `mapstructure:"PORT"`
	Delay  uint   `mapstructure:"DELAY"`
	Source string `mapstructure:"SOURCE"`
}

func (c *Config) String() string {
	return fmt.Sprintf("PORT: %s, DELAY: %s, SOURCE: %s, DB: %v",
		c.Port, c.Delay, c.Source, c.Db,
	)
}

// Load set configuration parameters.
// At first read config from file
// After that read environment variables
func Load(defaultConfigPath string) error {
	cfg, err := read(defaultConfigPath)
	if err != nil {
		return err
	}
	Cfg = cfg
	return validate()
}

func read(defaultConfigPath string) (*Config, error) {
	cfg := new(Config)

	// read config from file - it will be default values
	viper.SetConfigFile(defaultConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// read parameters from environment variables -> they override default values from file
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validate() error {
	if len(Cfg.Db.Host) == 0 {
		return errors.New("DB_HOST parameter is empty")
	}
	if len(Cfg.Db.Name) == 0 {
		return errors.New("DB_NAME parameter is empty")
	}
	if len(Cfg.Source) == 0 {
		return errors.New("SOURCE parameter is empty")
	}
	if Cfg.Delay == 0 {
		return errors.New("DELAY parameter can't be 0")
	}
	return nil
}
