package models

import "time"

// Menu 餐廳菜單品項
type Menu struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required,max=200"`
	Category  string    `json:"category" binding:"max=50"`
	Price     int       `json:"price" binding:"required,gte=0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
