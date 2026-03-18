package models

import "time"

// Borrowal 借閱紀錄（book_id + 借閱者，returned_at 為空表示尚未歸還）
type Borrowal struct {
	ID         int        `json:"id"`
	BookID     int        `json:"book_id" binding:"required,gte=1"`
	UserName   string     `json:"user_name" binding:"required,max=255"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	ReturnedAt *time.Time `json:"returned_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// BorrowalWithBookTitle 借閱紀錄＋書名（JOIN 查詢用）
type BorrowalWithBookTitle struct {
	ID         int        `json:"id"`
	BookID     int        `json:"book_id"`
	BookTitle  string     `json:"book_title"`
	UserName   string     `json:"user_name"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	ReturnedAt *time.Time `json:"returned_at,omitempty"`
}
