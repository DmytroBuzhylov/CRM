package postgres

import (
	"Test/internal/feature/task/entity"
	"Test/internal/feature/task/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task entity.Task) (entity.Task, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM organizations WHERE id = $1)`, task.OrganizationID).Scan(&exists)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to find organization: %w", err)
	}
	if !exists {
		return entity.Task{}, fmt.Errorf("invalid organization id")
	}
	if task.AssigneeID != uuid.Nil {
		err = r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM organizations_users WHERE user_id = $1)`, task.AssigneeID).Scan(&exists)
		if err != nil {
			return entity.Task{}, fmt.Errorf("failed to find assignee: %w", err)
		}
		if !exists {
			return entity.Task{}, fmt.Errorf("invalid assignee id")
		}
	}

	query := `INSERT INTO tasks (id, name, organization_id, description, priority, status, deadline, assignee_id, client_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = r.db.Exec(ctx, query, task.ID, task.Name, task.OrganizationID, task.Description, task.Priority, task.Status, task.Deadline, task.AssigneeID, task.ClientID, task.ID, task.CreatedAt)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to create task: %w", err)
	}
	return task, nil
}

func (r *TaskRepository) Get(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	query := `SELECT id, name, description, priority, status, deadline, assignee_id, client_id, created_at, updated_at FROM tasks WHERE id = $1`
	var task entity.Task
	err := r.db.QueryRow(ctx, query, id).Scan(&task.ID, &task.Name, &task.Description, &task.Priority, &task.Status, &task.Deadline, &task.AssigneeID, &task.ClientID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

func (r *TaskRepository) GetTasks(ctx context.Context, filter repository.TaskFilter, pagination repository.Pagination) ([]entity.Task, error) {
	query := `SELECT id, name, description, priority, status, deadline, assignee_id, client_id, created_at, updated_at FROM tasks WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name LIKE $%d", argCount)
		args = append(args, "%"+filter.Name+"%")
		argCount++
	}
	if filter.Priority > 0 {
		query += fmt.Sprintf(" AND priority = $%d", argCount)
		args = append(args, filter.Priority)
		argCount++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filter.Status)
		argCount++
	}
	if filter.Deadline != nil {
		query += fmt.Sprintf(" AND deadline::date = $%d::date", argCount)
		args = append(args, *filter.Deadline)
		argCount++
	}
	if filter.AssigneeID != uuid.Nil {
		query += fmt.Sprintf(" AND assignee_id = $%d", argCount)
		args = append(args, filter.AssigneeID)
		argCount++
	}
	if filter.ClientID != uuid.Nil {
		query += fmt.Sprintf(" AND client_id = $%d", argCount)
		args = append(args, filter.ClientID)
		argCount++
	}
	if filter.CreatedAt != nil {
		query += fmt.Sprintf(" AND created_at::date = $%d::date", argCount)
		args = append(args, *filter.CreatedAt)
		argCount++
	}
	if filter.UpdatedAt != nil {
		query += fmt.Sprintf(" AND updated_at::date = $%d::date", argCount)
		args = append(args, *filter.UpdatedAt)
		argCount++
	}

	if pagination.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
		args = append(args, pagination.Limit, (pagination.Page-1)*pagination.Limit)
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Priority, &task.Status, &task.Deadline, &task.AssigneeID, &task.ClientID, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task entity.Task) (entity.Task, error) {
	query := `UPDATE tasks SET name = $2, description = $3, priority = $4, status = $5, deadline = $6, assignee_id = $7, client_id = $8, updated_at = $9 WHERE id = $1 RETURNING updated_at`
	err := r.db.QueryRow(ctx, query, task.ID, task.Name, task.Description, task.Priority, task.Status, task.Deadline, task.AssigneeID, task.ClientID, time.Now()).Scan(&task.UpdatedAt)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to update task: %w", err)
	}
	return task, nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE tasks SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}

func (r *TaskRepository) UpdateDeadline(ctx context.Context, id uuid.UUID, deadline time.Time) error {
	query := `UPDATE tasks SET deadline = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, deadline, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update task deadline: %w", err)
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
