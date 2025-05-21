package usecase

import (
	"Test/internal/feature/task/entity"
	"Test/internal/feature/task/interface_adapters/dto"
	"Test/internal/feature/task/repository"
	"context"
	"time"
)

type updateTaskInteractor struct {
	taskUpdaterRepo repository.TaskUpdater
}

func NewUpdateTaskInteractor(taskUpdaterRepo repository.TaskUpdater) *updateTaskInteractor {
	return &updateTaskInteractor{taskUpdaterRepo: taskUpdaterRepo}
}

func (i *updateTaskInteractor) Update(ctx context.Context, req dto.UpdateTaskRequest) (dto.UpdateTaskResponse, error) {
	timeNow := time.Now()
	newTask := entity.Task{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		Deadline:    req.Deadline,
		AssigneeID:  req.AssigneeID,
		ClientID:    req.ClientID,
		UpdatedAt:   &timeNow,
	}

	updatedTask, err := i.taskUpdaterRepo.Update(ctx, newTask)
	if err != nil {
		return dto.UpdateTaskResponse{}, err
	}
	return dto.UpdateTaskResponse{
		ID:        updatedTask.ID,
		Status:    updatedTask.Status,
		UpdatedAt: updatedTask.UpdatedAt,
	}, nil
}

func (i *updateTaskInteractor) UpdateStatus(ctx context.Context, req dto.UpdateTaskStatusRequest) (dto.UpdateTaskStatusResponse, error) {

	err := i.taskUpdaterRepo.UpdateStatus(ctx, req.ID, req.Status)
	if err != nil {
		return dto.UpdateTaskStatusResponse{}, err
	}
	timeNow := time.Now()
	return dto.UpdateTaskStatusResponse{
		ID:        req.ID,
		Status:    req.Status,
		UpdatedAt: &timeNow,
	}, nil
}

func (i *updateTaskInteractor) UpdateDeadline(ctx context.Context, req dto.UpdateTaskDeadlineRequest) (dto.UpdateTaskDeadlineResponse, error) {

	err := i.taskUpdaterRepo.UpdateDeadline(ctx, req.ID, *req.Deadline)
	if err != nil {
		return dto.UpdateTaskDeadlineResponse{}, err
	}

	timeNow := time.Now()
	return dto.UpdateTaskDeadlineResponse{
		ID:        req.ID,
		UpdatedAt: &timeNow,
	}, nil
}
