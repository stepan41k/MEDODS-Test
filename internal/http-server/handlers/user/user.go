package user

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	resp "github.com/stepan41k/MEDODS-Test/internal/lib/api/response"
	"github.com/stepan41k/MEDODS-Test/internal/lib/sl"
)

type User interface {
	CreateUser(ctx context.Context) (guid []byte, err error)
}

type UserHandler struct {
	log *slog.Logger
	user User
}

func New(log *slog.Logger, user User) *UserHandler {
	return &UserHandler{
		log: log,
		user: user,
	}
}


func (uh *UserHandler) CreateUser(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Create"

		log := uh.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		guid, err := uh.user.CreateUser(ctx)
		if err != nil {
			log.Error("failed to create user", sl.Err(err))

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusBadRequest,
				Error: "failed to create user",
			})

			return
		}

		log.Info("user created")

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data: string(guid),
		})
	}
}