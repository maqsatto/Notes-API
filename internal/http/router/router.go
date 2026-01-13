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
	DB     any // replace with your real db type (*sql.DB or *sqlx.DB)
}

func New(d Deps) http.Handler {
	mux := http.NewServeMux()

	// register routes here (health, notes, users...)
	// mux.HandleFunc("/health", ...)

	// wrap with middleware here (recovery, logger, cors...)
	// handler := middleware.Recovery(d.Logger)(mux)
	// handler = middleware.Logger(d.Logger)(handler)

	return mux
}
