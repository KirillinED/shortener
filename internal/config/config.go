package config

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	DefaultHost            = "localhost"
	DefaultPort            = 80
	DefaultFileStoragePath = "/tmp/short-url-db.json"
)

type Config struct {
	Address
	BaseURL         string `env:"APP_BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() *Config {
	cfg := &Config{}

	cfg.Init()

	return cfg
}

func DefaultConfig() *Config {
	cfg := &Config{}

	SetDefault(cfg)

	return cfg
}

func (c *Config) Init() {
	ParseFlags(c)
	ParseEnv(c)
}

func SetDefault(c *Config) {
	c.Address = Address{}
	c.Address.Host = DefaultHost
	c.Address.Port = DefaultPort
	c.BaseURL = "http://" + c.Address.Host + ":" + strconv.Itoa(c.Address.Port) + "/"
	c.FileStoragePath = DefaultFileStoragePath
}

type Address struct {
	Host string `env:"APP_HOST"`
	Port int    `env:"APP_PORT"`
}

func (a *Address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a *Address) Set(flagValue string) error {
	a.Host = DefaultHost
	a.Port = DefaultPort

	hostPort := strings.Split(flagValue, ":")

	if len(hostPort[0]) != 0 {
		a.Host = hostPort[0]
	}

	if len(hostPort) > 1 && len(hostPort[1]) != 0 {
		port, err := strconv.Atoi(hostPort[1])
		if err != nil {
			return err
		}

		a.Port = port
	}

	return nil
}
