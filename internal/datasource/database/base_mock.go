package database

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

func MockRowsToSQLRows(mockRows *sqlmock.Rows) *sql.Rows {
	db, mock, _ := sqlmock.New()
	mock.ExpectQuery("select").WillReturnRows(mockRows)
	rows, _ := db.Query("select")
	return rows
}

func MockEmptyResultSet(cols []string) *sql.Rows {
	rows := sqlmock.NewRows(cols)
	return MockRowsToSQLRows(rows)
}
