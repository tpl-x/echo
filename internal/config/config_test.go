package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestCreateAConfig(t *testing.T) {
	cfg := AppConfig{
		Server: ServerConfig{
			BindPort:         9000,
			GraceExitTimeout: 5,
		},
		Log: LogConfig{
			FileName:    "server.log",
			MaxSize:     10,
			MaxBackups:  7,
			MaxKeepDays: 20,
			Compress:    true,
		},
	}
	f, err := os.OpenFile("testdata/config_example.yaml", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatal(err)
	}
	if err = yaml.NewEncoder(f).Encode(cfg); err != nil {
		t.Fatal(err)
	}
	t.Log("config created")
}
