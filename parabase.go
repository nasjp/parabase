package parabase

import (
	"database/sql"
	"testing"
	"time"
)

type Config struct {
	ParaNum            int
	Timeout            time.Duration
	DriverName         string
	DataSourceName     string
	ManagementDatabase ManagementDatabase
}

type ManagementDatabase interface {
	Connect(dbName string, cfg *Config) (*sql.DB, error)
	Setup(*sql.DB, *Config) error
	Get(db *sql.DB, cfg *Config) (*sql.DB, func(testing.TB), error)
}

func Use(cfg *Config) (*sql.DB, func(t testing.TB), error) {
	db, err := cfg.ManagementDatabase.Connect("", cfg)
	if err != nil {
		return nil, nil, err
	}

	if err := cfg.ManagementDatabase.Setup(db, cfg); err != nil {
		return nil, nil, err
	}

	db, teardown, err := cfg.ManagementDatabase.Get(db, cfg)
	if err != nil {
		return nil, nil, err
	}

	return db, teardown, nil
}
