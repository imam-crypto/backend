package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Email     string         `gorm:"not null;unique"`
	Password  string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
type DisplayUsers struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Paginations struct {
	Limit        int         `json:"limit"`
	Page         int         `json:"page"`
	Sort         string      `json:"sort"`
	TotalRows    int         `json:"total_rows"`
	FirstPage    string      `json:"first_page"`
	PreviousPage string      `json:"previous_page"`
	NextPage     string      `json:"next_page"`
	LastPage     string      `json:"last_page"`
	FromRow      int         `json:"from_row"`
	ToRow        int         `json:"to_row"`
	Rows         interface{} `json:"rows"`
}
