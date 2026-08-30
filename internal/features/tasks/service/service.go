package tasks_service

import (
	"context"
	"practice/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

func NewTasksService(tasksRepository TasksRepository) *TasksService {
	return &TasksService{tasksRepository: tasksRepository}
}

type TasksRepository interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
}
