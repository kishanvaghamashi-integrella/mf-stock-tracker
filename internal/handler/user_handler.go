package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/dto"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/service"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserService(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Create godoc
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserRequest true "Create user payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/ [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "UserHandler.Create", "method", r.Method, "path", r.URL.Path)

	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode request body", "handler", "UserHandler.Create", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(req); err != nil {
		slog.Warn("validation failed", "handler", "UserHandler.Create", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	if err := h.service.Create(r.Context(), &req); err != nil {
		util.HandleError(w, err, "UserHandler.Create")
		return
	}

	slog.Info("user created successfully", "handler", "UserHandler.Create")
	util.SendResponse(w, http.StatusOK, map[string]string{"message": "user created successfully"})
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequest true "Login payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "UserHandler.Login", "method", r.Method, "path", r.URL.Path)

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode request body", "handler", "UserHandler.Login", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(req); err != nil {
		slog.Warn("validation failed", "handler", "UserHandler.Login", "error", err)
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	user, err := h.service.Login(r.Context(), &req)
	if err != nil {
		util.HandleError(w, err, "UserHandler.Login")
		return
	}

	token, err := util.GenerateToken(user.ID, user.Email)
	if err != nil {
		slog.Error("failed to generate token", "handler", "UserHandler.Login", "userID", user.ID, "error", err)
		util.SendErrorResponse(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	slog.Info("user logged in successfully", "handler", "UserHandler.Login", "userID", user.ID)
	util.SendResponse(w, http.StatusOK, map[string]any{
		"message": "login successful",
		"user": dto.LoginResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		},
	})
}

// Verify godoc
// @Summary Verify token
// @Description Verify bearer token and return user info
// @Tags users
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/verify [get]
// @Security BearerAuth
func (h *UserHandler) Verify(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "UserHandler.Verify", "method", r.Method, "path", r.URL.Path)

	userID, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "UserHandler.Verify")
		util.SendErrorResponse(w, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		util.HandleError(w, err, "UserHandler.Verify")
		return
	}

	token := r.Header.Get("Authorization")

	token = strings.TrimPrefix(token, "Bearer ")

	slog.Info("token verified successfully", "handler", "UserHandler.Verify", "userID", userID)
	util.SendResponse(w, http.StatusOK, map[string]any{
		"message": "token is valid",
		"user": dto.LoginResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		},
	})
}

// Delete godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users [delete]
// @Security BearerAuth
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slog.Info("request started", "handler", "UserHandler.Delete", "method", r.Method, "path", r.URL.Path)

	userId, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "UserHandler.Delete")
		util.SendErrorResponse(w, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	if err := h.service.Delete(r.Context(), userId); err != nil {
		util.HandleError(w, err, "UserHandler.Delete")
		return
	}

	slog.Info("user deleted successfully", "handler", "UserHandler.Delete", "userID", userId)
	util.SendResponse(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}
