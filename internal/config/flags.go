package config

import (
	"flag"
	"strconv"
)

func ParseFlags(c *config) {
	flag.Var(&c.Address, "a", "Address to listen on host:port")
	flag.StringVar(&c.BaseURL, "b", "http://"+DefaultHost+":"+strconv.Itoa(DefaultPort)+"/", "Base URL for short link")

	flag.Parse()
}
