package core

import "time"

// TelegramAccount аккаунт через который производится работа
type TelegramAccount struct {
	AccountID  int64
	Phone      string
	Owner      string
	Status     string
	UserName   string
	LastName   string
	FirstName  string
	CreateTime string
	DeleteTime string
	TelegramApp
	TrustSessions   []TrustSession
	UntrustSessions []UntrustSession
	TelegramGroups  []TelegramGroup
	AdditionalAttributesAccount
}

type AdditionalAttributesAccount struct {
	Bot               bool
	Fake              bool
	Scam              bool
	Support           bool
	Premium           bool
	Verified          bool
	Restricted        bool
	RestrictionReason []string
}

func NewTelegramAccount(owner string) *TelegramAccount {
	return &TelegramAccount{
		Status:     "NEW",
		Owner:      owner,
		CreateTime: time.Now().Format(time.RFC1123),
	}
}

func (a *TelegramAccount) SetAttr(ID int64, Phone, Username, FirstName, LastName string, sessData []byte) {
	a.Phone = Phone
	a.AccountID = ID
	a.UserName = Username
	a.LastName = LastName
	a.FirstName = FirstName
	a.UntrustSessions = append(a.UntrustSessions, UntrustSession{Data: sessData})
}

func (a *TelegramAccount) SetAddAttr(Bot, Fake, Scam, Premium, Support, Verified bool) {
	a.Bot = Bot
	a.Fake = Fake
	a.Scam = Scam
	a.Premium = Premium
	a.Support = Support
	a.Verified = Verified
}

func (a *TelegramAccount) SetRestrict(Restricted bool, RestrictionReason []string) {

}
