package tasks_postgres_repository

import (
	"context"
	"fmt"
	core_error "practice/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE
	FROM todoapp.tasks
	WHERE id = $1;
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("task with id=%d: %w", id, core_error.ErrInvalidArgument)
	}
	return nil
}
