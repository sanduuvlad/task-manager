package handler

import (
	"encoding/json"
	"errors"
	"myprojects/internal/dto"
	"myprojects/internal/models"
	"myprojects/internal/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	Repo *repository.TaskRepository
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Repo.GetTasks()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Repo.CreateTask(task)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "id")

	idInt64, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.Repo.GetTaskByID(idInt64)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "id")

	idInt64, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var task models.Task

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Repo.UpdateTaskByID(idInt64, task)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "id")

	idInt64, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteTaskByID(idInt64)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) PatchTaskByID(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "id")

	idInt64, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var dto dto.UpdateTaskDTO

	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.Repo.PatchTaskByID(idInt64, dto)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
