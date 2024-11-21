package dto

import "database/sql"

type UserDTO struct {
	ID        sql.NullString `json:"id"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	Nickname  sql.NullString `json:"nickname"`
	Password  sql.NullString `json:"password"`
	Email     sql.NullString `json:"email"`
	Country   sql.NullString `json:"country"`
	CreatedAt sql.NullTime   `json:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
}

type UsersDTO []UserDTO

type GetUsersArgs struct {
	Page            int32
	PageSize        int32
	FilterID        string
	FilterFirstName string
	FilterLastName  string
	FilterNickname  string
	FilterEmail     string
	FilterCountry   string
}
