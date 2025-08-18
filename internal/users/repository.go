package users

import (
	"context"

	"github.com/ChargePi/ChargePi-go/internal/users/models"
)

type UserRepository interface {
	GetUser(ctx context.Context, username string) (*models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	AddUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, user models.User) (*models.User, error)
	DeleteUser(ctx context.Context, username string) error
}
