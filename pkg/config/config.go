package config

import (
	"log"
	"os"

	"fmt"

	"path/filepath"

	"path"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database Database
	Client   Client
}

func NewConfig(environment string) Config {
	filenameWithExtension := fmt.Sprintf("%s.toml", environment)
	absolutePath := filepath.Join(ProjectRoot(), "pkg", "config", "environments", filenameWithExtension)
	config := parseConfigFile(absolutePath)
	config.Client.IPCPath = filepath.Join(ProjectRoot(), config.Client.IPCPath)
	return config
}

func ProjectRoot() string {
	var _, filename, _, _ = runtime.Caller(0)
	return path.Join(path.Dir(filename), "..", "..")
}

func parseConfigFile(configfile string) Config {
	var cfg Config
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	if _, err := toml.DecodeFile(configfile, &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
