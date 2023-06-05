package models

import (
	"time"

	"github.com/google/uuid"
)

// поставки
type Deliveries struct {
	PkDeliveryId uuid.UUID `gorm:"column:pk_delivery_id;	type:uuid;	default:uuid_generate_v4();	primary_key"`

	FkProductIdReferer int
	FkProductsId       Products `gorm:"column:fk_product_id_referer;	foreignKey:FkProductIdReferer;	not null"`
	FkStoreIdReferer   int
	StoreId            Stores `gorm:"column:fk_store_id_referer;	foreignKey:FkStoreIdReferer;	not null"`

	DeliveryProductTotal int       `gorm:"column:delivery_product_total;	type:bigint;"`
	DeliveryCreateAt     time.Time `gorm:"column:delivery_create_at;	type:timestamp;"`
}
