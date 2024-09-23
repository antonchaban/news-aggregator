package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

const storTypeEnvVar = "STORAGE_TYPE"

type Config struct {
	Host     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(os.Getenv(storTypeEnvVar), fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s",
		os.Getenv(storTypeEnvVar), cfg.Username, cfg.Password, cfg.Host, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
