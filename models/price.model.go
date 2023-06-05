package models

import (
	"time"

	"github.com/google/uuid"
)

// изменение цены на товар
type PriceChanges struct {
	PkPriceChangeId uuid.UUID `gorm:"column:pk_price_change_id;	type:uuid;	default:uuid_generate_v4();	primary_key"`

	FkProductIdReferer int
	FkProductsId       Products `gorm:"column:fk_product_id_referer;	foreignKey:FkProductIdReferer;		not null"`

	PriceNew      int       `gorm:"column:price_new;	type:varchar(30);	not null"`
	PriceUpdateAt time.Time `gorm:"column:price_update_at;	type:timestamp;"`
}
