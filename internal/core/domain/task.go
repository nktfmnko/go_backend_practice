package domain

import (
	"fmt"
	core_error "practice/internal/core/errors"
	"time"
)

type Task struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func NewTask(id, version int, title string, description *string, completed bool, createdAt time.Time, completedAt *time.Time, author int) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: author,
	}
}

func NewTaskUninitialized(title string, description *string, author int) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		author,
	)
}

func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf("invalid title len %w", core_error.ErrInvalidArgument)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("invalid description len %w", core_error.ErrInvalidArgument)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf("completedAt cant be nil if completed is true %w", core_error.ErrInvalidArgument)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf("completedAt cant be before createdAt %w", core_error.ErrInvalidArgument)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf("completedAt must be nil if completed is false %w", core_error.ErrInvalidArgument)
		}
	}

	return nil
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func NewTaskPatch(title Nullable[string], description Nullable[string], completed Nullable[bool]) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

func (t *TaskPatch) Validate() error {
	if t.Title.Set && t.Title.Value == nil {
		return fmt.Errorf("title cant be patched to null: %w", core_error.ErrInvalidArgument)
	}

	if t.Completed.Set && t.Completed.Value == nil {
		return fmt.Errorf("completed cant be patched to null: %w", core_error.ErrInvalidArgument)
	}
	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}
	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}
	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value
		if tmp.Completed {
			tmp.CompletedAt = new(time.Now())
		} else {
			tmp.CompletedAt = nil
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	*t = tmp
	return nil
}
