package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port string
	DSN  string
}

//"dsn": "file:resources\\universities.db?cache=shared&mode=memory"

func GetConfig() (*Config, error) {
	configContent, err := os.ReadFile("resources/config.json")
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	// cfg.Port = "localhost:8080"
	// cfg.DSN = "file:resources\\universities.db"
	err = json.Unmarshal(configContent, cfg)
	fmt.Println(*cfg)
	return cfg, err
}
