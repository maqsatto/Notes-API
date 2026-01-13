// internal/http/router/router.go
package router

import (
	"net/http"

	"github.com/maqsatto/Notes-API/internal/config"
	"github.com/maqsatto/Notes-API/internal/http/middleware"
	"github.com/maqsatto/Notes-API/internal/logger"
	"github.com/maqsatto/Notes-API/internal/utils"
)

type Deps struct {
	Config *config.Config
	Logger *logger.Logger
	DB     any
}

func New(d Deps) http.Handler {
	mux := http.NewServeMux()


	//USERS
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	var h http.Handler = mux
	h = middleware.Recovery(d.Logger)(h)
	h = middleware.Logger(d.Logger)(h)
	h = middleware.CORS(h)

	return h
}
