package tasks_service

import (
	"context"
	"fmt"
)

func (s *TasksService) DeleteTask(ctx context.Context, id int) error {
	if err := s.tasksRepository.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("fail to delete task with id=%d : %w", id, err)
	}
	return nil
}
