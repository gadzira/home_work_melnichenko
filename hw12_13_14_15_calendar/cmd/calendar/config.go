package main

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf
	Server   server
	DataBase database
	Mode     mode
}

type LoggerConf struct {
	Level      string
	LogFile    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LocalTime  bool
	Compress   bool
}

type server struct {
	IP   string
	Port string
}

type mode struct {
	M string
}

type database struct {
	DSN string
}

func NewConfig(fileName string) Config {
	var confDir = "../../configs" // nolint:gofumpt
	conFile := filepath.Join(confDir, fileName)

	var config Config
	if _, err := toml.DecodeFile(conFile, &config); err != nil {
		log.Fatal("Can't load configuration file:", err)
	}

	return Config{
		Server:   config.Server,
		Logger:   config.Logger,
		DataBase: config.DataBase,
		Mode:     config.Mode,
	}
}
