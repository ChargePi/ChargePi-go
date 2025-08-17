package users

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/casbin/casbin/v2"

	"github.com/ChargePi/ChargePi-go/internal/users/models"
	"github.com/ChargePi/ChargePi-go/pkg/encryption"
)

var (
	ErrNoPermissions = errors.New("the user does not have sufficient permissions")
)

type (
	Service interface {
		GetUser(ctx context.Context, username string) (*models.User, error)
		GetUsers(ctx context.Context) ([]models.User, error)
		AddUser(ctx context.Context, username, password, role string) error
		UpdateUser(ctx context.Context, username string, password, role *string) (*models.User, error)
		DeleteUser(ctx context.Context, username string) error
		CheckPassword(ctx context.Context, username, password string) bool
	}

	UserService struct {
		database  UserRepository
		enforcer  *casbin.Enforcer
		encryptor encryption.Encryptor
		logger    *zap.Logger
	}
)

func NewUserService(logger *zap.Logger, db UserRepository) *UserService {
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
		// enforcer: e,
		logger: logger.Named("user_service"),
	}
}

func (u *UserService) GetUser(ctx context.Context, username string) (*models.User, error) {
	u.logger.With(zap.String("user", username)).Info("Getting user")
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	user, err := u.database.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	u.logger.Info("Getting users")
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	return u.database.GetUsers(ctx)
}

func (u *UserService) AddUser(ctx context.Context, username, password, role string) error {
	u.logger.With(zap.String("user", username)).Info("Adding a user")

	user := models.User{
		Username: username,
		Password: password,
		Role:     models.Role(role),
	}

	err := models.ValidateRole(role)
	if err != nil {
		return err
	}
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	// Encrypt the password before storing the user
	encrypt, err := u.encryptor.Encrypt(user.Password)
	if err != nil {
		return err
	}

	user.Password = *encrypt

	return u.database.AddUser(ctx, user)
}

func (u *UserService) UpdateUser(ctx context.Context, username string, password, role *string) (*models.User, error) {
	u.logger.With(zap.String("user", username)).Info("Updating a user")
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	if role != nil {
		err := models.ValidateRole(*role)
		if err != nil {
			return nil, err
		}
	}

	iUser := models.User{Username: username}

	// Encrypt the password before storing the user
	if password != nil {
		encryptedPass, err := u.encryptor.Encrypt(*password)
		if err != nil {
			return nil, err
		}

		iUser.Password = *encryptedPass
	}

	user, err := u.database.UpdateUser(ctx, iUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) DeleteUser(ctx context.Context, username string) error {
	u.logger.With(zap.String("user", username)).Info("Deleting a user")
	// todo check for access

	/*enforce, err := u.enforcer.Enforce(username)
	if err != nil {
		return nil, err
	}*/

	return u.database.DeleteUser(ctx, username)
}

func (u *UserService) CheckPassword(ctx context.Context, username, password string) bool {
	u.logger.With(zap.String("user", username)).Info("Checking user password")
	user, err := u.database.GetUser(ctx, username)
	if err != nil {
		return false
	}

	// Decrypt the password and compare
	decrypt, err := u.encryptor.Decrypt(user.Password)
	if err != nil {
		return false
	}

	return *decrypt == password
}

func (u *UserService) Pass() bool {
	return true
}

func (u *UserService) Name() string {
	return "user-service"
}
