package postgres

import (
	"Test/internal/feature/task/entity"
	"Test/internal/feature/task/repository"
	"Test/internal/pkg/provider"
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepo struct {
	db        *pgxpool.Pool
	connector provider.SQLConnector
	qb        sq.StatementBuilderType
}

func NewTaskRepository(db *pgxpool.Pool, connector provider.SQLConnector) *TaskRepo {
	return &TaskRepo{
		db:        db,
		connector: connector,
		qb:        sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

//func (r *TaskRepo) Create(ctx context.Context, task entity.Task) (entity.Task, error) {
//	var exists bool
//	err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM organizations WHERE id = $1)`, task.OrganizationID).Scan(&exists)
//	if err != nil {
//		return entity.Task{}, fmt.Errorf("failed to find organization: %w", err)
//	}
//	if !exists {
//		return entity.Task{}, fmt.Errorf("invalid organization id")
//	}
//	if task.AssigneeID != uuid.Nil {
//		err = r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM organizations_users WHERE user_id = $1)`, task.AssigneeID).Scan(&exists)
//		if err != nil {
//			return entity.Task{}, fmt.Errorf("failed to find assignee: %w", err)
//		}
//		if !exists {
//			return entity.Task{}, fmt.Errorf("invalid assignee id")
//		}
//	}
//
//	query := `INSERT INTO tasks (id, name, organization_id, description, priority, status, deadline, assignee_id, client_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
//	_, err = r.db.Exec(ctx, query, task.ID, task.Name, task.OrganizationID, task.Description, task.Priority, task.Status, task.Deadline, task.AssigneeID, task.ClientID, task.CreatedAt, task.UpdatedAt)
//	if err != nil {
//		return entity.Task{}, fmt.Errorf("failed to create task: %w", err)
//	}
//	return task, nil
//}

func (r *TaskRepo) Create(ctx context.Context, task entity.Task) (entity.Task, error) {
	query, _, err := r.qb.Select("SELECT EXISTS (SELECT 1 FROM organizations WHERE id = $1)").
		From("organizations").
		Limit(1).
		ToSql()
	if err != nil {
		return entity.Task{}, err
	}

	var exists bool
	err = r.connector.QueryRow(ctx, query, task.OrganizationID).Scan(&exists)
	if err != nil {
		return entity.Task{}, err
	}
	if !exists {
		return entity.Task{}, fmt.Errorf("invalid assignee id")
	}

	query, args, err := r.qb.Insert("tasks").
		Columns("id", "name", "organization_id", "description", "priority", "status", "deadline", "assignee_id", "client_id", "created_at", "updated_at").
		Values(task.ID, task.Name, task.OrganizationID, task.Description, task.Priority, task.Status, task.Deadline, task.AssigneeID, task.ClientID, task.CreatedAt, task.UpdatedAt).
		ToSql()
	if err != nil {
		return entity.Task{}, err
	}

	_, err = r.connector.Exec(ctx, query, args...)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *TaskRepo) Get(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	var task entity.Task
	query, args, err := r.qb.Select("id", "name", "description", "priority", "status", "deadline", "assignee_id", "client_id", "created_at", "updated_at").
		From("tasks").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Task{}, err
	}

	err = r.connector.QueryRow(ctx, query, args...).
		Scan(&task.ID, &task.Name, &task.Description, &task.Priority, &task.Status, &task.Deadline, &task.AssigneeID, &task.ClientID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *TaskRepo) GetTasks(ctx context.Context, filter repository.TaskFilter, pagination repository.Pagination) ([]entity.Task, error) {
	builder := r.qb.Select(
		"id",
		"name",
		"description",
		"priority",
		"status",
		"deadline",
		"assignee_id",
		"client_id",
		"created_at",
		"updated_at",
	).From("tasks")

	if filter.Name != "" {
		builder = builder.Where(sq.ILike{"name": "%" + filter.Name + "%"})
	}
	if filter.Priority > 0 {
		builder = builder.Where(sq.Eq{"priority": filter.Priority})
	}
	if filter.Status != "" {
		builder = builder.Where(sq.Eq{"status": filter.Status})
	}
	if filter.Deadline != nil {
		builder = builder.Where(sq.Expr("deadline::date = ?", *filter.Deadline))
	}
	if filter.AssigneeID != uuid.Nil {
		builder = builder.Where(sq.Eq{"assignee_id": filter.AssigneeID})
	}
	if filter.ClientID != uuid.Nil {
		builder = builder.Where(sq.Eq{"client_id": filter.ClientID})
	}
	if filter.CreatedAt != nil {
		builder = builder.Where(sq.Expr("created_at::date = ?", *filter.CreatedAt))
	}
	if filter.UpdatedAt != nil {
		builder = builder.Where(sq.Expr("updated_at::date = ?", *filter.UpdatedAt))
	}

	if pagination.Limit > 0 {
		offset := (pagination.Page - 1) * pagination.Limit
		builder = builder.Limit(uint64(pagination.Limit)).Offset(uint64(offset))
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.Deadline,
			&task.AssigneeID,
			&task.ClientID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepo) Update(ctx context.Context, task entity.Task) (entity.Task, error) {
	query := `UPDATE tasks SET name = $2, description = $3, priority = $4, status = $5, deadline = $6, assignee_id = $7, client_id = $8, updated_at = $9 WHERE id = $1 RETURNING updated_at`
	err := r.db.QueryRow(ctx, query, task.ID, task.Name, task.Description, task.Priority, task.Status, task.Deadline, task.AssigneeID, task.ClientID, time.Now()).Scan(&task.UpdatedAt)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to update task: %w", err)
	}
	return task, nil
}

func (r *TaskRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE tasks SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}

func (r *TaskRepo) UpdateDeadline(ctx context.Context, id uuid.UUID, deadline time.Time) error {
	query := `UPDATE tasks SET deadline = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, deadline, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update task deadline: %w", err)
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
