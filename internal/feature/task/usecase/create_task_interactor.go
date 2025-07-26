package usecase

import (
	"Test/internal/feature/task/entity"
	"Test/internal/feature/task/interface_adapters/dto"
	"Test/internal/feature/task/repository"
	"context"
	"time"
)

type createTaskInteractor struct {
	taskCreatorRepo repository.TaskCreator
}

func NewCreateTaskInteractor(taskCreatorRepo repository.TaskCreator) *createTaskInteractor {
	return &createTaskInteractor{taskCreatorRepo: taskCreatorRepo}
}

func (i *createTaskInteractor) Create(ctx context.Context, req dto.CreateTaskRequest) (dto.CreateTaskResponse, error) {
	timeNow := time.Now()
	newTask := entity.Task{
		Name:           req.Name,
		OrganizationID: *req.OrganizationID,
		Description:    req.Description,
		Priority:       req.Priority,
		Status:         "New",
		Deadline:       req.Deadline,
		AssigneeID:     req.AssigneeID,
		ClientID:       req.ClientID,
		CreatedAt:      &timeNow,
	}
	
	createdTask, err := i.taskCreatorRepo.Create(ctx, newTask)
	if err != nil {
		return dto.CreateTaskResponse{}, err
	}

	return dto.CreateTaskResponse{
		ID:        createdTask.ID,
		CreatedAt: createdTask.CreatedAt,
	}, nil
}
