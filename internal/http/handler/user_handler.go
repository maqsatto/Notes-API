package handler

import (
	"net/http"
	"strconv"

	"github.com/maqsatto/Notes-API/internal/domain"
	"github.com/maqsatto/Notes-API/internal/http/dto/request"
	"github.com/maqsatto/Notes-API/internal/http/dto/response"
	"github.com/maqsatto/Notes-API/internal/service"
	"github.com/maqsatto/Notes-API/internal/utils"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(UserService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: UserService,
	}
}

func toUserResponse(user *domain.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.RegisterRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Message: "Failed to parse request body",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
	}

	user, token, err := u.UserService.Register(r.Context(), req.Username, req.Email, req.Password);
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Registration failed",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusInternalServerError, errorRes)
		return
	}

	authRes := response.AuthResponse{
		User: toUserResponse(user),
		Token: token,
	}

	utils.WriteJSON(w, http.StatusCreated, authRes)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Message: "Failed to parse request body",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
	}

	user, token, err := u.UserService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Login failed",
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusUnauthorized, errorRes)
		return
	}
	authRes := response.AuthResponse{
		User: toUserResponse(user),
		Token: token,
	}
	utils.WriteJSON(w, http.StatusOK, authRes)
}

func (u *UserHandler) CheckEmail(w http.ResponseWriter, r *http.Request) {
	var req request.CheckEmailRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Message: "Failed to parse request body",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
		return
	}
	isTaken, err := u.UserService.IsEmailTaken(r.Context(), req.Email)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Check failed",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusInternalServerError, errorRes)
		return
	}
	res := response.AvailabilityResponse{
		Available: !isTaken,
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

func (u *UserHandler) CheckUsername(w http.ResponseWriter, r *http.Request) {
	var req request.CheckUsernameRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Message: "Failed to parse request body",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
		return
	}
	isTaken, err := u.UserService.IsUsernameTaken(r.Context(), req.Username)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Check failed",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusInternalServerError, errorRes)
		return
	}
	res := response.AvailabilityResponse{
		Available: !isTaken,
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

func (u *UserHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	if email == "" {
		errorRes := response.ErrorResponse{
			Error:   "Invalid request",
			Code:    http.StatusBadRequest,
			Message: "Email parameter is required",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
		return
	}

	user, err := u.UserService.GetByEmail(r.Context(), email)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "User not found",
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusNotFound, errorRes)
		return
	}
	utils.WriteJSON(w, http.StatusOK, toUserResponse(user))
}

func (u *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		errorRes := response.ErrorResponse{
			Error:   "Invalid request",
			Code:    http.StatusBadRequest,
			Message: "ID parameter is required",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
		return
	}
		// Convert string to uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Invalid ID",
			Code:    http.StatusBadRequest,
			Message: "ID must be a valid number",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
		return
	}
	user, err := u.UserService.GetByID(r.Context(), id)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "User not found",
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusNotFound, errorRes)
		return
	}
	utils.WriteJSON(w, http.StatusOK, toUserResponse(user))
}
func (u *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		errorRes := response.ErrorResponse{
			Error:   "Invalid request",
			Code:    http.StatusBadRequest,
			Message: "Username parameter is required",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
		return
	}
	user, err := u.UserService.GetByUsername(r.Context(), username)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "User not found",
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusNotFound, errorRes)
		return
	}
	utils.WriteJSON(w, http.StatusOK, toUserResponse(user))
}
