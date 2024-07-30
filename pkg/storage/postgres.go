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

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	fmt.Println("Vales for db:")
	fmt.Printf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect(os.Getenv(storTypeEnvVar), fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
