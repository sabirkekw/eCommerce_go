package postgres

import (
	"database/sql"
	"fmt"

	"github.com/sabirkekw/ecommerce_go/order-service/internal/cfg"
)

func ConnectToPostgres(cfg *cfg.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db, nil
}
