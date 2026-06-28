package repository

import (
	"context"
	"errors"
	"myprojects/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("User not found")

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (repo *UserRepository) GetUserByID(id int64) (models.User, error) {
	row := repo.pool.QueryRow(
		context.Background(),
		"SELECT id, username, email FROM users WHERE id = $1",
		id,
	)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) CreateUser(user models.User) error {
	_, err := repo.pool.Exec(
		context.Background(),
		`INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2, $3)`,
		user.Username,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User

	rows, err := repo.pool.Query(
		context.Background(),
		`SELECT id, username, email, created_at
		FROM users`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	query := `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`

	err := repo.pool.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil
}
