package dto

type UserDTO struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Country   string `json:"country"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UsersDTO []UserDTO
