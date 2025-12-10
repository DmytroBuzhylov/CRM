package usecase

import (
	"Test/internal/feature/staff/repository"
	repository2 "Test/internal/feature/task/repository"
	system_const "Test/internal/pkg/const"
	"context"

	"github.com/google/uuid"
)

type TaskDispatcherUseCaseInteractor struct {
	staffRepo repository.StaffRepository
	taskRepo  repository2.TaskRepository
}

func NewTaskDispatcherUseCaseInteractor() *TaskDispatcherUseCaseInteractor {
	return &TaskDispatcherUseCaseInteractor{}
}

func (i *TaskDispatcherUseCaseInteractor) FindAndAssign(ctx context.Context, orgID uuid.UUID, taskType system_const.TaskType) (uuid.UUID, error) {
	// find less busy staff()
	// assign and notify staff()

	return uuid.Nil, nil
}

func (i *TaskDispatcherUseCaseInteractor) FindAvailableStaffByPermission(ctx context.Context, orgID uuid.UUID, requiredPermission string) ([]uuid.UUID, error) {
	i.staffRepo.
}
