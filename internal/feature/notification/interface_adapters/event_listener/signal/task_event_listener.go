package signal

import (
	entity2 "Test/internal/feature/notification/entity"
	"Test/internal/feature/notification/interface_adapters/mapper"
	"Test/internal/feature/notification/usecase"
	"Test/internal/feature/task/entity"
	notification_const2 "Test/internal/pkg/const"
	"Test/internal/pkg/events"
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

type TaskListenerInteractor struct {
	uc usecase.NotificationUseCase
	//bus *events.EventBus
}

func NewTaskListenerInteractor(uc usecase.NotificationUseCase, bus *events.EventBus) *TaskListenerInteractor {
	return &TaskListenerInteractor{uc: uc, bus: bus}
}

func (i *TaskListenerInteractor) Listen(ctx context.Context) {

	entity.TaskCreatedSignal.AddListener(func(ctx context.Context, task *entity.Task) {
		i.handleTask(ctx, task, task.Description, i.bus.Events, i.bus.Errors)
	})

	entity.TaskUpdatedDeadlineSignal.AddListener(func(ctx context.Context, task *entity.Task) {
		i.handleTask(ctx, task, "Deadline updated", i.bus.Events, i.bus.Errors)
	})

	entity.TaskUpdatedSignal.AddListener(func(ctx context.Context, task *entity.Task) {
		i.handleTask(ctx, task, "Task updated", i.bus.Events, i.bus.Errors)
	})
}

func (i *TaskListenerInteractor) Close() error {
	return nil
}

func (i *TaskListenerInteractor) handleTask(
	ctx context.Context,
	task *entity.Task,
	notificationText string,
	evCh chan<- *events.Event,
	errCh chan<- *events.EventError,
) {
	resType := notification_const2.ResourceTask

	resp, err := i.uc.CreateAndDeliver(ctx, &entity2.Notification{
		SendToApp:         true,
		SendToEmail:       true,
		OrganizationID:    task.OrganizationID,
		UserID:            task.AssigneeID,
		Type:              notification_const2.NotificationTask,
		Text:              notificationText,
		RelatedResourceID: &task.ID,
		ResourceType:      &resType,
	})
	if err != nil {
		log.Warn().
			Str("event_listener", "handleTask").
			Str("method", "CreateAndDeliver").
			Err(err).
			Send()
		errCh <- events.NewEventError(err, task.ID.String())
		return
	}

	protoData, err := mapper.NotificationEntityToProto(resp).ProtoData()

	if err != nil {
		errCh <- events.NewEventError(err, task.AssigneeID.String())
		return
	}

	evCh <- &events.Event{
		Data:   protoData,
		Source: "task-listener",
		Time:   time.Now(),
		Headers: map[string]string{
			"user_id": resp.UserID.String(),
		},
	}

}
