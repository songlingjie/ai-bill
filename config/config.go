package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey         string
	BaseURL        string
	Model          string
	Port           string
	DBDSN          string
	JWTSecret      string
	WeChatAppID    string
	WeChatSecret   string
	WeChatMockMode bool
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		APIKey:         os.Getenv("OPENAI_API_KEY"),
		BaseURL:        os.Getenv("OPENAI_BASE_URL"),
		Model:          os.Getenv("OPENAI_MODEL"),
		Port:           os.Getenv("PORT"),
		DBDSN:          os.Getenv("DBDSN"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		WeChatAppID:    os.Getenv("WECHAT_APP_ID"),
		WeChatSecret:   os.Getenv("WECHAT_APP_SECRET"),
		WeChatMockMode: strings.EqualFold(os.Getenv("WECHAT_MOCK_MODE"), "true"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.holysheep.ai/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "gpt-4.1-mini"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "replace-me-with-a-long-random-secret"
	}

	return cfg
}
