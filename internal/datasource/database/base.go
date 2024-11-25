// Package database contains the common interfaces and structs for interacting with different database sources.
package database

import (
	"database/sql"
	"fmt"
)

type Client interface {
	PingDatabase() error
	Close() error
}

type BaseClient struct {
	DB *sql.DB
}

func (d BaseClient) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

func (d *BaseClient) PingDatabase() error {
	if err := d.DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}
