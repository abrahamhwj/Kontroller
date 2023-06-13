// Package config provides configuration settings for the application
package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// Log represents the log settings
type Log struct {
	Level int32 `yaml:"level"`
}

// Manager represents the manager settings
type Manager struct {
	ThreadNumber            int32         `yaml:"threadNumber"`
	ControllerMaxRetryTimes int32         `yaml:"controllerMaxRetryTimes"`
	ThreadTimeout           time.Duration `yaml:"threadTimeout"`
	ReSyncPeriod            time.Duration `yaml:"reSyncPeriod"`
}

// Config represents the overall configuration
type Config struct {
	Log     Log     `yaml:"log"`
	Manager Manager `yaml:"manager"`
}

// Cfg is the global configuration variable
var Cfg Config

// init initializes the configuration
func init() {
	// Set the configuration file name and path
	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("..")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	// Set default values for configuration settings
	viper.SetDefault("log.level", 4)
	viper.SetDefault("manager.threadNumber", 1)
	viper.SetDefault("manager.controllerMaxRetryTimes", 5)
	viper.SetDefault("manager.threadTimeout", 5)
	viper.SetDefault("manager.reSyncPeriod", 300)
	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config err with: %s", err))
	}
	// Unmarshal the configuration into the Cfg variable
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(fmt.Errorf("decode config err with: %s", err))
	}
	// Convert time.Duration settings from seconds to actual duration
	Cfg.Manager.ThreadTimeout *= time.Second
	Cfg.Manager.ReSyncPeriod *= time.Second
}
