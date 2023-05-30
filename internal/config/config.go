package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Address       string   `toml:"address"`
	WidgetID      string   `toml:"widget_id"`
	EnableAvatars []string `toml:"enable_avatars"`
	VRChat        struct {
		Port int `toml:"port"`
	} `toml:"vrchat"`
	Logger struct {
		Level *int `toml:"level"`
	} `toml:"logger"`
}

func LoadFormFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if _, err = toml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
