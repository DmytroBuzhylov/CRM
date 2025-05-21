package usecase

import (
	"Test/internal/feature/task/interface_adapters/dto"
	"Test/internal/feature/task/repository"
	"context"
)

type deleteTaskInteractor struct {
	taskDeleterRepo repository.TaskDeleter
}

func NewDeleteTaskInteractor(taskDeleterRepo repository.TaskDeleter) *deleteTaskInteractor {
	return &deleteTaskInteractor{taskDeleterRepo: taskDeleterRepo}
}

func (i *deleteTaskInteractor) Delete(ctx context.Context, req dto.DeleteTaskRequest) (dto.DeleteTaskResponse, error) {
	err := i.taskDeleterRepo.Delete(ctx, req.ID)
	if err != nil {
		return dto.DeleteTaskResponse{
			Status: false,
		}, err
	}

	return dto.DeleteTaskResponse{
		Status: true,
	}, nil
}
