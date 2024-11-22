package server

import (
	"fmt"

	"github.com/EFG/api"
)

type requiredFields map[string]string

func validateCreateUserRequest(req *api.CreateUserRequest) error {
	fields := requiredFields{
		"FirstName": req.FirstName,
		"LastName":  req.LastName,
		"Email":     req.Email,
		"Password":  req.Password,
		"Country":   req.Country,
		"NickName":  req.Nickname,
	}

	for field, value := range fields {
		if value == "" {
			return fmt.Errorf("%s cannot be empty", field)
		}
	}

	return nil
}

func validateExistingUserRequest(req *api.ModifyUserRequest) error {
	fields := requiredFields{
		"Id": req.Id,
	}

	for field, value := range fields {
		if value == "" {
			return fmt.Errorf("%s cannot be empty", field)
		}
	}

	return nil
}
