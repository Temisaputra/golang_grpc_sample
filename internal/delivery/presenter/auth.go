package presenter

import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Role     string `json:"role"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

var validate = validator.New()

// validation login
func ValidateLogin(req LoginRequest) error {
	return validate.Struct(req)
}

// validation register
func ValidateRegister(req RegisterRequest) error {
	return validate.Struct(req)
}
