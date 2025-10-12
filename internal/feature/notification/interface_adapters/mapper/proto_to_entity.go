package mapper

import (
	"Test/internal/feature/notification/entity"
	notification_const "Test/internal/feature/notification/entity/const"
	pb "Test/internal/feature/notification/interface_adapters/proto/gen"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func PBResourceTypeToEntity(pbType pb.ResourceType) (notification_const.ResourceType, error) {
	switch pbType {
	case pb.ResourceType_RESOURCE_TYPE_ORDER:
		return notification_const.ResourceOrder, nil
	case pb.ResourceType_RESOURCE_TYPE_UNSPECIFIED:
		return notification_const.ResourceUnspecified, nil
	case pb.ResourceType_RESOURCE_TYPE_TABLE:
		return notification_const.ResourceTable, nil
	case pb.ResourceType_RESOURCE_TYPE_TASK:
		return notification_const.ResourceTask, nil
	case pb.ResourceType_RESOURCE_TYPE_MENU_ITEM:
		return notification_const.ResourceMenuItem, nil
	case pb.ResourceType_RESOURCE_TYPE_INVENTORY:
		return notification_const.ResourceInventory, nil
	case pb.ResourceType_RESOURCE_TYPE_SHIFT:
		return notification_const.ResourceShift, nil
	case pb.ResourceType_RESOURCE_TYPE_EMPLOYEE:
		return notification_const.ResourceEmployee, nil
	case pb.ResourceType_RESOURCE_TYPE_KITCHEN_STATION:
		return notification_const.ResourceKitchenStation, nil
	case pb.ResourceType_RESOURCE_TYPE_INCIDENT:
		return notification_const.ResourceIncident, nil

	default:
		return "", fmt.Errorf("the specified type does not exist")
	}
}

func ResourceEntityToPBType(resType notification_const.ResourceType) (pb.ResourceType, error) {
	switch resType {
	case notification_const.ResourceOrder:
		return pb.ResourceType_RESOURCE_TYPE_ORDER, nil
	case notification_const.ResourceUnspecified:
		return pb.ResourceType_RESOURCE_TYPE_UNSPECIFIED, nil
	case notification_const.ResourceTable:
		return pb.ResourceType_RESOURCE_TYPE_TABLE, nil
	case notification_const.ResourceTask:
		return pb.ResourceType_RESOURCE_TYPE_TASK, nil
	case notification_const.ResourceMenuItem:
		return pb.ResourceType_RESOURCE_TYPE_MENU_ITEM, nil
	case notification_const.ResourceInventory:
		return pb.ResourceType_RESOURCE_TYPE_INVENTORY, nil
	case notification_const.ResourceShift:
		return pb.ResourceType_RESOURCE_TYPE_SHIFT, nil
	case notification_const.ResourceEmployee:
		return pb.ResourceType_RESOURCE_TYPE_EMPLOYEE, nil
	case notification_const.ResourceKitchenStation:
		return pb.ResourceType_RESOURCE_TYPE_KITCHEN_STATION, nil
	case notification_const.ResourceIncident:
		return pb.ResourceType_RESOURCE_TYPE_INCIDENT, nil

	default:
		return -1, fmt.Errorf("the specified type does not exist")
	}
}

func PBNotificationTypeToEntity(pbType pb.NotificationType) (notification_const.NotificationType, error) {
	switch pbType {
	case pb.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED:
		return notification_const.NotificationUnspecified, nil

	case pb.NotificationType_NOTIFICATION_TYPE_ORDER:
		return notification_const.NotificationOrder, nil

	case pb.NotificationType_NOTIFICATION_TYPE_TASK:
		return notification_const.NotificationTask, nil

	case pb.NotificationType_NOTIFICATION_TYPE_SYSTEM:
		return notification_const.NotificationSystem, nil

	case pb.NotificationType_NOTIFICATION_TYPE_CLIENT_REQUEST:
		return notification_const.NotificationClientRequest, nil

	case pb.NotificationType_NOTIFICATION_TYPE_SHIFT_STARTED:
		return notification_const.NotificationShiftStarted, nil

	case pb.NotificationType_NOTIFICATION_TYPE_SHIFT_ENDED:
		return notification_const.NotificationShiftEnded, nil

	case pb.NotificationType_NOTIFICATION_TYPE_SHIFT_OVERDUE:
		return notification_const.NotificationShiftOverdue, nil

	case pb.NotificationType_NOTIFICATION_TYPE_FOOD_READY:
		return notification_const.NotificationFoodReady, nil

	case pb.NotificationType_NOTIFICATION_TYPE_INCIDENT:
		return notification_const.NotificationIncident, nil

	case pb.NotificationType_NOTIFICATION_TYPE_EQUIPMENT_ALERT:
		return notification_const.NotificationEquipmentAlert, nil

	case pb.NotificationType_NOTIFICATION_TYPE_SCREEN_MESSAGE:
		return notification_const.NotificationScreenMessage, nil

	case pb.NotificationType_NOTIFICATION_TYPE_MANAGER_ALERT:
		return notification_const.NotificationManagerAlert, nil

	default:
		return "", fmt.Errorf("the specified type does not exist")
	}
}

func NotificationEntityToPBType(notType notification_const.NotificationType) (pb.NotificationType, error) {
	switch notType {
	case notification_const.NotificationUnspecified:
		return pb.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED, nil

	case notification_const.NotificationOrder:
		return pb.NotificationType_NOTIFICATION_TYPE_ORDER, nil

	case notification_const.NotificationTask:
		return pb.NotificationType_NOTIFICATION_TYPE_TASK, nil

	case notification_const.NotificationSystem:
		return pb.NotificationType_NOTIFICATION_TYPE_SYSTEM, nil

	case notification_const.NotificationClientRequest:
		return pb.NotificationType_NOTIFICATION_TYPE_CLIENT_REQUEST, nil

	case notification_const.NotificationShiftStarted:
		return pb.NotificationType_NOTIFICATION_TYPE_SHIFT_STARTED, nil

	case notification_const.NotificationShiftEnded:
		return pb.NotificationType_NOTIFICATION_TYPE_SHIFT_ENDED, nil

	case notification_const.NotificationShiftOverdue:
		return pb.NotificationType_NOTIFICATION_TYPE_SHIFT_OVERDUE, nil

	case notification_const.NotificationFoodReady:
		return pb.NotificationType_NOTIFICATION_TYPE_FOOD_READY, nil

	case notification_const.NotificationIncident:
		return pb.NotificationType_NOTIFICATION_TYPE_INCIDENT, nil

	case notification_const.NotificationEquipmentAlert:
		return pb.NotificationType_NOTIFICATION_TYPE_EQUIPMENT_ALERT, nil

	case notification_const.NotificationScreenMessage:
		return pb.NotificationType_NOTIFICATION_TYPE_SCREEN_MESSAGE, nil

	case notification_const.NotificationManagerAlert:
		return pb.NotificationType_NOTIFICATION_TYPE_MANAGER_ALERT, nil

	default:
		return -1, fmt.Errorf("the specified type does not exist")
	}
}

type NotificationProto struct {
	data *pb.Notification
}

func NotificationEntityToProto(n *entity.Notification) *NotificationProto {
	nType, err := NotificationEntityToPBType(n.Type)
	if err != nil {
		nType = pb.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED
	}
	relResId := n.RelatedResourceID.String()
	rType, err := ResourceEntityToPBType(*n.ResourceType)
	if err != nil {
		rType = pb.ResourceType_RESOURCE_TYPE_UNSPECIFIED
	}

	protoData := &pb.Notification{
		Id:                n.ID.String(),
		OrganizationId:    n.OrganizationID.String(),
		UserId:            n.UserID.String(),
		Type:              nType,
		Text:              n.Text,
		Viewed:            n.Viewed,
		RelatedResourceId: &relResId,
		ResourceType:      &rType,
		CreatedAt:         n.CreatedAt,
	}

	return &NotificationProto{data: protoData}
}

func (p *NotificationProto) Struct() *pb.Notification {
	return p.data
}

func (p *NotificationProto) ProtoData() ([]byte, error) {
	return proto.Marshal(p.data)
}
