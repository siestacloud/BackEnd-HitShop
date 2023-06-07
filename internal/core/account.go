package core

import (
	"time"

	"github.com/google/uuid"
)

// User имплементирует клиента
type Account struct {
	UUID uuid.UUID `json:"-" db:"pk_account_id"`

	Role     string `json:"role" validate:"required"`
	Email    string `json:"email"  validate:"required"`
	Status   string
	Password string `json:"password" validate:"required,min=8"`

	AccountPhoneNumber     string
	AccountPhoneNumberHash string

	Verified         bool `json:"verify"`
	VerificationCode string

	AccountStatistic
	Favorites

	CreateAt time.Time
	UpdateAt time.Time
	DeleteAt time.Time
}

type AccountStatistic struct {
	PurchaseCount int
	Description   string
}

type Favorites struct {
	Products []Product
}

type SignUpInput struct {
	Role            string `json:"role" validate:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required" validate:"required"`
}

type SignInInput struct {
	Role     string `json:"role" validate:"required"`
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}
