package utils

import "fmt"

type NotificationType struct {
	value string
}

func (t *NotificationType) String() string {
	return t.value
}

func (t *NotificationType) IsSystem() bool {
	return t.value == "SYSTEM"
}

func NewNotificationType(value string) (NotificationType, error) {
	switch value {
	case "UNSPECIFIED", "ORDER", "TASK", "SYSTEM", "CLIENT_REQUEST",
		"SHIFT_STARTED", "SHIFT_ENDED", "SHIFT_OVERDUE", "FOOD_READY",
		"INCIDENT", "EQUIPMENT_ALERT", "SCREEN_MESSAGE", "MANAGER_ALERT":
		return NotificationType{value: value}, nil
	default:
		return NotificationType{value: "UNSPECIFIED"}, fmt.Errorf("unknown type: %s", value)
	}
}
