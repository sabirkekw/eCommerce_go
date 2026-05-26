package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sabirkekw/ecommerce_go/cart-service/internal/config"
)

func ConnectToPostgres(cfg *config.Config) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic("failed to connect to postgres")
	}
	if err := db.Ping(); err != nil {
		panic("failed to ping postgres")
	}
	return db
}
