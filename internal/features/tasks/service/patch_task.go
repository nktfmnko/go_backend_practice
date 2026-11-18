package tasks_service

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
	core_error "practice/internal/core/errors"
	core_auth "practice/internal/core/transport/auth"
)

func (s *TasksService) PatchTask(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error) {
	claims, err := core_auth.MustUserFromContext(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("user unauthenticated: %w", err)
	}

	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("fail to get task: %w", err)
	}

	if claims.UserID != task.AuthorUserID && !claims.IsAdmin {
		return domain.Task{}, fmt.Errorf("user %d is not author of task %d: %w", claims.UserID, id, core_error.ErrPermissionDenied)
	}

	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply patch task: %w", err)
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("fail to patch task: %w", err)
	}

	return patchedTask, nil
}
