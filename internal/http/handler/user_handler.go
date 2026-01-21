package handler

import (
	"net/http"
	"strconv"

	"github.com/maqsatto/Notes-API/internal/domain"
	"github.com/maqsatto/Notes-API/internal/http/dto/request"
	"github.com/maqsatto/Notes-API/internal/http/dto/response"
	"github.com/maqsatto/Notes-API/internal/http/middleware"
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

func (u *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		errorRes := response.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in context",
		}
		utils.WriteJSON(w, http.StatusUnauthorized, errorRes)
		return
	}

	user, err := u.UserService.GetByID(r.Context(), userID)
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

func (u *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req request.UpdateProfileRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Message: "Failed to parse request body",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
		return
	}
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		errorRes := response.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in context",
		}
		utils.WriteJSON(w, http.StatusUnauthorized, errorRes)
		return
	}
	user, err := u.UserService.UpdateProfile(r.Context(), userID, req.Username, req.Email)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Update failed",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusInternalServerError, errorRes)
		return
	}

	utils.WriteJSON(w, http.StatusOK, toUserResponse(user))
}

func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req request.ChangePasswordRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Message: "Failed to parse request body",
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
		return
	}
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		errorRes := response.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in context",
		}
		utils.WriteJSON(w, http.StatusUnauthorized, errorRes)
		return
	}
	if err := u.UserService.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Password change failed",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusBadRequest, errorRes)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response.MessageResponse{
		Message: "Password changed successfully",
	})
}

func (u *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		errorRes := response.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in context",
		}
		utils.WriteJSON(w, http.StatusUnauthorized, errorRes)
		return
	}
	if err := u.UserService.DeleteAccount(r.Context(), userID); err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Delete failed",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusInternalServerError, errorRes)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response.MessageResponse{
		Message: "Account deleted successfully",
	})
}

func (u *UserHandler) PermanentDeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		errorRes := response.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in context",
		}
		utils.WriteJSON(w, http.StatusUnauthorized, errorRes)
		return
	}
	if err := u.UserService.PermanentDeleteAccount(r.Context(), userID); err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Permanent delete failed",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusInternalServerError, errorRes)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response.MessageResponse{
		Message: "Account permanently deleted",
	})
}

func (u *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	var req request.ListUsersRequest
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			errorRes := response.ErrorResponse{
				Error:   "Invalid parameter",
				Code:    http.StatusBadRequest,
				Message: "Limit must be a valid number",
			}
			utils.WriteJSON(w, http.StatusBadRequest, errorRes)
			return
		}
		req.Limit = limit
	}
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			errorRes := response.ErrorResponse{
				Error:   "Invalid parameter",
				Code:    http.StatusBadRequest,
				Message: "Offset must be a valid number",
			}
			utils.WriteJSON(w, http.StatusBadRequest, errorRes)
			return
		}
		req.Offset = offset
	}

	req.SetDefaults()

	users, total, err := u.UserService.ListUsers(r.Context(), req.Limit, req.Offset)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Failed to retrieve users",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusInternalServerError, errorRes)
		return
	}
	userResponses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = toUserResponse(user)
	}
	listResponse := response.UserListResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	listResponse.CalculateTotalPages()

	utils.WriteJSON(w, http.StatusOK, listResponse)
}

func (u *UserHandler) GetTotalUserCount(w http.ResponseWriter, r *http.Request) {
	total, err := u.UserService.GetTotalUserCount(r.Context())
	if err != nil {
		errorRes := response.ErrorResponse{
			Error:   "Failed to get user count",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		utils.WriteJSON(w, http.StatusInternalServerError, errorRes)
		return
	}
	statsResponse := response.StatsResponse{
		TotalUsers: total,
	}
	utils.WriteJSON(w, http.StatusOK, statsResponse)
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
