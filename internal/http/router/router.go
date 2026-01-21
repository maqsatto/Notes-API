package router

import (
	"database/sql"
	"net/http"

	"github.com/maqsatto/Notes-API/internal/auth"
	"github.com/maqsatto/Notes-API/internal/config"
	"github.com/maqsatto/Notes-API/internal/http/handler"
	"github.com/maqsatto/Notes-API/internal/http/middleware"
	"github.com/maqsatto/Notes-API/internal/logger"
	"github.com/maqsatto/Notes-API/internal/repository"
	"github.com/maqsatto/Notes-API/internal/service"
	"github.com/maqsatto/Notes-API/internal/utils"
)

type Deps struct {
	Config *config.Config
	Logger *logger.Logger
	DB     any
	JWT    *auth.JWTManager
}

func New(d Deps) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Initialize repository, service, and handler
	db := d.DB.(*sql.DB)
	userRepo := repository.NewUserRepo(db)
	noteRepo := repository.NewNoteRepo(db)
	userSvc := service.NewUserService(*userRepo, *noteRepo, d.JWT)
	userHandler := handler.NewUserHandler(*userSvc)

	// Public Routes (No Auth)
	mux.HandleFunc("POST /api/users/register", userHandler.Register)
	mux.HandleFunc("POST /api/users/login", userHandler.Login)

	// Email and username availability checks
	mux.HandleFunc("POST /api/users/check-email", userHandler.CheckEmail)
	mux.HandleFunc("POST /api/users/check-username", userHandler.CheckUsername)

	// Protected Routes (Auth Required)
	authMW := middleware.AuthMiddleware(d.JWT)

	// // // User profile management
	// mux.Handle("GET /api/users/me", authMW(http.HandlerFunc(userHandler.GetProfile)))
	// mux.Handle("PUT /api/users/me", authMW(http.HandlerFunc(userHandler.UpdateProfile)))
	// mux.Handle("PUT /api/users/me/password", authMW(http.HandlerFunc(userHandler.ChangePassword)))
	// mux.Handle("DELETE /api/users/me", authMW(http.HandlerFunc(userHandler.DeleteAccount)))
	// mux.Handle("DELETE /api/users/me/permanent", authMW(http.HandlerFunc(userHandler.PermanentDeleteAccount)))

	// User lookup by ID, email, username
	mux.Handle("GET /api/users/{id}", authMW(http.HandlerFunc(userHandler.GetByID)))
	mux.Handle("GET /api/users/email/{email}", authMW(http.HandlerFunc(userHandler.GetByEmail)))
	mux.Handle("GET /api/users/username/{username}", authMW(http.HandlerFunc(userHandler.GetByUsername)))

	// // List users with pagination
	// mux.Handle("GET /api/users", authMW(http.HandlerFunc(userHandler.ListUsers)))

	// // User statistics
	// mux.Handle("GET /api/users/stats/total", authMW(http.HandlerFunc(userHandler.GetTotalUserCount)))

	// Global Middleware Chain
	var h http.Handler = mux
	h = middleware.Recovery(d.Logger)(h)
	h = middleware.Logger(d.Logger)(h)
	h = middleware.CORS(h)

	return h
}
