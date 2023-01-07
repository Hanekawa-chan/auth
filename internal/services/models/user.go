package models

type CreateUserRequest struct {
	Email   string
	Country string
}

type User struct {
	Id      string
	Name    string
	Email   string
	Country string
}
