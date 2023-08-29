package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Udo struct {
	AddrPort string
	AddrHost string
}

type (
	Config struct {
		Server `mapstructure:",squash"`
	}
	Server struct {
		AddrPort string `mapstructure:"APP_PORT"`
		AddrHost string `mapstructure:"APP_HOST"`
	}
)

func (c *Config) LoadEnv(path string) error {

	if path == "" {
		files, err := os.ReadDir(".")
		if err != nil {
			return errors.Wrap(err, "не удалось найти конфигурацию")
		}

		for _, file := range files {
			file.Info()

			filename := file.Name()

			if ext := filepath.Ext(filename); ext != ".env" {
				continue
			}

			if err := c.load("./" + filename); err != nil {
				return errors.Wrap(err, "не удалось загрузить конфигурацию")
			}
			return nil
		}

	}

	if err := c.load(path); err != nil {
		return errors.Wrap(err, "не удалось загрузить конфигурацию")
	}

	return nil
}

func (c *Config) load(path string) error {
	dir, file := filepath.Split(path)
	filename := filepath.Base(path)
	ext := filepath.Ext(file)
	name := filename[0 : len(filename)-len(ext)]

	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(name)
	v.SetConfigType(ext[1:])
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "не удалось прочитать конфигурационный файл")
	}

	err = v.Unmarshal(&c)
	if err != nil {
		return errors.Wrap(err, "не удалось расшифровать конфигурацию для struct")
	}

	return nil
}
