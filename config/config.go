package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
	"udo_mass/pkg/logger"
)

type Config struct {
	Server Server
}

type Server struct {
	AddrPort     string
	AddrHost     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	ShutdownTime time.Duration
}

func New() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	// Переменные для docker-compose
	addrPOrt := os.Getenv("APP_PORT")
	if addrPOrt == "" {
		addrPOrt = viper.GetString("APP_PORT")
	}
	addrHost := os.Getenv("APP_HOST")
	if addrHost == "" {
		addrHost = viper.GetString("APP_HOST")
	}

	readTimeoutStr := os.Getenv("READ_TIMEOUT")
	var readTimeout time.Duration
	if readTimeoutStr != "" {
		var err error
		readTimeout, err = time.ParseDuration(readTimeoutStr)
		if err != nil {
			logger.Error("ошибки парсинга времени", err)
		}
	} else {
		readTimeout, _ = time.ParseDuration(viper.GetString("READ_TIMEOUT"))
	}

	writeTimeoutStr := os.Getenv("WRITE_TIMEOUT")
	var writeTimeout time.Duration
	if writeTimeoutStr != "" {
		var err error
		writeTimeout, err = time.ParseDuration(writeTimeoutStr)
		if err != nil {
			logger.Error("ошибки парсинга времени", err)
		} else {
			writeTimeout, _ = time.ParseDuration(viper.GetString("WRITE_TIMEOUT"))
		}
	}

	idleTimeoutStr := os.Getenv("IDLE_TIMEOUT")
	var idleTimeout time.Duration
	if idleTimeoutStr != "" {
		var err error
		idleTimeout, err = time.ParseDuration(idleTimeoutStr)
		if err != nil {
			logger.Error("ошибки парсинга времени", err)
		} else {
			idleTimeout, _ = time.ParseDuration(viper.GetString("IDLE_TIMEOUT"))
		}
	}

	shutdownTimeStr := os.Getenv("SHUTDOWN_TIMEOUT")
	var shutdownTime time.Duration
	if shutdownTimeStr != "" {
		var err error
		shutdownTime, err = time.ParseDuration(shutdownTimeStr)
		if err != nil {
			logger.Error("ошибки парсинга времени", err)
		} else {
			shutdownTime, _ = time.ParseDuration(viper.GetString("SHUTDOWN_TIMEOUT"))
		}
	}

	return &Config{
		Server: Server{
			AddrPort:     addrPOrt,
			AddrHost:     addrHost,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
			ShutdownTime: shutdownTime,
		},
	}
}
