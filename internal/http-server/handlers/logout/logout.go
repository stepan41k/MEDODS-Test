package logout

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	resp "github.com/stepan41k/MEDODS-Test/internal/lib/api/response"
	"github.com/stepan41k/MEDODS-Test/internal/lib/cookie"
)


func Delete(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "session.FinishSession"

		log = log.With(
			"op", op,
		)

		cookie.DeleteCookie(w)

		log.Info("cookie deleted")

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data: "successful logout",
		})
	}
}  
