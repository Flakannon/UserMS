package dto

import (
	"database/sql"

	"github.com/EFG/api"
	"github.com/EFG/internal/utils"
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
	g.Page = utils.ToNullInt32(req.Page)
	g.PageSize = utils.ToNullInt32(req.PageSize)
	g.FilterID = utils.ToNullString(req.FilterId)
	g.FilterFirstName = utils.ToNullString(req.FilterFirstName)
	g.FilterLastName = utils.ToNullString(req.FilterLastName)
	g.FilterNickname = utils.ToNullString(req.FilterNickname)
	g.FilterEmail = utils.ToNullString(req.FilterEmail)
	g.FilterCountry = utils.ToNullString(req.FilterCountry)
}
