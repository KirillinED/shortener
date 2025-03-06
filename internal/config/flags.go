package config

import (
	"flag"
	"strconv"
)

func ParseFlags(c *Config) {
	flag.Var(&c.Address, "a", "Address to listen on host:port")
	flag.StringVar(&c.BaseURL, "b", "http://"+DefaultHost+":"+strconv.Itoa(DefaultPort)+"/", "Base URL for short link")
	flag.StringVar(&c.FileStoragePath, "f", DefaultFileStoragePath, "Storage file path")
	flag.StringVar(&c.DatabaseDSN, "dsn", "", "Database DSN")
	flag.StringVar(&c.DatabaseDriver, "d", "", "Database driver")
	flag.StringVar(&c.CookieHashKey, "cookie-hash-key", "", "Cookie hash key")
	flag.Parse()
}
