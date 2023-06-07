package models

type User struct {
	Name    string
	Friends []User
}
