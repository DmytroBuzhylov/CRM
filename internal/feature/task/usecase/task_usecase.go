package usecase

import (
	"Test/internal/feature/task/interface_adapters/dto"
	"context"
)

type TaskUseCase interface {
	Create(ctx context.Context, req dto.CreateTaskRequest) (dto.CreateTaskResponse, error)
	Get(ctx context.Context, req dto.GetTaskRequest) (dto.GetTaskResponse, error)
	GetTasks(ctx context.Context, req dto.GetTasksRequest) (dto.GetTasksResponse, error)
	Update(ctx context.Context, req dto.UpdateTaskRequest) (dto.UpdateTaskResponse, error)
	UpdateStatus(ctx context.Context, req dto.UpdateTaskStatusRequest) (dto.UpdateTaskStatusResponse, error)
	UpdateDeadline(ctx context.Context, req dto.UpdateTaskDeadlineRequest) (dto.UpdateTaskDeadlineResponse, error)
	Delete(ctx context.Context, req dto.DeleteTaskRequest) (dto.DeleteTaskResponse, error)
}
