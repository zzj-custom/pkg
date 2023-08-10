package pEtcd

import (
	"time"
)

type Config struct {
	Endpoints   []string      `toml:"endpoints"`
	DialTimeout time.Duration `toml:"dial_timeout"`
	Username    string        `toml:"username"`
	Password    string        `toml:"password"`
	Prefix      string        `toml:"prefix"`
}
