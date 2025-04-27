package models

import (
	"time"
)

type Main struct {
	ID        int32      `json:"id"`
	Title     string     `json:"title"`
	SubID     *int32     `json:"sub_id"`
	SubObj    *string    `json:"sub_obj"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Tools  []Tool  `json:"tools,omitempty"`
	Tables []Table `json:"tables,omitempty"`
	Chairs []Chair `json:"chairs,omitempty"`
}
type Tool struct {
	ID          int32      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	MainID      int32      `json:"main_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type Table struct {
	ID        int32      `json:"id"`
	Name      string     `json:"name"`
	MainID    int32      `json:"main_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Chair struct {
	ID        int32      `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	MainID    int32      `json:"main_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
type MainInput struct {
	Title  string  `json:"title"`
	SubID  *int32  `json:"sub_id"`
	SubObj *string `json:"sub_obj"`
}

type MainUpdateInput struct {
	Title     *string `json:"title"`
	SubID     *int32  `json:"sub_id"`
	SubObj    *string `json:"sub_obj"`
	DeletedAt *string `json:"deleted_at"`
}

type ToolInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type TableInput struct {
	Name string `json:"name"`
}

type ChairInput struct {
	Name string `json:"name"`
	Type string `json:"type"` // "abc" или "cde"
}
