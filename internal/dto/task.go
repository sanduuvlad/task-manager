package dto

type UpdateTaskDTO struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
