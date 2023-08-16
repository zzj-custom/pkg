package pMysql

import (
	"time"
)

type Database struct {
	DSN             string        `json:"dsn" toml:"dsn" yaml:"DSN" mapstructure:"DSN"`
	Username        string        `json:"username" toml:"username" yaml:"username" mapstructure:"username"`
	Password        string        `json:"password" toml:"password" yaml:"password" mapstructure:"password"`
	Host            string        `json:"host" toml:"host" yaml:"host" mapstructure:"host"`
	Database        string        `json:"database" toml:"database" yaml:"database" mapstructure:"database"`
	MaxOpenConn     int           `json:"max-open-conn" toml:"max-open-conn" yaml:"max-open-conn" mapstructure:"max-open-conn"`
	MaxIdleConn     int           `json:"max-idle-conn" toml:"max-idle-conn" yaml:"max-idle-conn" mapstructure:"max-idle-conn"`
	ConnMaxIdleTime int           `json:"conn-max-idletime" toml:"conn-max-idletime" yaml:"conn-max-idletime" mapstructure:"conn-max-idletime"`
	ConnMaxLifetime int           `json:"conn-max-lifetime" toml:"conn-max-lifetime" yaml:"conn-max-lifetime" mapstructure:"conn-max-lifetime"`
	UseLog          bool          `json:"use-log" toml:"use-log" yaml:"use-log" mapstructure:"use-log"`
	LogLevel        int           `json:"log-level" toml:"log-level" yaml:"log-level" mapstructure:"log-level"`
	SlowLog         time.Duration `json:"slow-log" toml:"slow-log" yaml:"slow-log" mapstructure:"slow-log"`
	TablePrefix     string        `json:"table-prefix" toml:"table-prefix" yaml:"table-prefix" mapstructure:"table-prefix"`
	Singular        bool          `json:"singular" toml:"singular" yaml:"singular" mapstructure:"singular"`
}
