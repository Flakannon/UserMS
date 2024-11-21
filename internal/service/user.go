package service

import (
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

func (u *User) hashPassword() {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
	}
	u.Password = string(hashed)
}

func (u *User) FromAPI(apiUser *api.CreateUserRequest) {
	u.FirstName = apiUser.FirstName
	u.LastName = apiUser.LastName
	u.Nickname = apiUser.Nickname
	u.Password = apiUser.Password
	u.Email = apiUser.Email
	u.Country = apiUser.Country
}

func (u *User) toDTO() dto.UserDTO {
	return dto.UserDTO{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		Password:  u.Password,
		Email:     u.Email,
		Country:   u.Country,
	}
}
