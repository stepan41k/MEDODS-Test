package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type User interface {
	CreateUser(ctx context.Context, guid []byte) (guidR []byte, err error)
}

type UserService struct {
	log *slog.Logger
	user User
}

func New(log *slog.Logger, user User) *UserService {
	return &UserService{
		log: log,
		user: user,
	} 
}


func (us *UserService) CreateUser(ctx context.Context) ([]byte, error) {
	const op = "service.user.Create"

	log := us.log.With(
		slog.String("op", op),
	)

	log.Info("start creating")

	guid, err := uuid.NewRandom()
	if err != nil {
		log.Error("failed to generate guid")

		return nil,fmt.Errorf("%s: %w", op, err)
	}

	log.Info("GUID generated")

	guidR, err := us.user.CreateUser(ctx, []byte(guid.String()))
	if err != nil {
		log.Error("failed to create new user")

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user created")

	return guidR, nil
}