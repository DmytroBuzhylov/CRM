package usecase

import (
	"Test/internal/feature/task/interface_adapters/dto"
	"context"
)

type CreateTask interface {
	Create(ctx context.Context, req dto.CreateTaskRequest) (dto.CreateTaskResponse, error)
}

type GetTask interface {
	Get(ctx context.Context, req dto.GetTaskRequest) (dto.GetTaskResponse, error)
	GetTasks(ctx context.Context, req dto.GetTasksRequest) (dto.GetTasksResponse, error)
}

type UpdateTask interface {
	Update(ctx context.Context, req dto.UpdateTaskRequest) (dto.UpdateTaskResponse, error)
	UpdateStatus(ctx context.Context, req dto.UpdateTaskStatusRequest) (dto.UpdateTaskStatusResponse, error)
	UpdateDeadline(ctx context.Context, req dto.UpdateTaskDeadlineRequest) (dto.UpdateTaskDeadlineResponse, error)
}

type DeleteTask interface {
	Delete(ctx context.Context, req dto.DeleteTaskRequest) (dto.DeleteTaskResponse, error)
}
