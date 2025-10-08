package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
	GoogleAI GoogleAIConfig `mapstructure:"googleai"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
}

type OpenAIConfig struct {
	APIKey       string `mapstructure:"api_key"`
	BaseURL      string `mapstructure:"base_url"`
	Timeout      int    `mapstructure:"timeout"`
	MaxRetries   int    `mapstructure:"max_retries"`
	DefaultModel string `mapstructure:"default_model"`
}

type GoogleAIConfig struct {
	APIKey       string `mapstructure:"api_key"`
	ProjectID    string `mapstructure:"project_id"`
	Location     string `mapstructure:"location"`
	Timeout      int    `mapstructure:"timeout"`
	MaxRetries   int    `mapstructure:"max_retries"`
	DefaultModel string `mapstructure:"default_model"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 设置环境变量
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		// 注意：这里不能使用统一日志工具，因为日志器还未初始化
		// 使用fmt.Printf作为临时解决方案
		fmt.Printf("Warning: Could not read config file: %v\n", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	viper.SetDefault("database.driver", "sqlite3")
	viper.SetDefault("database.dsn", "./data/admin.db")

	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expire_time", 24)

	viper.SetDefault("openai.api_key", "")
	viper.SetDefault("openai.base_url", "https://api.openai.com/v1")
	viper.SetDefault("openai.timeout", 30)
	viper.SetDefault("openai.max_retries", 3)
	viper.SetDefault("openai.default_model", "gpt-3.5-turbo")

	viper.SetDefault("googleai.api_key", "")
	viper.SetDefault("googleai.project_id", "")
	viper.SetDefault("googleai.location", "us-central1")
	viper.SetDefault("googleai.timeout", 30)
	viper.SetDefault("googleai.max_retries", 3)
	viper.SetDefault("googleai.default_model", "gemini-1.5-flash")
}

func (c *Config) GetDatabaseDSN() string {
	return c.Database.DSN
}

func (c *Config) GetDatabaseDriver() string {
	return c.Database.Driver
}
