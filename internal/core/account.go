package core

import (
	"time"

	"github.com/google/uuid"
)

// User имплементирует клиента
type Account struct {
	UUID uuid.UUID `json:"-" db:"pk_account_id"`

	Role                   string `json:"role" validate:"required"`
	Email                  string `json:"email"  validate:"required"`
	Status                 string
	Verify                 bool   `json:"verify"`
	Password               string `json:"password" validate:"required,min=5"`
	AccountPhoneNumber     string
	AccountPhoneNumberHash string

	AccountStatistic
	Favorites
}

type AccountStatistic struct {
	PurchaseCount int
	Description   string
	CreateAt      time.Time
	DeleteAt      time.Time
}

type Favorites struct {
	Products []Product
}
