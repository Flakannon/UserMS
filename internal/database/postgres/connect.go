package postgres

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	"github.com/EFG/internal/database"
	"github.com/EFG/internal/env"
)

type PostgresConnector struct {
	Config env.DatabaseConfig
	*database.BaseClient
}

// Connects to database via postgres driver
func (d *PostgresConnector) Connect() error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.Config.Username, d.Config.Password, d.Config.Host, d.Config.Port, d.Config.Database)
	log.Println("Now attempting to connect to postgres database")

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Println(err)
		return err
	}

	d.DB = db

	err = d.PingDatabase()
	if err != nil {
		return fmt.Errorf("error pinging postgres database: %v", err)
	}

	log.Println("Successfully connected to postgres database")

	return nil
}
