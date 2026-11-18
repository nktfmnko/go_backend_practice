package tasks_service

import (
	"context"
	"fmt"
	core_error "practice/internal/core/errors"
	core_auth "practice/internal/core/transport/auth"
)

func (s *TasksService) DeleteTask(ctx context.Context, id int) error {
	claims, err := core_auth.MustUserFromContext(ctx)
	if err != nil {
		return fmt.Errorf("user unauthenticated: %w", err)
	}

	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return fmt.Errorf("fail to get task: %w", err)
	}

	if claims.UserID != task.AuthorUserID && !claims.IsAdmin {
		return fmt.Errorf("user %d is not author of task %d: %w", claims.UserID, id, core_error.ErrPermissionDenied)
	}

	if err := s.tasksRepository.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("fail to delete task with id=%d : %w", id, err)
	}
	return nil
}
