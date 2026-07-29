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

	router.Get("/users/{id}", userHandler.GetUserByID)
	router.Post("/users", userHandler.CreateUser)
	router.Get("/users", userHandler.GetAllUsers)
	router.Post("/login", userHandler.Login)

	router.Group(func(r chi.Router) {
		r.Use(middleware.JWT(cfg.JWTSecret))

		r.Get("/tasks", taskHandler.GetTasks)
		r.Post("/tasks", taskHandler.CreateTask)
		r.Get("/tasks/{id}", taskHandler.GetTaskByID)
		r.Put("/tasks/{id}", taskHandler.UpdateTaskByID)
		r.Delete("/tasks/{id}", taskHandler.DeleteTaskByID)
		r.Patch("/tasks/{id}", taskHandler.PatchTaskByID)
	})
	log.Println("Server started on :8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
