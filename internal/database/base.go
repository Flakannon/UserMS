// Package database contains the interfaces and structs for interacting with different data sources.
package database

import (
	"database/sql"
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

func (d BaseClient) PingDatabase() error {
	return d.DB.Ping()
}