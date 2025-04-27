package models

import "time"

type Main struct {
	ID        int32
	Title     string
	SubID     int32
	SubObj    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Tool struct {
	ID          int32
	Title       string
	Description *string
	MainID      int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type Table struct {
	ID        int32
	Name      string
	MainID    int32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Chair struct {
	ID        int32
	Name      string
	Type      string
	MainID    int32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
