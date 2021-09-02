package db

import (
	"database/sql"
	"errors"
	"fmt"
	"smsc/pkg/log"
	"time"

	//github.com/jackc/pgx/v4/stdlib
	_ "github.com/jackc/pgx/v4/stdlib"
	//github.com/go-sql-driver/mysql
	// _ "github.com/go-sql-driver/mysql"
)

//Config database struct
type Config struct {
	Driver          string
	Host            string
	Port            string
	Dbname          string
	SslMode         string
	User            string
	Pass            string
	ConnMaxLifetime int
	MaxOpenConns    int
	MaxIdleConns    int
	ApplicationName string
}

var cfg *Config

var db *sql.DB
var err error

//SetConfigDB func
func SetConfigDB(conf *Config) {
	cfg = conf
}

//Init func
func Init() {

	if cfg == nil {
		log.Info("config is nil Error config Database not set", cfg)
		panic(errors.New("config is nil"))
	}

	dbConnString := ""
	if cfg.Driver == "pgx" {
		dbConnString = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s application_name=%s",
			cfg.User,
			cfg.Pass,
			cfg.Host,
			cfg.Port,
			cfg.Dbname,
			cfg.SslMode,
			cfg.ApplicationName,
		)
	} else {
		dbConnString = fmt.Sprintf("%s:%s@/%s", cfg.User, cfg.Pass, cfg.Dbname)
	}

	db, err = sql.Open(cfg.Driver, dbConnString)

	if err != nil {
		log.Error("Failed to connect to database >> ", dbConnString, err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Millisecond * time.Duration(cfg.ConnMaxLifetime))

}

//GetDB - get DB
func GetDB() *sql.DB {
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
