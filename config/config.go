package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
	"udo_mass/pkg/logger"
)

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	AddrPort     string
	AddrHost     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	ShutdownTime time.Duration
}

type Database struct {
	ConnStr string //postgres://postgres:postgrespw@localhost:49153/Account

	Host     string // postgres
	User     string // postgres
	Password string // postgrespw
	Name     string // Account
	Port     string // 49153
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
		Database: Database{
			ConnStr:  initDb(),
			Host:     viper.GetString("HOST_DB"),
			User:     viper.GetString("USER_DB"),
			Password: viper.GetString("PASSWORD_DB"),
			Name:     viper.GetString("NAME_DB"),
			Port:     viper.GetString("PORT_DB"),
		},
	}
}

func initDb() string {
	connStr := os.Getenv("DB_POSTGRES_URL")
	if connStr == "" {
		c := &Config{
			Database: Database{
				Host:     viper.GetString("HOST_DB"),
				User:     viper.GetString("USER_DB"),
				Password: viper.GetString("PASSWORD_DB"),
				Name:     viper.GetString("NAME_DB"),
				Port:     viper.GetString("PORT_DB"),
			},
		}
		connStr = fmt.Sprintf(
			"%s://%s:%s@localhost:%s/%s",
			c.Database.Host, c.Database.User, c.Database.Password, c.Database.Port, c.Database.Name,
		)
	}
	return connStr
}
