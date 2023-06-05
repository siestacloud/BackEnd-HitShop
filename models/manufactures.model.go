package models

import (
	"time"

	"github.com/google/uuid"
)

// производитель товара
type Manufactures struct {
	PkManufactureId        uuid.UUID `gorm:"column:pk_manufacture_id;								type:uuid;				default:uuid_generate_v4();	primary_key"`
	ManufactureName        string    `gorm:"column:manufacture_status_name;					type:varchar(30);	not null"`
	ManufactureDescription string    `gorm:"column:manufacture_status_description;	type:varchar(255);"`

	CreateAt time.Time `gorm:"column:manufacture_create_at;	type:timestamp;"`
	DeleteAt time.Time `gorm:"column:manufacture_delete_at;	type:timestamp;"`
}
