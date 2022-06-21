package main

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var defaultConfigFile = "./config.yaml"

func readFromFile(v *viper.Viper) {
	// default config path
	configPath := v.GetString("config-path")
	if configPath == "" {
		configPath = defaultConfigFile
	}
	// if file doesn't exist just continue
	_, err := os.Stat(configPath)
	if err != nil {
		return
	}
	// set config type from file extension
	ext := filepath.Ext(configPath)[1:]
	v.SetConfigType(ext)
	// set file path
	v.SetConfigFile(configPath)
	// try read config or fatal
	err = v.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config (%s): %+v", configPath, err)
	}
}

func NewViper() *viper.Viper {
	v := viper.New()
	// enable parse from environment variable
	v.AutomaticEnv()
	// replace "." and "-" with "_" for envs
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	// try read config where its possible
	readFromFile(v)
	return v
}
