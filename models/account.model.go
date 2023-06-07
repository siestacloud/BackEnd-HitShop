package models

import (
	"time"

	"github.com/google/uuid"
)

// роли аккаунта
type AccountRoles struct {
	PkAccountRoleId        uuid.UUID `gorm:"column:pk_account_role_id; 			 type:uuid;					default:uuid_generate_v4();	primary_key"`
	AccountRoleName        string    `gorm:"column:account_role_name;  			 type:varchar(30);	unique; not null"`
	AccountRoleDescription string    `gorm:"column:account_role_description; type:varchar(255);"`
}

// статусы аккаунта
type AccountStatuses struct {
	PkAccountStatusId        uuid.UUID `gorm:"column:pk_account_status_id;				type:uuid;				default:uuid_generate_v4();	primary_key"`
	AccountStatusName        string    `gorm:"column:account_status_name;					type:varchar(30);	not null"`
	AccountStatusDescription string    `gorm:"column:account_status_description;	type:varchar(255);"`
}

// cтатистика аккаунта
type AccountStat struct {
	PkAccountStatId uuid.UUID `gorm:"column:pk_account_stat_id;	type:uuid;	default:uuid_generate_v4();	primary_key"`

	FkAccountStatusIdReferer int
	FkAccountStatusId        AccountStatuses `gorm:"column:fk_account_status_id_referer;	foreignKey:FkAccountStatusIdReferer;	constraint:OnDelete:CASCADE;	not null"`

	AccountStatPurchaseCount int64  `gorm:"column:account_stat_purchase_count;	type:bigint;"`
	AccountStatDescription   string `gorm:"column:account_stat_description;		type:varchar(255);"`
}

// Избранные товары клиентов, с ссылками
type AccountFavorites struct {
	PkAccountFavoriteId uuid.UUID `gorm:"column:pk_account_favorite_id;	type:uuid;	default:uuid_generate_v4();	primary_key"`

	FkAccountIdReferer int
	FkAccountId        Accounts `gorm:"type:uuid; column:fk_account_id_referer;	foreignKey:FkAccountIdReferer;	not null"`

	FkProductIdReferer int
	FkProductId        Products `gorm:"column:fk_product_id_referer;	foreignKey:FkProductIdReferer;	not null"`
}

// аккаунты
type Accounts struct {
	PkAccountId uuid.UUID `gorm:"column:pk_account_id;	type:uuid;	default:uuid_generate_v4();	primary_key"`

	FkRoleIdReferer int
	RoleId          AccountRoles `gorm:"column:fk_role_id_referer; foreignKey:FkRoleIdReferer;"`

	FkAccountStatIdReferer int
	AccountStatId          AccountStatuses `gorm:"column:fk_account_stat_id_referer;	foreignKey:FkAccountStatIdReferer;	constraint:OnDelete:CASCADE;	not null"`

	AccountEmail      string `gorm:"column:account_email;					type:varchar(30);	unique; 	not null"`
	AccountVerify     bool   `gorm:"column:account_verify;					type:boolean;							not null"`
	AccountVerifyCode string `gorm:"column:account_verify_code;					type:varchar(90);			not null"`

	AccountPasswordHash    string `gorm:"column:account_password_hash;	type:varchar(90);					not null"`
	AccountPhoneNumberHash string `gorm:"column:account_phone_number;		type:varchar(30);	unique;"`

	CreateAt time.Time `gorm:"column:account_create_at;			type:timestamp;"`
	UpdateAt time.Time `gorm:"column:account_update_at;			type:timestamp;"`
	DeleteAt time.Time `gorm:"column:account_delete_at;			type:timestamp;"`
}
