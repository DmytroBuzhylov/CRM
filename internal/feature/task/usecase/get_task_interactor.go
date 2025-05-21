package usecase

import (
	"Test/internal/feature/task/interface_adapters/dto"
	"Test/internal/feature/task/repository"
	"context"
)

type getTaskInteractor struct {
	taskGetterRepo repository.TaskGetter
}

func NewGetTaskInteractor(taskGetterRepo repository.TaskGetter) *getTaskInteractor {
	return &getTaskInteractor{taskGetterRepo: taskGetterRepo}
}

func (i *getTaskInteractor) Get(ctx context.Context, req dto.GetTaskRequest) (dto.GetTaskResponse, error) {
	gotTask, err := i.taskGetterRepo.Get(ctx, req.ID)
	if err != nil {
		return dto.GetTaskResponse{}, err
	}

	return dto.GetTaskResponse{
		ID:          gotTask.ID,
		Name:        gotTask.Name,
		Description: gotTask.Description,
		Priority:    gotTask.Priority,
		Status:      gotTask.Status,
		Deadline:    gotTask.Deadline,
		AssigneeID:  gotTask.AssigneeID,
		ClientID:    gotTask.ClientID,
		CreatedAt:   gotTask.CreatedAt,
		UpdatedAt:   gotTask.UpdatedAt,
	}, nil
}

func (i *getTaskInteractor) GetTasks(ctx context.Context, req dto.GetTasksRequest) (dto.GetTasksResponse, error) {
	filter := repository.TaskFilter{
		Name:       req.Filter.Name,
		Priority:   req.Filter.Priority,
		Status:     req.Filter.Status,
		Deadline:   req.Filter.Deadline,
		AssigneeID: req.Filter.AssigneeID,
		ClientID:   req.Filter.ClientID,
		CreatedAt:  req.Filter.CreatedAt,
		UpdatedAt:  req.Filter.UpdatedAt,
	}
	pagination := repository.Pagination{
		Page:  req.Page,
		Limit: 0,
	}

	gotTasks, err := i.taskGetterRepo.GetTasks(ctx, filter, pagination)
	if err != nil {
		return dto.GetTasksResponse{}, err
	}
	tasks := make([]dto.GetTaskResponse, len(gotTasks))
	for i, task := range gotTasks {
		tasks[i] = dto.GetTaskResponse{
			ID:          task.ID,
			Name:        task.Name,
			Description: task.Description,
		}
	}
	return dto.GetTasksResponse{Tasks: tasks}, nil
}
