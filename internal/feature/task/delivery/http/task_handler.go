package http

import (
	"Test/internal/feature/task/interface_adapters/dto"
	"Test/internal/feature/task/usecase"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

type TaskHandler struct {
	createTaskUC usecase.CreateTask
	getTaskUC    usecase.GetTask
	updateTaskUC usecase.UpdateTask
	deleteTaskUC usecase.DeleteTask
}

func NewTaskHandler(
	createTaskUC usecase.CreateTask,
	getTaskUC usecase.GetTask,
	updateTaskUC usecase.UpdateTask,
	deleteTaskUC usecase.DeleteTask,
) *TaskHandler {

	return &TaskHandler{
		createTaskUC: createTaskUC,
		getTaskUC:    getTaskUC,
		updateTaskUC: updateTaskUC,
		deleteTaskUC: deleteTaskUC,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		log.Warn().Err(err).Send()
		return
	}
	organizationID, ok := c.Get("organization_id")
	if !ok || organizationID.(uuid.UUID) == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	req.OrganizationID = organizationID.(uuid.UUID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	response, err := h.createTaskUC.Create(ctx, req)
	if err != nil {
		log.Info().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	req := dto.GetTaskRequest{ID: id}
	response, err := h.getTaskUC.Get(ctx, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func parseUint(str string) (uint, error) {
	res, err := strconv.ParseUint(str, 10, 64)
	return uint(res), err
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	var req dto.GetTasksRequest
	if err := req.GetParameters(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	response, err := h.getTaskUC.GetTasks(ctx, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	var req dto.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	response, err := h.updateTaskUC.UpdateStatus(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) UpdateTaskDeadline(c *gin.Context) {
	var req dto.UpdateTaskDeadlineRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	response, err := h.updateTaskUC.UpdateDeadline(ctx, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	response, err := h.updateTaskUC.Update(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	var (
		req dto.DeleteTaskRequest
		err error
	)

	idStr := c.Param("id")
	req.ID, err = uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	response, err := h.deleteTaskUC.Delete(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}
