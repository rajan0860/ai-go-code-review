package config

import (
	"os"
	"strconv"
)

type Config struct {
	OllamaModel      string
	QualityThreshold float64
	EnableSecurity   bool
}

func LoadConfig() *Config {
	cfg := &Config{
		OllamaModel:      "qwen2.5-coder:latest",
		QualityThreshold: 7.5,
		EnableSecurity:   true,
	}

	if model := os.Getenv("OLLAMA_MODEL"); model != "" {
		cfg.OllamaModel = model
	}

	if threshold := os.Getenv("QUALITY_THRESHOLD"); threshold != "" {
		if val, err := strconv.ParseFloat(threshold, 64); err == nil {
			cfg.QualityThreshold = val
		}
	}

	return cfg
}
