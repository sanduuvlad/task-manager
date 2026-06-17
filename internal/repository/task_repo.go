package repository

import (
	"context"
	"errors"
	"fmt"
	"myprojects/internal/dto"
	"myprojects/internal/models"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		pool: pool,
	}
}

func (repo *TaskRepository) GetTasks() ([]models.Task, error) {
	var tasks []models.Task

	rows, err := repo.pool.Query(
		context.Background(),
		"SELECT id, title, description, created_at FROM tasks",
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var title, description string
		var createdAt time.Time

		err = rows.Scan(&id, &title, &description, &createdAt)
		if err != nil {
			return nil, err
		}

		task := models.Task{
			ID:          id,
			Title:       title,
			Description: description,
			CreatedAt:   createdAt,
		}

		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *TaskRepository) CreateTask(task models.Task) error {
	_, err := repo.pool.Exec(
		context.Background(),
		"INSERT INTO tasks (title, description) VALUES ($1, $2)",
		task.Title,
		task.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *TaskRepository) GetTaskByID(id int64) (models.Task, error) {
	row := repo.pool.QueryRow(
		context.Background(),
		"SELECT id, title, description, created_at FROM tasks WHERE id = $1",
		id,
	)

	var task models.Task

	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
	)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (repo *TaskRepository) UpdateTaskByID(id int64, task models.Task) error {
	_, err := repo.pool.Exec(
		context.Background(),
		"UPDATE tasks SET title = $1, description = $2 WHERE id = $3",
		task.Title,
		task.Description,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

var ErrTaskNotFound = errors.New("task not found")

func (repo *TaskRepository) DeleteTaskByID(id int64) error {
	tag, err := repo.pool.Exec(
		context.Background(),
		"DELETE FROM tasks WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (repo *TaskRepository) PatchTaskByID(id int64, dto dto.UpdateTaskDTO) error {
	var setClauses []string
	var args []any

	argPos := 1

	if dto.Title != nil {
		setClauses = append(
			setClauses,
			fmt.Sprintf("title = $%d", argPos),
		)

		args = append(args, *dto.Title)
		argPos++
	}

	if dto.Description != nil {
		setClauses = append(
			setClauses,
			fmt.Sprintf("description = $%d", argPos),
		)

		args = append(args, *dto.Description)
		argPos++
	}

	setQuery := strings.Join(setClauses, ", ")

	query := fmt.Sprintf(
		"UPDATE tasks SET %s WHERE id = $%d",
		setQuery,
		argPos,
	)

	args = append(args, id)

	fmt.Println(query)
	fmt.Println(args)

	_, err := repo.pool.Exec(
		context.Background(),
		query,
		args...,
	)
	if err != nil {
		return err
	}

	return nil
}
