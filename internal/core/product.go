package core

import "time"

type Product struct {
	UUID        string
	Category    string
	Manufacture string
	Description string
	CreateAt    time.Time
	DeleteAt    time.Time
}
