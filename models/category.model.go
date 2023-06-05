package models

import (
	"time"

	"github.com/google/uuid"
)

// категория товара
type Category struct {
	PkCategoryId        uuid.UUID `gorm:"column:pk_category_id;				type:uuid;				default:uuid_generate_v4();	primary_key"`
	CategoryName        string    `gorm:"column:category_name;				type:varchar(30);	unique; not null"`
	CategoryDescription string    `gorm:"column:category_description;	type:varchar(255);"`

	CategoryCreateAt time.Time `gorm:"column:category_stat_create_at;	type:timestamp;"`
	CategoryDeleteAt time.Time `gorm:"column:category_stat_delete_at;	type:timestamp;"`
}
