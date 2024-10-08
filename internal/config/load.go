package config

import (
	"flag"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Load provides a loaded configuration for server.
func Load() (*Config, error) {
	var config Config

	cfg, err := readConfigFile()
	if err != nil {
		return nil, err
	}

	if err = cfg.Unmarshal(&config, setupConfigUnmarshaler); err != nil {
		log.Println("can not extract configuration")
		log.Println(err.Error())
		return nil, err
	}

	return &config, nil
}

// setupConfigUnmarshaler configures the Config loader to properly unmarshal
// special types we use for the API server
func setupConfigUnmarshaler(cfg *mapstructure.DecoderConfig) {
	// add the decoders missing here
	cfg.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		cfg.DecodeHook)
}

// readConfigFile reads the config file and provides instance of the loaded configuration.
func readConfigFile() (*viper.Viper, error) {
	log.Print("loading app configuration")

	// Get the config reader
	cfg := reader()

	// set default values
	applyDefaults(cfg)

	// Try to read the file
	if err := cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Print("can not read the server configuration")
			return nil, err
		}

		log.Print("configuration file not found, using default values")
	}

	return cfg, nil
}

// reader provides instance of the config reader.
// It accepts an explicit path to a config file if it was requested by `cfg` flag.
func reader() *viper.Viper {
	cfg := viper.New()

	cfg.SetConfigName(configFileName)

	cfg.AddConfigPath(defaultConfigDir())
	cfg.AddConfigPath(".")

	var cfgPath string
	flag.StringVar(&cfgPath, keyConfigFilePath, "", "Path to a configuration file")
	flag.Parse()

	cfg.SetConfigFile(cfgPath)

	return cfg
}
