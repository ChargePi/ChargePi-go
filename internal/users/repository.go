package users

import (
	"github.com/ChargePi/ChargePi-go/internal/users/models"
)

type UserRepository interface {
	GetUser(username string) (*models.User, error)
	GetUsers() ([]models.User, error)
	AddUser(user models.User) error
	UpdateUser(models.User) (*models.User, error)
	DeleteUser(username string) error
}
