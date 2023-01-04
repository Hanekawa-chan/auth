package models

type CreateUserRequest struct {
	Username string
	Country  string
}

type User struct {
	Id string
}
