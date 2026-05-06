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
	JWT      JWTConfig      `yaml:"jwt"`
}

type LogConfig struct {
	FileName    string `yaml:"file_name"`
	Level       string `yaml:"level"`
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

type JWTConfig struct {
	AccessSecret       string `yaml:"access_secret"`
	AccessTokenExpSec  int64  `yaml:"access_token_exp_sec"`
	RefreshSecret      string `yaml:"refresh_secret"`
	RefreshTokenExpSec int64  `yaml:"refresh_token_exp_sec"`
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
	if err != nil {
		return nil, err
	}
	defer f.Close()
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
