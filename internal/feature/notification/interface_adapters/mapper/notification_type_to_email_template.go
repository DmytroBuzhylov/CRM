package mapper

import (
	"Test/internal/pkg/const"
	"errors"
)

func PrepareEmailContent(notType notification_const.NotificationType) (template string, subject string, err error) {
	switch notType {
	case notification_const.NotificationUnspecified:
		template, subject, err = "", "", nil
	case notification_const.NotificationOrder:
		template, subject, err = "order.html", "Уведомление о заказе", nil

	case notification_const.NotificationTask:
		template, subject, err = "task.html", "Уведомление о задаче", nil
	case notification_const.NotificationSystem:
		template, subject, err = "system.html", "Системное уведомление", nil
	case notification_const.NotificationClientRequest:
		template, subject, err = "client_request.html", "Новый запрос от клиента", nil

	case notification_const.NotificationShiftStarted:
		template, subject, err = "shift_started.html", "Смена началась", nil

	case notification_const.NotificationShiftEnded:
		template, subject, err = "shift_ended.html", "Смена завершена", nil

	case notification_const.NotificationShiftOverdue:
		template, subject, err = "shift_overdue", "Внимание: просроченная смена", nil

	case notification_const.NotificationFoodReady:
		template, subject, err = "food_ready.html", "Ваш заказ готов!", nil

	case notification_const.NotificationIncident:
		template, subject, err = "incident.html", "Уведомление об инциденте", nil

	case notification_const.NotificationEquipmentAlert:
		template, subject, err = "equipment_alert.html", "Оповещение оборудования", nil

	case notification_const.NotificationScreenMessage:
		template, subject, err = "screen_message.html", "Новое сообщение на экране", nil

	case notification_const.NotificationManagerAlert:
		template, subject, err = "manager_alert.html", "Срочное оповещение для менеджера", nil

	default:
		template, subject, err = "", "", errors.New("incorrect notification type")

	}
	return template, subject, err
}
