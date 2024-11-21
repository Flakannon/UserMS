package user

import "context"

type UserReader interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
}

type UserWriter interface {
	CreateUser(ctx context.Context, user *User) error
	ModifyUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID int) error
}

type DataSource interface {
	UserReader
	UserWriter
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
}
