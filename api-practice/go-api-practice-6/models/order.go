package models

import "time"

// OrderStatus 訂單狀態：pending 待處理, completed 已完成, cancelled 已取消
const (
	OrderStatusPending   = "pending"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

// Order 訂單（關聯菜單）
type Order struct {
	ID        int       `json:"id"`
	MenuID    int       `json:"menu_id" binding:"required,gte=1"`
	Quantity  int       `json:"quantity" binding:"required,gte=1"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderWithMenuName 訂單＋菜單名稱（JOIN 查詢用）
type OrderWithMenuName struct {
	ID         int       `json:"id"`
	MenuID     int       `json:"menu_id"`
	MenuName   string    `json:"menu_name"`
	Quantity   int       `json:"quantity"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
