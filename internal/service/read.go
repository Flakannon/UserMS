package service

import "context"

type Reader interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
}
