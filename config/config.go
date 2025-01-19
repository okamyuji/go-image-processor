package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds the configuration values
type Config struct {
	DefaultWidth  int `yaml:"default_width"`
	DefaultHeight int `yaml:"default_height"`
	DefaultAngle  int `yaml:"default_angle"`
	JpegQuality   int `yaml:"jpeg_quality"`
}

// LoadConfig reads the config file and returns a Config struct
func LoadConfig(filename string) (*Config, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetConfig loads the configuration or returns default values
func GetConfig() *Config {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		slog.Warn("error loading config file, using default values",
			"error", err)
		return &Config{
			DefaultWidth:  800,
			DefaultHeight: 600,
			DefaultAngle:  90,
			JpegQuality:   75,
		}
	}
	return config
}
