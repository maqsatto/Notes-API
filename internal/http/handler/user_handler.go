package handler

import (
	"net/http"

	"github.com/maqsatto/Notes-API/internal/http/dto/request"
	"github.com/maqsatto/Notes-API/internal/service"
	"github.com/maqsatto/Notes-API/internal/utils"
)

type Handler struct {
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req request.RegisterRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	u, err := h.userService.Register(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, u)
}
