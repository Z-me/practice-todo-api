package model

import "time"

type Todo struct {
	ID        uint		`gorm:"primaryKey"`
	Title     string
	Status		string
	Details		string
	Priority	string
	CreatedAt	time.Time
	UpdatedAt	time.Time
}

type TodoList []Todo

type Payload struct {
  Title     string
  Status    string
  Details   string
  Priority  string
}

type Status struct {
  Status    string
}
