package core

type User struct {
	Id       int    `json:"-" db:"id"`
	Login    string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
