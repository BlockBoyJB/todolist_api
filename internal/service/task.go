package service

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
	"todolist_api/internal/model/dbmodel"
	"todolist_api/internal/repo"
	"todolist_api/internal/repo/pgerrs"
)

const (
	taskServicePrefixLog = "/service/task"
)

type taskService struct {
	task repo.Task
}

func newTaskService(task repo.Task) *taskService {
	return &taskService{task: task}
}

func (s *taskService) Create(ctx context.Context, input TaskCreateInput) (TaskOutput, error) {
	task := &dbmodel.Task{
		Username:    input.Username,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
	}
	if err := s.task.Create(ctx, task); err != nil {
		if errors.Is(err, pgerrs.ErrForeignKey) {
			return TaskOutput{}, ErrUserNotFound
		}
		log.Errorf("%s/Create error create task: %s", taskServicePrefixLog, err)
		return TaskOutput{}, err
	}
	return TaskOutput{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate.Format(time.RFC3339),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *taskService) Find(ctx context.Context, username string) ([]TaskOutput, error) {
	tasks, err := s.task.Find(ctx, username)
	if err != nil {
		log.Errorf("%s/Find error find user tasks: %s", taskServicePrefixLog, err)
		return nil, err
	}

	result := make([]TaskOutput, 0)
	for _, t := range tasks {
		result = append(result, TaskOutput{
			Id:          t.Id,
			Title:       t.Title,
			Description: t.Description,
			DueDate:     t.DueDate.Format(time.RFC3339),
			CreatedAt:   t.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

func (s *taskService) FindById(ctx context.Context, id int, username string) (TaskOutput, error) {
	task, err := s.task.FindById(ctx, id, username)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return TaskOutput{}, ErrTaskNotFound
		}
		log.Errorf("%s/FindById error find task by id: %s", taskServicePrefixLog, err)
		return TaskOutput{}, err
	}
	return TaskOutput{
		Id:          id,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate.Format(time.RFC3339),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *taskService) Update(ctx context.Context, input TaskUpdateInput) (TaskOutput, error) {
	task := &dbmodel.Task{
		Id:          input.Id,
		Username:    input.Username,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
	}

	if err := s.task.Update(ctx, task); err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return TaskOutput{}, ErrTaskNotFound
		}
		log.Errorf("%s/Update error update user task: %s", taskServicePrefixLog, err)
		return TaskOutput{}, err
	}

	return TaskOutput{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate.Format(time.RFC3339),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *taskService) Delete(ctx context.Context, id int, username string) error {
	if err := s.task.Delete(ctx, id, username); err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Errorf("%s/Delete error delete user task: %s", userServicePrefixLog, err)
		return err
	}
	return nil
}
