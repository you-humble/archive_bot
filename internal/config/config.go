package config

import (
	"flag"
	"archive_bot/pkg/er"
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel    string `yaml:"log_level"`
	IsWebhook   int    `yaml:"is_webhook"`
	AdminID     int64  `yaml:"admin_id"`
	PostgresURL string `yaml:"postgres_url"`
	Redis       Redis  `yaml:"redis"`
	Bot         Bot    `yaml:"bot"`
}

type Redis struct {
	Address    string `yaml:"address"`
	Password   string `yaml:"password"`
	User       string `yaml:"user"`
	Db         int    `yaml:"db"`
	MaxRetries int    `yaml:"max_retries"`
}

func (r Redis) Options() *redis.Options {
	return &redis.Options{
		Addr:       r.Address,
		Password:   r.Password,
		DB:         r.Db,
		Username:   r.User,
		MaxRetries: r.MaxRetries,
	}
}

type Bot struct {
	Token              string `yaml:"token"`
	WebhookURL         string `yaml:"webhook_url"`
	WebhookSecretToken string `yaml:"webhook_secret_token"`
	Port               string `yaml:"port"`
}

func New() (*Config, error) {
	const op = "config.New"

	configPath, err := parceFlags()
	if err != nil {
		return nil, err
	}
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, er.New("failed to open "+configPath, op, err)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, er.New(configPath+" failed to open the file", op, err)
	}

	return config, nil
}

func parceFlags() (string, error) {
	const op = "config.parceFlags"
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")

	flag.Parse()

	s, err := os.Stat(configPath)
	if err != nil {
		return "", er.New(configPath+" failed to get FileInfo", op, err)
	}
	if s.IsDir() {
		return "", er.New(configPath+" is a directory, must be a file", op, err)
	}

	return configPath, nil
}
