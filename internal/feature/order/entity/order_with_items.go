package entity

type OrderWithItems struct {
	Order *Order
	Items []*OrderItem
}
