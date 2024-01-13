package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type AppConfig struct {
	Server ServerConfig
	Log    LogConfig
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
