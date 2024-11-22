package dto

import (
	"database/sql"

	"github.com/EFG/api"
)

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
	Page            sql.NullInt32
	PageSize        sql.NullInt32
	FilterID        sql.NullString
	FilterFirstName sql.NullString
	FilterLastName  sql.NullString
	FilterNickname  sql.NullString
	FilterEmail     sql.NullString
	FilterCountry   sql.NullString
}

func (g *GetUsersArgs) FromAPI(req *api.GetUsersRequest) {
	g.Page = toNullInt32(req.Page)
	g.PageSize = toNullInt32(req.PageSize)
	g.FilterID = toNullString(req.FilterId)
	g.FilterFirstName = toNullString(req.FilterFirstName)
	g.FilterLastName = toNullString(req.FilterLastName)
	g.FilterNickname = toNullString(req.FilterNickname)
	g.FilterEmail = toNullString(req.FilterEmail)
	g.FilterCountry = toNullString(req.FilterCountry)
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func toNullInt32(i int32) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: i, Valid: true}
}
