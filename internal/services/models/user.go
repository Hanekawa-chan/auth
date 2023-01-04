package models

type CreateUserRequest struct {
	Username string
	Country  string
}

type User struct {
	Id       string
	Name     string
	Username string
	Email    string
	Country  string
}
