package postgres

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/EFG/internal/datasource/database"
	"github.com/EFG/internal/env"

	_ "github.com/lib/pq"
)

type Client struct {
	Config env.DatabaseConfig
	*database.BaseClient
}

func NewClient(config env.DatabaseConfig) *Client {
	return &Client{
		Config:     config,
		BaseClient: &database.BaseClient{},
	}
}

// Connects to database via postgres driver
func (d *Client) Connect() error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
		d.Config.Username, d.Config.Password, d.Config.Host, d.Config.Port, d.Config.Database)

	slog.Info("Now attempting to connect to postgres database")

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Println(err)
		return err
	}
	db.SetMaxOpenConns(80)
	db.SetMaxIdleConns(15)
	db.SetConnMaxLifetime(30 * time.Minute)

	d.DB = db

	err = d.PingDatabase()
	if err != nil {
		return fmt.Errorf("error pinging postgres database: %v", err)
	}

	slog.Info("Successfully connected to postgres database")

	return nil
}
