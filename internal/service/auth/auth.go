package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/stepan41k/MEDODS-Test/internal/lib/jwt"
	"github.com/stepan41k/MEDODS-Test/internal/storage"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Auth interface {
	Create(ctx context.Context, guid []byte, refreshToken []byte) error
	GetRefresh(ctx context.Context, guid []byte) (refreshToken string, err error)
}

type AuthService struct {
	log *slog.Logger
	auth Auth
}

func New(log *slog.Logger, auth Auth) *AuthService {
	return &AuthService{
		log: log,
		auth: auth,
	} 
}


func (as *AuthService) Create(ctx context.Context, guid []byte, ip string) (string, error) {
	const op = "service.auth.Create"

	log := as.log.With(
		slog.String("op", op),
		slog.String("guid", string(guid)),
	)

	newKey := uuid.NewString()

	accessToken, err := jwt.NewAccesssTokens(string(guid), newKey, ip)
	if err != nil {
		log.Error("failed to genrate new access token")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := jwt.NewRefreshToken(newKey, ip)
	if err != nil {
		log.Error("failed to genrate new refresh token")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := as.auth.Create(ctx, guid, []byte(refreshToken)); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("created")

	return accessToken, nil
}


func (as *AuthService) Refresh(ctx context.Context, ip string, accessCookie string) (token string, err error) {
	const op = "service.auth.Refresh"

	guid := jwt.GetGUID(accessCookie)

	log := as.log.With(
		slog.String("op", op),
		slog.String("GUID", string(guid)),
	)

	newKey := uuid.NewString()

	refreshToken, err := as.auth.GetRefresh(ctx, guid)
	if err != nil {
		log.Error("failed to get refresh token")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	flag, err := jwt.CheckRefresh(refreshToken) 
	if err != nil || !flag  {
		log.Error("failed to check refresh token")

		return "", fmt.Errorf("%s: %w", op, err)
	}
	
	newAccess, err := jwt.NewAccesssTokens(string(guid), newKey, ip)
	if err != nil {
		log.Error("failed to genrate new access token")

		return "", fmt.Errorf("%s: %w", op, err)
	}
	newRefresh, err := jwt.NewRefreshToken(newKey, ip)
	if err != nil {
		log.Error("failed to genrate new refresh token")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("create new tokens")

	if err = as.auth.Create(ctx, guid, []byte(newRefresh)); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("refreshed")

	return newAccess, nil
}