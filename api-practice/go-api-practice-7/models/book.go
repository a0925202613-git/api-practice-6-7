package models

import "time"

// Book 書籍
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" binding:"required,max=500"`
	ISBN      string    `json:"isbn" binding:"max=20"`
	Available bool      `json:"available"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
