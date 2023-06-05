package models

import "github.com/google/uuid"

// филиалы
type Stores struct {
	PkStoreId    uuid.UUID `gorm:"column:pk_store_id;			type:uuid;	default:uuid_generate_v4();	primary_key"`
	StoreName    string    `gorm:"column:store_name;			type:varchar(30);"`
	StoreAddress string    `gorm:"column:store_address;		type:varchar(255);"`
}
