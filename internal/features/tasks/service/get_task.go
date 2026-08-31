package tasks_service

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
)

func (s *TasksService) GetTask(ctx context.Context, id int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("fail get task from repository by id:%d: %w", id, err)
	}

	return task, nil
}
