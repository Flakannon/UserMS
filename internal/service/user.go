package service

import (
	"database/sql"
	"time"

	"github.com/EFG/api"
	"github.com/EFG/internal/datasource/dto"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type Datasource interface {
	Reader
	Writer
}

type User struct {
	ID        string
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Users []User

func (u *User) hashPassword() {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
	}
	u.Password = string(hashed)
}

func (u *User) FromAPICreate(apiUser *api.CreateUserRequest) {
	u.FirstName = apiUser.FirstName
	u.LastName = apiUser.LastName
	u.Nickname = apiUser.Nickname
	u.Password = apiUser.Password
	u.Email = apiUser.Email
	u.Country = apiUser.Country
}

func (u *User) FromAPIUpdate(apiUser *api.ModifyUserRequest) {
	u.ID = apiUser.Id
	u.FirstName = apiUser.FirstName
	u.LastName = apiUser.LastName
	u.Nickname = apiUser.Nickname
	u.Password = apiUser.Password
	u.Email = apiUser.Email
	u.Country = apiUser.Country
}

func (u *User) toDTO() dto.UserDTO {
	return dto.UserDTO{
		ID:        toNullString(u.ID),
		FirstName: toNullString(u.FirstName),
		LastName:  toNullString(u.LastName),
		Nickname:  toNullString(u.Nickname),
		Password:  toNullString(u.Password),
		Email:     toNullString(u.Email),
		Country:   toNullString(u.Country),
		CreatedAt: sql.NullTime{Time: u.CreatedAt, Valid: !u.CreatedAt.IsZero()},
		UpdatedAt: sql.NullTime{Time: u.UpdatedAt, Valid: !u.UpdatedAt.IsZero()},
	}
}

func (u *User) FromDTO(userDTO dto.UserDTO) {
	u.ID = userDTO.ID.String
	u.FirstName = userDTO.FirstName.String
	u.LastName = userDTO.LastName.String
	u.Nickname = userDTO.Nickname.String
	u.Password = userDTO.Password.String
	u.Email = userDTO.Email.String
	u.Country = userDTO.Country.String
	u.CreatedAt = userDTO.CreatedAt.Time
	u.UpdatedAt = userDTO.UpdatedAt.Time
}

func FromDTOToAPI(userDTO dto.UsersDTO) []*api.User {
	users := make([]*api.User, 0, len(userDTO))
	for _, u := range userDTO {
		users = append(users, &api.User{
			Id:        u.ID.String,
			FirstName: u.FirstName.String,
			LastName:  u.LastName.String,
			Nickname:  u.Nickname.String,
			Email:     u.Email.String,
			Country:   u.Country.String,
		})
	}

	return users
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
