package pEtcd

import (
	"time"
)

type Config struct {
	Endpoints   []string      `toml:"endpoints" json:"endpoints,omitempty" yaml:"endpoints" mapstructure:"endpoints"`
	DialTimeout time.Duration `toml:"dial-timeout" json:"dial-timeout,omitempty" yaml:"dial-timeout" mapstructure:"dial-timeout"`
	Username    string        `toml:"username" json:"username,omitempty" yaml:"username" mapstructure:"username"`
	Password    string        `toml:"password" json:"password,omitempty" yaml:"password" mapstructure:"password"`
	Prefix      string        `toml:"prefix" json:"prefix,omitempty" yaml:"prefix" mapstructure:"prefix"`
}
