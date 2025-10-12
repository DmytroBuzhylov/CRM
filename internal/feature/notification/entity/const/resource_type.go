package notification_const

type (
	ResourceType string
)

const (
	ResourceUnspecified    ResourceType = "UNSPECIFIED"
	ResourceOrder          ResourceType = "ORDER"
	ResourceTable          ResourceType = "TABLE"
	ResourceTask           ResourceType = "TASK"
	ResourceInventory      ResourceType = "INVENTORY"
	ResourceMenuItem       ResourceType = "MENU_ITEM"
	ResourceShift          ResourceType = "SHIFT"
	ResourceEmployee       ResourceType = "EMPLOYEE"
	ResourceKitchenStation ResourceType = "KITCHEN_STATION"
	ResourceIncident       ResourceType = "INCIDENT"
)
