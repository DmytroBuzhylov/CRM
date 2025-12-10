package usecase

import (
	"Test/internal/feature/task/entity"
	"Test/internal/feature/task/interface_adapters/dto"
	"Test/internal/feature/task/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskUseCaseInteractor struct {
	repo repository.TaskRepository
}

func NewTaskUseCaseInteractor(repo repository.TaskRepository) *TaskUseCaseInteractor {
	return &TaskUseCaseInteractor{repo: repo}
}

func (i *TaskUseCaseInteractor) Create(ctx context.Context, req dto.CreateTaskRequest) (dto.CreateTaskResponse, error) {
	timeNow := time.Now()
	newTask := entity.Task{
		ID:             uuid.New(),
		Name:           req.Name,
		OrganizationID: req.OrganizationID,
		Description:    req.Description,
		Priority:       req.Priority,
		Status:         "New",
		Deadline:       req.Deadline,
		AssigneeID:     req.AssigneeID,
		ClientID:       req.ClientID,
		CreatedAt:      &timeNow,
		UpdatedAt:      &timeNow,
	}

	createdTask, err := i.repo.Create(ctx, newTask)
	if err != nil {
		return dto.CreateTaskResponse{}, err
	}

	entity.TaskCreatedSignal.Emit(ctx, &newTask)

	return dto.CreateTaskResponse{
		ID:        createdTask.ID,
		CreatedAt: createdTask.CreatedAt,
	}, nil
}

func (i *TaskUseCaseInteractor) Delete(ctx context.Context, req dto.DeleteTaskRequest) (dto.DeleteTaskResponse, error) {
	err := i.repo.Delete(ctx, req.ID)
	if err != nil {
		return dto.DeleteTaskResponse{
			Status: false,
		}, err
	}

	return dto.DeleteTaskResponse{
		Status: true,
	}, nil
}

func (i *TaskUseCaseInteractor) Get(ctx context.Context, req dto.GetTaskRequest) (dto.GetTaskResponse, error) {
	gotTask, err := i.repo.Get(ctx, req.ID)
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

func (i *TaskUseCaseInteractor) GetTasks(ctx context.Context, req dto.GetTasksRequest) (dto.GetTasksResponse, error) {
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

	gotTasks, err := i.repo.GetTasks(ctx, filter, pagination)
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

func (i *TaskUseCaseInteractor) Update(ctx context.Context, req dto.UpdateTaskRequest) (dto.UpdateTaskResponse, error) {
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

	updatedTask, err := i.repo.Update(ctx, newTask)
	if err != nil {
		return dto.UpdateTaskResponse{}, err
	}

	entity.TaskUpdatedSignal.Emit(ctx, &newTask)
	return dto.UpdateTaskResponse{
		ID:        updatedTask.ID,
		Status:    updatedTask.Status,
		UpdatedAt: updatedTask.UpdatedAt,
	}, nil
}

func (i *TaskUseCaseInteractor) UpdateStatus(ctx context.Context, req dto.UpdateTaskStatusRequest) (dto.UpdateTaskStatusResponse, error) {

	err := i.repo.UpdateStatus(ctx, req.ID, req.Status)
	if err != nil {
		return dto.UpdateTaskStatusResponse{}, err
	}
	timeNow := time.Now()

	entity.TaskUpdatedStatusSignal.Emit(ctx, &entity.Task{
		ID:     req.ID,
		Status: req.Status,
	})

	return dto.UpdateTaskStatusResponse{
		ID:        req.ID,
		Status:    req.Status,
		UpdatedAt: &timeNow,
	}, nil
}

func (i *TaskUseCaseInteractor) UpdateDeadline(ctx context.Context, req dto.UpdateTaskDeadlineRequest) (dto.UpdateTaskDeadlineResponse, error) {

	err := i.repo.UpdateDeadline(ctx, req.ID, *req.Deadline)
	if err != nil {
		return dto.UpdateTaskDeadlineResponse{}, err
	}

	timeNow := time.Now()

	entity.TaskUpdatedDeadlineSignal.Emit(ctx, &entity.Task{
		ID:       req.ID,
		Deadline: req.Deadline,
	})

	return dto.UpdateTaskDeadlineResponse{
		ID:        req.ID,
		UpdatedAt: &timeNow,
	}, nil
}

func (i *TaskUseCaseInteractor) CreateAutomaticTask(ctx context.Context, task *entity.Task) error {
	_, err := i.repo.Create(ctx, *task)
	return err
}
