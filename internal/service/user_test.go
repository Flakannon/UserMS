package service

import (
	"reflect"
	"testing"
	"time"

	"github.com/EFG/api"
	"github.com/EFG/internal/datasource/dto"
	"github.com/EFG/internal/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUser_hashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "valid password",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "empty password",
			password:    "",
			expectError: false, // bcrypt should hash an empty string without error
		},
		{
			name:        "password with special characters",
			password:    "p4ssw0rd!@#",
			expectError: false,
		},
		{
			name:        "password with spaces",
			password:    "password 123",
			expectError: false,
		},
		{
			name:        "password with spaces and special characters",
			password:    "p4ssw0rd!@# 123",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Password: tt.password,
			}

			err := u.hashPassword()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.NotEqual(t, tt.password, u.Password, "hashed password should be the same as the original password")

				err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(tt.password))
				assert.NoError(t, err, "plaintext equivalent of hash should match the original password")
			}
		})
	}
}

func TestNewUserFromCreateRequest(t *testing.T) {
	type args struct {
		req *api.CreateUserRequest
	}
	tests := []struct {
		name string
		args args
		want User
	}{
		{
			name: "valid request",
			args: args{
				req: &api.CreateUserRequest{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					Password:  "password123",
					Country:   "USA",
					Nickname:  "johndoe",
				},
			},
			want: User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Country:   "USA",
				Nickname:  "johndoe",
			},
		},
		{
			name: "empty request",
			args: args{
				req: &api.CreateUserRequest{},
			},
			want: User{
				FirstName: "",
				LastName:  "",
				Email:     "",
				Password:  "",
				Country:   "",
				Nickname:  "",
			},
		},
		{
			name: "nil request",
			args: args{
				req: nil,
			},
			want: User{
				FirstName: "",
				LastName:  "",
				Email:     "",
				Password:  "",
				Country:   "",
				Nickname:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserFromCreateRequest(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserFromCreateRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserFromModifyRequest(t *testing.T) {
	type args struct {
		req *api.ModifyUserRequest
	}
	tests := []struct {
		name string
		args args
		want User
	}{
		{
			name: "valid request",
			args: args{
				req: &api.ModifyUserRequest{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					Password:  "password123",
					Country:   "USA",
					Nickname:  "johndoe",
				},
			},
			want: User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Country:   "USA",
				Nickname:  "johndoe",
			},
		},
		{
			name: "empty request",
			args: args{
				req: &api.ModifyUserRequest{},
			},
			want: User{
				FirstName: "",
				LastName:  "",
				Email:     "",
				Password:  "",
				Country:   "",
				Nickname:  "",
			},
		},
		{
			name: "nil request",
			args: args{
				req: nil,
			},
			want: User{
				FirstName: "",
				LastName:  "",
				Email:     "",
				Password:  "",
				Country:   "",
				Nickname:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserFromModifyRequest(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserFromModifyRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_toDTO(t *testing.T) {
	timestamp := time.Now()
	tests := []struct {
		name string
		user User
		want dto.UserDTO
	}{
		{
			name: "valid user",
			user: User{
				ID:        "123",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Password:  "password123",
				Email:     "john.doe@email.com",
				Country:   "USA",
				CreatedAt: timestamp,
				UpdatedAt: timestamp,
			},
			want: dto.UserDTO{
				ID:        utils.ToNullString("123"),
				FirstName: utils.ToNullString("John"),
				LastName:  utils.ToNullString("Doe"),
				Nickname:  utils.ToNullString("johndoe"),
				Password:  utils.ToNullString("password123"),
				Email:     utils.ToNullString("john.doe@email.com"),
				Country:   utils.ToNullString("USA"),
				CreatedAt: utils.ToNullTime(timestamp),
				UpdatedAt: utils.ToNullTime(timestamp),
			},
		},
		{
			name: "empty user",
			user: User{},
			want: dto.UserDTO{
				ID:        utils.ToNullString(""),
				FirstName: utils.ToNullString(""),
				LastName:  utils.ToNullString(""),
				Nickname:  utils.ToNullString(""),
				Password:  utils.ToNullString(""),
				Email:     utils.ToNullString(""),
				Country:   utils.ToNullString(""),
				CreatedAt: utils.ToNullTime(time.Time{}),
				UpdatedAt: utils.ToNullTime(time.Time{}),
			},
		},
		{
			name: "nil user",
			user: User{},
			want: dto.UserDTO{
				ID:        utils.ToNullString(""),
				FirstName: utils.ToNullString(""),
				LastName:  utils.ToNullString(""),
				Nickname:  utils.ToNullString(""),
				Password:  utils.ToNullString(""),
				Email:     utils.ToNullString(""),
				Country:   utils.ToNullString(""),
				CreatedAt: utils.ToNullTime(time.Time{}),
				UpdatedAt: utils.ToNullTime(time.Time{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.toDTO(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.toDTO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromDTOToAPI(t *testing.T) {
	type args struct {
		userDTO dto.UsersDTO
	}
	tests := []struct {
		name string
		args args
		want []*api.User
	}{
		{
			name: "valid user",
			args: args{
				userDTO: dto.UsersDTO{
					{
						ID:        utils.ToNullString("123"),
						FirstName: utils.ToNullString("John"),
						LastName:  utils.ToNullString("Doe"),
						Nickname:  utils.ToNullString("johndoe"),
						Password:  utils.ToNullString("password123"),
						Email:     utils.ToNullString("john.doe@example.com"),
						Country:   utils.ToNullString("USA"),
					},
				},
			},
			want: []*api.User{
				{
					Id:        "123",
					FirstName: "John",
					LastName:  "Doe",
					Nickname:  "johndoe",
					Email:     "john.doe@example.com",
					Country:   "USA",
				},
			},
		},
		{
			name: "empty user",
			args: args{
				userDTO: dto.UsersDTO{
					{
						ID:        utils.ToNullString(""),
						FirstName: utils.ToNullString(""),
						LastName:  utils.ToNullString(""),
						Nickname:  utils.ToNullString(""),
						Password:  utils.ToNullString(""),
						Email:     utils.ToNullString(""),
						Country:   utils.ToNullString(""),
					},
				},
			},
			want: []*api.User{
				{
					Id:        "",
					FirstName: "",
					LastName:  "",
					Nickname:  "",
					Email:     "",
					Country:   "",
				},
			},
		},
		{
			name: "multiple users",
			args: args{
				userDTO: dto.UsersDTO{
					{
						ID:        utils.ToNullString("123"),
						FirstName: utils.ToNullString("John"),
						LastName:  utils.ToNullString("Doe"),
						Nickname:  utils.ToNullString("johndoe"),
						Password:  utils.ToNullString("password123"),
						Email:     utils.ToNullString("john.doe@example.com"),
						Country:   utils.ToNullString("USA"),
					},
					{
						ID:        utils.ToNullString("456"),
						FirstName: utils.ToNullString("Jane"),
						LastName:  utils.ToNullString("Doe"),
						Nickname:  utils.ToNullString("janedoe"),
						Password:  utils.ToNullString("password456"),
						Email:     utils.ToNullString("jane.doe@example.com"),
						Country:   utils.ToNullString("USA"),
					},
				},
			},
			want: []*api.User{
				{
					Id:        "123",
					FirstName: "John",
					LastName:  "Doe",
					Nickname:  "johndoe",
					Email:     "john.doe@example.com",
					Country:   "USA",
				},
				{
					Id:        "456",
					FirstName: "Jane",
					LastName:  "Doe",
					Nickname:  "janedoe",
					Email:     "jane.doe@example.com",
					Country:   "USA",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromDTOToAPI(tt.args.userDTO); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromDTOToAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateUserChangeNotification(t *testing.T) {
	timeNow := time.Now()
	type args struct {
		changeType string
		userID     string
		eventTime  time.Time
	}
	tests := []struct {
		name string
		args args
		want UserChange
	}{
		{
			name: "valid user change",
			args: args{
				changeType: "create",
				userID:     "123",
				eventTime:  timeNow,
			},
			want: UserChange{
				ChangeType: "create",
				EventTime:  timeNow.Format(time.RFC3339),
				UserID:     "123",
			},
		},
		{
			name: "empty user change",
			args: args{
				changeType: "",
				userID:     "",
				eventTime:  time.Time{},
			},
			want: UserChange{
				ChangeType: "",
				EventTime:  time.Time{}.Format(time.RFC3339),
				UserID:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateUserChangeNotification(tt.args.changeType, tt.args.userID, tt.args.eventTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUserChangeNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}
