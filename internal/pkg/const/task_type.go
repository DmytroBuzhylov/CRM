package system_const

type TaskType string

const (
	TaskOrderCreated   TaskType = "TaskOrderCreated"
	TaskOrderConfirmed TaskType = "TaskOrderConfirmed"
	TaskOrderUpdated   TaskType = "TaskOrderUpdated"
	TaskOrderCanceled  TaskType = "TaskOrderCanceled"
)
