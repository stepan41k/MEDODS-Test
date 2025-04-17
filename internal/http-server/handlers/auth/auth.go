package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	resp "github.com/stepan41k/MEDODS-Test/internal/lib/api/response"
	"github.com/stepan41k/MEDODS-Test/internal/lib/cookie"
	"github.com/stepan41k/MEDODS-Test/internal/lib/sl"
	"github.com/stepan41k/MEDODS-Test/internal/service/auth"
)

type Auth interface {
	Create(ctx context.Context, guid []byte, ip string) (token string, err error)
	Refresh(ctx context.Context, ip string, accessCookie string) (token string, err error)
}

type AuthHandler struct {
	log *slog.Logger
	auth Auth
}

func New(log *slog.Logger, auth Auth) *AuthHandler {
	return &AuthHandler{
		log: log,
		auth: auth,
	}
}


func (ah *AuthHandler) Create(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.Create"

		log := ah.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		guid := chi.URLParam(r, "guid")
		if guid == "" {
			log.Warn("guid is empty")

			render.Status(r, http.StatusConflict)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusConflict,
				Error: "invalid request",
			})

			return
		}

		ip := r.RemoteAddr
		accessCookie, err := cookie.TakeCookie(w, r)
		if accessCookie != "" || err == nil {
			log.Error("tokens already created")

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusBadRequest,
				Error: "tokens already created",
			})

			return 
		}

		token, err := ah.auth.Create(ctx, []byte(guid), ip)
		if err != nil {
			if errors.Is(err, auth.ErrUserNotFound) {
				log.Error("user not found:", sl.Err(auth.ErrUserNotFound))

				render.Status(r, http.StatusBadRequest)

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusBadRequest,
					Error: "user not found",
				})

				return
			}

			log.Error("failed to create tokens", sl.Err(err))

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusBadRequest,
				Error: "failed to create tokens",
			})

			return
		}

		cookie.SetCookie(w, token)

		log.Info("tokens created")

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data: "tokens created",
		})
	}
}


func (ah *AuthHandler) Refresh(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.Refresh"

		log := ah.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("start refreshing tokens")

		ip := r.RemoteAddr

		accessCookie, err := cookie.TakeCookie(w, r)
		if err != nil {
			if errors.Is(err, cookie.ErrCookieNotSet) {
				log.Error("token not set", sl.Err(cookie.ErrCookieNotSet))

				render.Status(r, http.StatusBadRequest)

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusUnauthorized,
					Error: "first you need to generate a token",
				})

				return
			}

			log.Error("failed to take cookie", sl.Err(err))

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusBadRequest,
				Error: "failed to take cookie",
			})

			return 
		}

		log.Info("took cookie")

		token, err := ah.auth.Refresh(ctx, ip, accessCookie)
		if err != nil {
			log.Error("failed to refresh tokens", sl.Err(err))

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusBadRequest,
				Error: "failed to refresh tokens",
			})

			return
		}

		cookie.SetCookie(w, token)

		log.Info("tokens refreshed")

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data: "tokens refreshed",
		})
	}
}