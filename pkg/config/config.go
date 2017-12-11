package config

import (
	"os"

	"fmt"

	"path/filepath"

	"path"
	"runtime"

	"errors"

	"net/url"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database Database
	Client   Client
}

var NewErrConfigFileNotFound = func(environment string) error {
	return errors.New(fmt.Sprintf("No configuration found for environment: %v", environment))
}

var NewErrBadConnectionString = func(connectionString string) error {
	return errors.New(fmt.Sprintf("connection string is invalid: %v", connectionString))
}

func NewConfig(environment string) (Config, error) {
	filenameWithExtension := fmt.Sprintf("%s.toml", environment)
	absolutePath := filepath.Join(ProjectRoot(), "environments", filenameWithExtension)
	config, err := parseConfigFile(absolutePath)
	if err != nil {
		return Config{}, NewErrConfigFileNotFound(environment)
	} else {
		if !filepath.IsAbs(config.Client.IPCPath) && !isUrl(config.Client.IPCPath) {
			config.Client.IPCPath = filepath.Join(ProjectRoot(), config.Client.IPCPath)
		}
		return config, nil
	}
}

func ProjectRoot() string {
	var _, filename, _, _ = runtime.Caller(0)
	return path.Join(path.Dir(filename), "..", "..")
}

func isUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err == nil {
		return true
	}
	return false
}

func fileExists(s string) bool {
	_, err := os.Stat(s)
	if err == nil {
		return true
	}
	return false
}

func parseConfigFile(filePath string) (Config, error) {
	var cfg Config
	if !isUrl(filePath) && !fileExists(filePath) {
		return Config{}, NewErrBadConnectionString(filePath)
	} else {
		_, err := toml.DecodeFile(filePath, &cfg)
		if err != nil {
			return Config{}, err
		}
		return cfg, nil
	}
}
