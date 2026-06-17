package handler

import (
	"encoding/json"
	"myprojects/internal/models"
	"myprojects/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "id")

	idInt64, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByID(idInt64)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError) // status code 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.CreateUser(user)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "user created",
	})
}
