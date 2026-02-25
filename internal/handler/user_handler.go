package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/model"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/service"
	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserService(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := util.Validate.Struct(user); err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	if err := h.service.Create(r.Context(), &user); err != nil {
		util.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.SendResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("User created with id %d", user.ID)})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "Error while parsing the userId")
		return
	}

	if err := h.service.Delete(r.Context(), userId); err != nil {
		if appErr, ok := err.(*util.AppError); ok {
			util.SendErrorResponse(w, appErr.Code, appErr.Message)
		} else {
			util.SendErrorResponse(w, http.StatusInternalServerError, "unexpected error")
		}
		return
	}

	util.SendResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully."})
}
