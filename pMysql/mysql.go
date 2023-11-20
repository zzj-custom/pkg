package pMysql

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"net/url"
	"sync"
	"time"
)

var (
	dbOnce    sync.Once
	dbManager map[string]*gorm.DB
	err       error
)

func Dbs() map[string]*gorm.DB {
	return dbManager
}

func GetDb(key string) (*gorm.DB, error) {
	if dbManager == nil {
		return nil, errors.New("db connections not initialized")
	}
	conn, ok := dbManager[key]
	if !ok {
		return nil, errors.New("no such db connection defined")
	}
	return conn, nil
}

func Setup(dbs map[string]*Database) (map[string]*gorm.DB, error) {
	connects := map[string]*gorm.DB{}
	for k, v := range dbs {
		var dsn = buildDSN(v)
		dialectal := gormMySQL.New(gormMySQL.Config{
			DSN:                       dsn,
			SkipInitializeWithVersion: false,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			DontSupportForShareClause: true,
		})
		cfg := &gorm.Config{
			DisableAutomaticPing: false,
		}
		if v.TablePrefix != "" {
			cfg.NamingStrategy = schema.NamingStrategy{TablePrefix: v.TablePrefix, SingularTable: v.Singular}
		}

		config := logger.Config{
			SlowThreshold:             0,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  0,
		}
		if v.UseLog {
			if v.LogLevel > 0 && v.LogLevel <= 4 {
				config.LogLevel = logger.LogLevel(v.LogLevel)
			} else {
				config.LogLevel = logger.Error
			}

			if v.SlowLog.String() != "" {
				config.SlowThreshold = v.SlowLog * time.Millisecond
			}
		}

		conn, err := gorm.Open(dialectal, cfg)
		if err != nil {
			log.Error("Failed to create connection for db: ", k)
			continue
		}

		sqlDB, err := conn.DB()
		if err != nil {
			log.Error("Failed to Connect to db server.")
			continue
		}
		sqlDB.SetMaxIdleConns(v.MaxIdleConn)
		sqlDB.SetMaxOpenConns(v.MaxOpenConn)
		sqlDB.SetConnMaxIdleTime(time.Duration(v.ConnMaxIdleTime) * time.Second)
		sqlDB.SetConnMaxLifetime(time.Duration(v.ConnMaxLifetime) * time.Second)

		connects[k] = conn
	}

	return connects, nil
}

func DbInit(dbs map[string]*Database) (map[string]*gorm.DB, error) {
	dbOnce.Do(func() {
		dbManager, err = Setup(dbs)
	})
	return dbManager, err
}

func buildDSN(v *Database) string {
	var dsn string
	if v.DSN != "" {
		dsn = v.DSN
	} else {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true",
			url.QueryEscape(v.Username),
			url.QueryEscape(v.Password),
			v.Host,
			v.Database,
		)
	}
	return dsn
}
