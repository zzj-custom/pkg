package pRedis

import (
	"time"
)

type DialConfig struct {
	Host            string        `toml:"host" json:"host,omitempty" yaml:"host" mapstructure:"host"`
	Port            int           `toml:"port" json:"port,omitempty" yaml:"port" mapstructure:"port"`
	Database        int           `toml:"database" json:"database,omitempty" yaml:"database" mapstructure:"database"`
	Password        string        `toml:"password" json:"password,omitempty" yaml:"password" mapstructure:"password"`
	MaxIdle         int           `toml:"max-idle" json:"max-idle,omitempty" yaml:"max-idle" mapstructure:"max-idle"`
	MaxActive       int           `toml:"max-active" json:"max-active,omitempty" yaml:"max-active" mapstructure:"max-active"`
	Wait            bool          `toml:"wait" json:"wait,omitempty" yaml:"wait" mapstructure:"wait"`
	ConnectTimeout  time.Duration `toml:"connect-timeout" json:"connect-timeout,omitempty" yaml:"connect-timeout" mapstructure:"connect-timeout"`
	ReadTimeout     time.Duration `toml:"read-timeout" json:"read-timeout,omitempty" yaml:"read-timeout" mapstructure:"read-timeout"`
	MaxConnLifetime time.Duration `toml:"max-conn-lifetime" json:"max-conn-lifetime,omitempty" yaml:"max-conn-lifetime" mapstructure:"max-conn-lifetime"`
	IdleTimeout     time.Duration `toml:"idle-timeout" json:"idle-timeout,omitempty" yaml:"idle-timeout" mapstructure:"idle-timeout"`
}

type MultiDialConfig struct {
	Name    string      `toml:"name"`
	Default bool        `toml:"default"`
	Config  *DialConfig `toml:"config"`
}
