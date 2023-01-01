package users

import (
	"github.com/xBlaz3kx/ChargePi-go/internal/users/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/users/pkg/models"
)

type (
	Service interface {
		GetUser(username string) (*models.User, error)
		GetUsers() ([]models.User, error)
		AddUser(username, password, role string) error
		UpdateUser(username string, password, role *string) (*models.User, error)
		DeleteUser(username string) error
	}
)

type UserService struct {
	database database.Database
}

func NewUserService(db database.Database) *UserService {
	return &UserService{
		database: db,
	}
}

func (u *UserService) GetUser(username string) (*models.User, error) {
	_, err := u.database.GetUser(username)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UserService) GetUsers() ([]models.User, error) {
	return u.database.GetUsers(), nil
}

func (u *UserService) AddUser(username, password, role string) error {
	/*err := u.database.AddUser()
	if err != nil {
		return err
	}*/

	return nil
}

func (u *UserService) UpdateUser(username string, password, role *string) (*models.User, error) {
	/*_, err := u.database.UpdateUser()
	if err != nil {
		return nil, err
	}*/

	return nil, nil
}

func (u *UserService) DeleteUser(username string) error {
	return nil
}
