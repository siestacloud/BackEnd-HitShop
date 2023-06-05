package models

import (
	"time"

	"github.com/google/uuid"
)

// Продукты
type Products struct {
	PkProductId uuid.UUID `gorm:"column:pk_account_id;	type:uuid;	default:uuid_generate_v4();	primary_key"`

	FkCategoryIdReferer int
	FkCategoryId        Category `gorm:"column:fk_category_id_referer;	foreignKey:FkCategoryIdReferer;	not null"`

	FkManufactureIdReferer int
	FkManufactureId        Manufactures `gorm:"column:fk_manufacture_id_referer;	foreignKey:FkManufactureIdReferer;		not null"`

	ProductName        string `gorm:"column:product_name;					type:varchar(30);"`
	ProductDescription string `gorm:"column:product_description;	type:varchar(30);"`

	ProductCreateAt time.Time `gorm:"column:product_create_at;	type:timestamp;"`
	ProductDeleteAt time.Time `gorm:"column:product_delete_at;	type:timestamp;"`
}
