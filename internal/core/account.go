package core

import (
	"time"

	"github.com/google/uuid"
)

// User имплементирует клиента
type Account struct {
	UUID uuid.UUID `json:"-" db:"pk_account_id"`

	Role     string `json:"role" validate:"required" `
	Email    string `json:"email"  validate:"required,email" db:"account_email"`
	Status   string
	Password string `json:"-" validate:"required,min=8,containsany=!@#?*" db:"account_password_hash"`

	AccountPhoneNumber     string `json:"-"`
	AccountPhoneNumberHash string `json:"-"`

	Verified         bool   `json:"verify" db:"account_verify"`
	VerificationCode string `json:"-" db:"account_verify_code"`

	AccountStatistic
	Favorites

	CreateAt time.Time `db:"account_create_at"`
	UpdateAt time.Time `db:"account_update_at"`
	DeleteAt time.Time `json:"-" db:"account_delete_at"`
}

type AccountStatistic struct {
	PurchaseCount int    `json:"-"`
	Description   string `json:"-"`
}

type Favorites struct {
	Products []Product
}

type SignUpInput struct {
	Role            string `json:"role" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,containsany=!@#?*"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required" validate:"required,min=8,containsany=!@#?*"`
}

type SignInInput struct {
	Role     string `json:"role" validate:"required"`
	Email    string `json:"email"  binding:"required" validate:"required,email"`
	Password string `json:"password"  binding:"required" validate:"required,min=8,containsany=!@#?*"`
}

type ChangePassInput struct {
	Password    string `json:"password"  binding:"required" validate:"required,min=8,containsany=!@#?*"`
	PasswordNew string `json:"password_new"  binding:"required,min=8" validate:"required,min=8,containsany=!@#?*"`
}
