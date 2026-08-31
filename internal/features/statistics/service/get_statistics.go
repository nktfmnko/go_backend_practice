package statistics_service

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
	core_error "practice/internal/core/errors"
	"time"
)

func (s *StatisticsService) GetStatistics(ctx context.Context, userID *int, from, to *time.Time) (domain.Statistics, error) {
	if from != nil && to != nil {
		if from.Before(*to) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("to must be after from: %w", core_error.ErrInvalidArgument)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repo: %w", err)
	}

	stats := calcStatistics(tasks)

	return stats, err
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	tasksCreated := len(tasks)
	tasksCompleted := 0
	if tasksCreated == 0 {
		return domain.Statistics{
			TasksCreated:               tasksCreated,
			TasksCompleted:             tasksCompleted,
			TasksCompletedRate:         nil,
			TasksAverageCompletionTime: nil,
		}
	}

	var totalCompletedDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}
		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletedDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletedDuration != 0 {
		avg := totalCompletedDuration / time.Duration(tasksCompleted)
		tasksAverageCompletionTime = &avg
	}

	return domain.Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         &tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
