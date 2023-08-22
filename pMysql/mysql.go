package pMysql

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"net/url"
	"sync"
	"time"
)

var (
	dbOnce  sync.Once
	clients map[string]*gorm.DB
	err     error
)

func Connects() map[string]*gorm.DB {
	return clients
}

func Connect(key string) (*gorm.DB, error) {
	if clients == nil {
		return nil, errors.New("db connections not initialized")
	}
	conn, ok := clients[key]
	if !ok {
		return nil, errors.New("no such db connection defined")
	}
	return conn, nil
}

func NewClient(dbs map[string]*Database) (map[string]*gorm.DB, error) {
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
			var logLvl logger.LogLevel
			if v.LogLevel > 0 && v.LogLevel <= 4 {
				logLvl = logger.LogLevel(v.LogLevel)
			} else {
				logLvl = logger.Error
			}

			if v.SlowLog.String() != "" {
				config.SlowThreshold = v.SlowLog * time.Millisecond
			}
			if logLvl > 0 {
				config.LogLevel = logLvl
			}
		}

		//l := logger.New(log.StandardLogger(), config)
		//cfg.Logger = l

		conn, err := gorm.Open(dialectal, cfg)
		if err != nil {
			log.Error("Failed to create connection for db: ", k)
			continue
		}

		zapL, _ := zap.NewProduction()
		SetGormDBLogger(conn, New(zapL, WithCustomFields(func(ctx context.Context) zap.Field {
			item := ctx.Value("requestId")
			if vv, ok := item.(string); ok {
				return zap.String("requestId", vv)
			}
			return zap.Skip()
		},
		), WithConfig(config)))

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

func Client(dbs map[string]*Database) (map[string]*gorm.DB, error) {
	dbOnce.Do(func() {
		clients, err = NewClient(dbs)
	})
	return clients, err
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
