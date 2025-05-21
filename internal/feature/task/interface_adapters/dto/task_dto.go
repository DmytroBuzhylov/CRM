package dto

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type CreateTaskRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Priority    uint       `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	AssigneeID  uint64     `json:"assignee_id"`
	ClientID    uint64     `json:"client_id"`
}

type CreateTaskResponse struct {
	ID        uint64     `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
}

type UpdateTaskRequest struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Priority    uint       `json:"priority"`
	Status      string     `json:"status"`
	Deadline    *time.Time `json:"deadline"`
	AssigneeID  uint64     `json:"assignee_id"`
	ClientID    uint64     `json:"client_id"`
}

type UpdateTaskResponse struct {
	ID        uint64     `json:"id"`
	Status    string     `json:"status"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UpdateTaskStatusRequest struct {
	ID     uint64 `json:"id"`
	Status string `json:"status"`
}

type UpdateTaskStatusResponse struct {
	ID        uint64     `json:"id"`
	Status    string     `json:"status"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UpdateTaskDeadlineRequest struct {
	ID       uint64     `json:"id"`
	Deadline *time.Time `json:"deadline"`
}

type UpdateTaskDeadlineResponse struct {
	ID        uint64     `json:"id"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type DeleteTaskRequest struct {
	ID uint64 `json:"id"`
}

type DeleteTaskResponse struct {
	Status bool `json:"status"`
}

type GetTaskRequest struct {
	ID uint64 `json:"id"`
}

type GetTaskResponse struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Priority    uint       `json:"priority"`
	Status      string     `json:"status"`
	Deadline    *time.Time `json:"deadline"`
	AssigneeID  uint64     `json:"assignee_id"`
	ClientID    uint64     `json:"client_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type GetTasksRequest struct {
	Page   uint       `json:"page"`
	Filter FilterTask `json:"filter"`
}

func (r *GetTasksRequest) GetParameters(c *gin.Context) error {
	var err error
	pageStr := c.Query("page")
	if pageStr != "" {
		page, err := strconv.ParseUint(pageStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid page format: %w", err)
		}
		r.Page = uint(page)
	} else {
		r.Page = 1
	}
	r.Filter.Name = c.Query("name")
	r.Filter.Status = c.Query("status")
	priorityStr := c.Query("priority")
	if priorityStr != "" {
		priority, err := strconv.ParseUint(priorityStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid priority format: %w", err)
		}
		r.Filter.Priority = uint(priority)
	}
	timeParse := func(queryParam string) (*time.Time, error) {
		timeStr := c.Query(queryParam)
		if timeStr != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
			if err != nil {
				return nil, fmt.Errorf("invalid %s format: %w", queryParam, err)
			}
			return &parsedTime, nil
		}
		return nil, nil
	}

	r.Filter.Deadline, err = timeParse("deadline")
	if err != nil {
		return err
	}
	r.Filter.CreatedAt, err = timeParse("created_at")
	if err != nil {
		return err
	}
	r.Filter.UpdatedAt, err = timeParse("updated_at")
	if err != nil {
		return err
	}
	r.Filter.ClientID, err = strconv.ParseUint(c.Query("client_id"), 10, 64)
	if err != nil {
		return err
	}
	r.Filter.AssigneeID, err = strconv.ParseUint(c.Query("client_id"), 10, 64)
	if err != nil {
		return err
	}
	return nil
}

type FilterTask struct {
	Name       string     `json:"name"`
	Priority   uint       `json:"priority"`
	Status     string     `json:"status"`
	Deadline   *time.Time `json:"deadline"`
	AssigneeID uint64     `json:"assignee_id"`
	ClientID   uint64     `json:"client_id"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type GetTasksResponse struct {
	Tasks []GetTaskResponse `json:"tasks"`
}

type AcceptTaskRequest struct {
	ID         uint64 `json:"id"`
	AssigneeID uint64 `json:"assignee_id"`
}

type AcceptTaskResponse struct {
	ID     uint64 `json:"id"`
	Status string `json:"status"`
}

type CompletedTaskRequest struct {
	ID uint64 `json:"id"`
}

type CompletedTaskResponse struct {
	ID     uint64 `json:"id"`
	Status string `json:"status"`
}
