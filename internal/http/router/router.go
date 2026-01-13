// internal/http/router/router.go
package router

import (
	"net/http"

	"github.com/maqsatto/Notes-API/internal/config"
	"github.com/maqsatto/Notes-API/internal/logger"
)

type Deps struct {
	Config *config.Config
	Logger *logger.Logger
	DB     any
}

func New(d Deps) http.Handler {
	mux := http.NewServeMux()


	//USERS
	mux.Handle("GET /api/users", handler)
}
