package users

import (
	"github.com/casbin/casbin/v2"
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

	UserService struct {
		database database.Database
		enforcer *casbin.Enforcer
	}
)

func NewUserService(db database.Database) *UserService {
	/*opts := badgerhold.DefaultOptions
	store, err := badgerhold.Open(opts)
	if err != nil {

	}

	a, err := badgeradapter.NewAdapter(store, "")
	if err != nil {

	}

	e, err := casbin.NewEnforcer("path/to/model.conf", a)
	if err != nil {

	}

	e.EnableEnforce(true)
	e.EnableLog(true)
	e.EnableAutoSave(true)*/

	return &UserService{
		database: db,
	}
}

func (u *UserService) GetUser(username string) (*models.User, error) {
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	user, err := u.database.GetUser(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) GetUsers() ([]models.User, error) {
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	return u.database.GetUsers(), nil
}

func (u *UserService) AddUser(username, password, role string) error {
	user := models.User{
		Username: username,
		Password: password,
		Role:     role,
	}

	err := validateRole(role)
	if err != nil {
		return err
	}
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	return u.database.AddUser(user)
}

func (u *UserService) UpdateUser(username string, password, role *string) (*models.User, error) {
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	if role != nil {
		err := validateRole(*role)
		if err != nil {
			return nil, err
		}
	}

	iUser := models.User{Username: username, Password: *password}
	user, err := u.database.UpdateUser(iUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) DeleteUser(username string) error {
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	return u.database.DeleteUser(username)
}

func validateRole(role string) error {
	switch models.Role(role) {
	case models.Manufacturer, models.Technician, models.Observer:
		return nil
	default:
		return nil
	}
}
