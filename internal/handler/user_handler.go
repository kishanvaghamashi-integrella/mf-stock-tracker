package handler

import (
	"encoding/json"
	"net/http"

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
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(req); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	if err := h.service.Create(r.Context(), &req); err != nil {
		handleError(w, err)
		return
	}

	util.SendResponse(w, http.StatusOK, map[string]string{"message": "user created successfully"})
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequest true "Login payload"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.Validate.Struct(req); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	user, err := h.service.Login(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	token, err := util.GenerateToken(user.ID, user.Email)
	if err != nil {
		util.SendErrorResponse(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	util.SendResponse(w, http.StatusOK, map[string]any{
		"message": "user created successfully",
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
// @Param userId path int64 true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/{userId} [delete]
// @Security BearerAuth
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserIDFromContext(r.Context())
	if ok == false {
		util.SendErrorResponse(w, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	if err := h.service.Delete(r.Context(), userId); err != nil {
		handleError(w, err)
		return
	}

	util.SendResponse(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}
