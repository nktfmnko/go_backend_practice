package tasks_service

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
	core_error "practice/internal/core/errors"
)

func (s *TasksService) GetTasks(ctx context.Context, userID, limit, offset *int) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_error.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_error.ErrInvalidArgument)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository: %w", err)
	}

	return tasks, nil
}
