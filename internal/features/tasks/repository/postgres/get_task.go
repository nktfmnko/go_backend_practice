package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"practice/internal/core/domain"
	core_error "practice/internal/core/errors"
	core_postgres_pool "practice/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(ctx context.Context, id int) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	WHERE id = $1;
	`

	var taskModel TaskModel
	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id=%d :, %w", id, core_error.ErrNotFound)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	task := taskDomainFromModel(taskModel)

	return task, nil
}
