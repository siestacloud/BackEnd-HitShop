package models

import "github.com/google/uuid"

// статусы заказов
type OrdersStatuses struct {
	PkOrderStatusId        uuid.UUID `gorm:"column:pk_order_status_id;					type:uuid;	default:uuid_generate_v4();	primary_key"`
	OrderStatusName        string    `gorm:"column:order_status_name;						type:varchar(30);"`
	OrderStatusDescription string    `gorm:"column:order_status_description;		type:varchar(255);"`
}

// заказы клиетов с указанием на филиал и статус заказа
type Orders struct {
	PkOrderId uuid.UUID `gorm:"column:pk_order_id;	type:uuid;	default:uuid_generate_v4();	primary_key"`

	FkStoreIdReferer int
	StoreId          Stores `gorm:"column:fk_store_id_referer;	foreignKey:FkStoreIdReferer;not null"`

	FkOrderStatusIdReferer int
	OrderStatusId          OrdersStatuses `gorm:"column:fk_order_status_id_referer;	foreignKey:FkOrderStatusIdReferer;not null"`
}
