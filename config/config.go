package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server  Server   `mapstructure:"server"`
	Runners []Runner `mapstructure:"runners"`
}

type Server struct {
	ListenAddress string `mapstructure:"listen_address"`
}

type Runner struct {
	Name   string            `mapstructure:"name"`
	Group  string            `mapstructure:"group"`
	Enable bool              `mapstructure:"enable"`
	Mode   string            `mapstructure:"mode"` // prod / test
	Labels map[string]string `mapstructure:"labels"`

	Logs struct {
		Runner string `mapstructure:"runner"`
		Worker string `mapstructure:"worker"`
		Event  string `mapstructure:"event"`
	} `mapstructure:"logs"`

	Test struct {
		RunnerPath string `mapstructure:"runner_path"`
		EventPath  string `mapstructure:"event_path"`
		WorkerPath string `mapstructure:"worker_path"`
	} `mapstructure:"test"`

	Metrics struct {
		EnableRunner bool `mapstructure:"enable_runner"`
		EnableJob    bool `mapstructure:"enable_job"`
		EnableEvent  bool `mapstructure:"enable_event"`
	} `mapstructure:"metrics"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("github-runner") // github-runner.yaml
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/github-runner-exporter/")
	v.AutomaticEnv()

	// Defaults
	v.SetDefault("server.listen_address", ":9200")
	v.SetDefault("mode", "prod")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}
