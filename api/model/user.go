package model

import "time"

type User struct {
	ID			uint		`gorm:"primaryKey"`
	Name		string
	Password	string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ID uint
