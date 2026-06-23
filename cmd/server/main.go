package main

import (
	"log"
	"myprojects/internal/config"
	"myprojects/internal/database"
	"myprojects/internal/handler"
	"myprojects/internal/middleware"
	"myprojects/internal/repository"
	"myprojects/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	pool, err := database.Connect(cfg.DATABASEUrl)
	if err != nil {
		panic(err)
	}

	taskRepo := repository.NewTaskRepository(pool)

	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo, cfg)

	taskHandler := handler.TaskHandler{
		Repo: taskRepo,
	}

	userHandler := handler.UserHandler{
		Service: userService,
	}

	router := chi.NewRouter()

	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)

	router.Get("/tasks", taskHandler.GetTasks)
	router.Post("/tasks", taskHandler.CreateTask)
	router.Get("/tasks/{id}", taskHandler.GetTaskByID)
	router.Put("/tasks/{id}", taskHandler.UpdateTaskByID)
	router.Delete("/tasks/{id}", taskHandler.DeleteTaskByID)
	router.Patch("/tasks/{id}", taskHandler.PatchTaskByID)

	router.Get("/users", userHandler.GetAllUsers)
	router.Post("/users", userHandler.CreateUser)
	router.Get("/users/{id}", userHandler.GetUserByID)

	log.Println("Server started on :8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
