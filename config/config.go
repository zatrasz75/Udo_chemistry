package config

import "os"

type Config struct {
	Udo Udo
}

type Udo struct {
	AddrPort string
	AddrHost string
}

func New() *Config {
	return &Config{
		Udo: Udo{
			AddrPort: getEnv("APP_PORT", ""),
			AddrHost: getEnv("APP_HOST", ""),
		},
	}
}

// Простая вспомогательная функция для считывания окружения или возврата значения по умолчанию.
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
