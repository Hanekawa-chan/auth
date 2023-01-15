package models

type CreateUserRequest struct {
	Email   string
	Country string
}

type CreateUserResponse struct {
	UserId string
}
