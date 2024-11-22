package integrationtest

import (
	"database/sql"
	"fmt"

	"github.com/EFG/internal/datasource/dto"
	_ "github.com/lib/pq"
)

type PostgresClient struct {
	DB *sql.DB
}

func (p *PostgresClient) Connect() error {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	p.DB = db

	if err := p.DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (p *PostgresClient) Disconnect() error {
	return p.DB.Close()
}

func (p *PostgresClient) GetUserCount() (int, error) {
	query := "SELECT COUNT(*) FROM users"
	row := p.DB.QueryRow(query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

func (p *PostgresClient) GetUserById(id string) (dto.UserDTO, error) {
	query := `SELECT id, first_name, last_name, email, password, country, nick_name, created_at, updated_at
              FROM users WHERE id = $1`
	row := p.DB.QueryRow(query, id)

	var user dto.UserDTO
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Country,
		&user.Nickname,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.UserDTO{}, fmt.Errorf("user with id %s not found", id)
		}
		return dto.UserDTO{}, fmt.Errorf("query failed: %w", err)
	}

	return user, nil
}

func (p *PostgresClient) ResetUserStore() error {
	_, err := p.DB.Exec("TRUNCATE TABLE users")
	if err != nil {
		return fmt.Errorf("failed to truncate users table: %w", err)
	}

	return nil
}
