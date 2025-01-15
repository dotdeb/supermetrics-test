package data

type Database interface {
	GetUsers() []User
}
