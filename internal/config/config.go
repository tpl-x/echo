package config

import (
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
	Database DatabaseConfig `yaml:"database"`
}

type LogConfig struct {
	FileName    string `yaml:"file_name"`
	MaxSize     int    `yaml:"max_size"`
	MaxBackups  int    `yaml:"max_backups"`
	MaxKeepDays int    `yaml:"max_keep_days"`
	Compress    bool   `yaml:"compress"`
}

type ServerConfig struct {
	BindPort         int32 `yaml:"bind_port"`
	GraceExitTimeout int   `yaml:"grace_exit_timeout"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	MaxIdle  int    `yaml:"max_idle"`
	MaxOpen  int    `yaml:"max_open"`
}

// LoadFromReader  load config from reader
func LoadFromReader(reader io.Reader) (*AppConfig, error) {
	config := &AppConfig{}
	if err := yaml.NewDecoder(reader).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

// LoadFromFile Load config from file path
func LoadFromFile(path string) (*AppConfig, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return LoadFromReader(f)
}

// ProvideConfig provides config with fx
func ProvideConfig(path string) func() (*AppConfig, error) {
	return func() (*AppConfig, error) {
		return LoadFromFile(path)
	}
}

// Module provides config as fx module
var Module = fx.Module("config")
