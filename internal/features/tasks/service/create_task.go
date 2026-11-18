package tasks_service

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
	core_auth "practice/internal/core/transport/auth"
)

func (s *TasksService) CreateTask(ctx context.Context, title string, description *string) (domain.Task, error) {
	claims, err := core_auth.MustUserFromContext(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("user unauthenticated: %w", err)
	}

	domainTask := domain.NewTaskUninitialized(title, description, claims.UserID)

	if err := domainTask.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	task, err := s.tasksRepository.CreateTask(ctx, domainTask)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
