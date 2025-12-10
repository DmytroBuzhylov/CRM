package system_const

type NotificationType string

const (
	NotificationUnspecified NotificationType = "UNSPECIFIED"
	NotificationOrder       NotificationType = "ORDER"
	NotificationTask        NotificationType = "TASK"
	NotificationSystem      NotificationType = "SYSTEM"

	NotificationClientRequest NotificationType = "CLIENT_REQUEST"

	NotificationShiftStarted NotificationType = "SHIFT_STARTED"
	NotificationShiftEnded   NotificationType = "SHIFT_ENDED"
	NotificationShiftOverdue NotificationType = "SHIFT_OVERDUE"

	NotificationFoodReady NotificationType = "FOOD_READY"

	NotificationIncident       NotificationType = "INCIDENT"
	NotificationEquipmentAlert NotificationType = "EQUIPMENT_ALERT"

	NotificationScreenMessage NotificationType = "SCREEN_MESSAGE"
	NotificationManagerAlert  NotificationType = "MANAGER_ALERT"
)
