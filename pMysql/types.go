package pMysql

import (
	"time"
)

type Database struct {
	DSN             string        `json:"dsn" toml:"dsn"`
	Username        string        `json:"username" toml:"username"`
	Password        string        `json:"password" toml:"password"`
	Host            string        `json:"host" toml:"host"`
	Database        string        `json:"database" toml:"database"`
	MaxOpenConn     int           `json:"max_open_conn" toml:"max_open_conn"`
	MaxIdleConn     int           `json:"max_idle_conn" toml:"max_idle_conn"`
	ConnMaxIdleTime int           `json:"conn_max_idletime" toml:"conn_max_idletime"`
	ConnMaxLifetime int           `json:"conn_max_lifetime" toml:"conn_max_lifetime"`
	UseLog          bool          `json:"use_log" toml:"use_log"`
	LogLevel        int           `json:"log_level" toml:"log_level"`
	SlowLog         time.Duration `json:"slow_log" toml:"slow_log"`
	TablePrefix     string        `json:"table_prefix" toml:"table_prefix""`
	Singular        bool          `json:"singular" toml:"singular"`
}
