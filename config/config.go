package config

import (
	"log"
	"os"

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
		APIKey:         "sk-gr-8a67f6189efe32e3046174a1d1a7e4c397da4ce7",
		BaseURL:        "https://endpoint.greatrouter.com",
		Model:          "gpt-5.4-mini",
		Port:           "8080",
		DBDSN:          os.Getenv("DB_DSN"),
		JWTSecret:      "8c16e4203f2d674409009hhyb56a5bf",
		WeChatAppID:    "wx69dffb6e777c0b96",
		WeChatSecret:   "8c16e4203f2d6744195e33315b56a5bf",
		WeChatMockMode: false,
	}
	log.Printf("config %+v", cfg)

	return cfg
}
