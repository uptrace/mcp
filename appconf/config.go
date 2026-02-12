package appconf

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Service ServiceConfig `yaml:"service"`
	Logging LoggingConfig `yaml:"logging"`
	Default DefaultConfig `yaml:"default"`
	Uptrace UptraceConfig `yaml:"uptrace"`
}
type DefaultConfig struct {
	Limit int    `yaml:"limit"`
	Query string `yaml:"query"`
}
type ServiceConfig struct {
	StartTimeout time.Duration `yaml:"start_timeout"`
	StopTimeout  time.Duration `yaml:"stop_timeout"`
}

type LoggingConfig struct {
	Level       string `yaml:"level"`
	MaxBodySize int    `yaml:"max_body_size"`
}

type UptraceConfig struct {
	DSN       string `yaml:"dsn"`
	APIURL    string `yaml:"api_url"`
	APIToken  string `yaml:"api_token"`
	ProjectID int64  `yaml:"project_id"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	return Parse(data)
}

// Parse parses YAML data into a Config struct.
func Parse(data []byte) (*Config, error) {
	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	conf.setDefaults()
	return &conf, nil
}

func (c *Config) setDefaults() {
	if c.Service.StartTimeout == 0 {
		c.Service.StartTimeout = 15 * time.Second
	}
	if c.Service.StopTimeout == 0 {
		c.Service.StopTimeout = 15 * time.Second
	}
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
}
