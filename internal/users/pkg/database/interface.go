package database

import "github.com/ChargePi/ChargePi-go/internal/users/pkg/models"

type Database interface {
	GetUser(username string) (*models.User, error)
	GetUsers() []models.User
	AddUser(user models.User) error
	UpdateUser(models.User) (*models.User, error)
	DeleteUser(username string) error
}
