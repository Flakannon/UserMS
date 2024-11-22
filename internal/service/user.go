package service

import (
	"time"

	"github.com/EFG/api"
	"github.com/EFG/internal/datasource/dto"
	"github.com/EFG/internal/utils"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) hashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)

	return nil
}

func NewUserFromCreateRequest(req *api.CreateUserRequest) User {
	if req == nil {
		return User{}
	}
	return User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
	}
}

func NewUserFromModifyRequest(req *api.ModifyUserRequest) User {
	if req == nil {
		return User{}
	}
	return User{
		ID:        req.Id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
	}
}

func (u *User) toDTO() dto.UserDTO {
	return dto.UserDTO{
		ID:        utils.ToNullString(u.ID),
		FirstName: utils.ToNullString(u.FirstName),
		LastName:  utils.ToNullString(u.LastName),
		Nickname:  utils.ToNullString(u.Nickname),
		Password:  utils.ToNullString(u.Password),
		Email:     utils.ToNullString(u.Email),
		Country:   utils.ToNullString(u.Country),
		CreatedAt: utils.ToNullTime(u.CreatedAt),
		UpdatedAt: utils.ToNullTime(u.UpdatedAt),
	}
}

// func NewUserFromDTO(userDTO dto.UserDTO) User {
// 	return User{
// 		ID:        userDTO.ID.String,
// 		FirstName: userDTO.FirstName.String,
// 		LastName:  userDTO.LastName.String,
// 		Nickname:  userDTO.Nickname.String,
// 		Password:  userDTO.Password.String,
// 		Email:     userDTO.Email.String,
// 		Country:   userDTO.Country.String,
// 		CreatedAt: userDTO.CreatedAt.Time,
// 		UpdatedAt: userDTO.UpdatedAt.Time,
// 	}
// }

func FromDTOToAPI(userDTO dto.UsersDTO) []*api.User {
	users := make([]*api.User, len(userDTO))
	for i, u := range userDTO {
		users[i] = &api.User{
			Id:        u.ID.String,
			FirstName: u.FirstName.String,
			LastName:  u.LastName.String,
			Nickname:  u.Nickname.String,
			Email:     u.Email.String,
			Country:   u.Country.String,
			CreatedAt: u.CreatedAt.Time.String(),
			UpdatedAt: u.UpdatedAt.Time.String(),
		}
	}
	return users
}
