package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func ConnectToPostgres(host string, port int, username string, password string, dbname string, logger *zap.SugaredLogger) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		username, password, host, port, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.New("Failed to connect to database: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New("Failed to ping database: " + err.Error())
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)
	return db, nil
}
