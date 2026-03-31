package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"go-api-practice-7/database"
	"go-api-practice-7/models"

	"github.com/gin-gonic/gin"
)

// GetBorrowals 取得所有借閱紀錄的列表。
// 每筆借閱紀錄要帶出「對應的書名」（借閱只存 book_id，回傳時要讓前端看到書名），列表依借閱 id 排序。
// 成功時回傳 200 與借閱紀錄陣列（含書名）；若沒有資料就回傳空陣列。注意：歸還時間可能尚未填寫（未還書），需能正確處理空值。
func GetBorrowals(c *gin.Context) {
	// TODO: 查詢所有借閱紀錄，一併取得每筆對應的書名，結果依借閱 id 排序
	query := `
	SELECT b.id, b.book_id, bk.title, b.user_name, b.borrowed_at, b.returned_at
	FROM borrowals b
	JOIN books bk ON b.book_id = bk.id
	ORDER BY b.id
`

	borrowals := []models.BorrowalWithBookTitle{}

	rows, err := database.DB.Query(query)
	if err != nil {
		respondError(c, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var b models.BorrowalWithBookTitle
		if err := rows.Scan(&b.ID, &b.BookID, &b.BookTitle, &b.UserName, &b.BorrowedAt, &b.ReturnedAt); err != nil {
			respondError(c, err)
			return
		}
		borrowals = append(borrowals, b)
	}

	if err := rows.Err(); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, borrowals)
}

// GetBorrowalByID 依網址上的 id 取得單一筆借閱紀錄。
// 回傳的資料要包含「該筆借閱對應的書名」（不只 book_id，要有書名給前端顯示）。若該 id 沒有對應的借閱紀錄，回傳 404；找到就回傳 200 與該筆資料（含書名）。歸還時間可能為空（尚未還書）。
func GetBorrowalByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	// TODO: 用 id 查出一筆借閱紀錄，並一併取得該筆對應的書名
	query := `
	SELECT b.id, b.book_id, bk.title, b.user_name, b.borrowed_at, b.returned_at
	FROM borrowals b
	JOIN books bk ON b.book_id = bk.id
	WHERE b.id = $1
`
	var b models.BorrowalWithBookTitle

	err := database.DB.QueryRow(query, id).Scan(&b.ID, &b.BookID, &b.BookTitle, &b.UserName, &b.BorrowedAt, &b.ReturnedAt)
	if err != nil {
		if err == sql.ErrNoRows { //找不到對應的借閱紀錄
			c.JSON(http.StatusNotFound, gin.H{"error": "借閱紀錄不存在"})
		} else {
			respondError(c, err) //其他查詢錯誤
		}
		return
	}

	// TODO: 查不到就回 404；查到就回 200 與該筆借閱資料（含書名）
	c.JSON(http.StatusOK, b)
}

// CreateBorrowal 建立一筆借閱紀錄（借書）（此 API 需帶 token）。
// 請求 body 需提供 book_id、user_name（必填），由 ShouldBindJSON 綁定；驗證失敗時用 formatValidationError 回傳 400。
// 實作時要先確認該 book_id 的書存在，且目前狀態為「可借閱」；若書不存在或已被借出，回傳 400，錯誤訊息為「此書目前已被借出」或「此書目前已被借出或不存在」。
// 若可借閱，則：新增一筆借閱紀錄（book_id、user_name），並把該書的「可借閱」狀態改為「已借出」，再回傳 201 與「剛建立的那筆借閱紀錄」的資料（建議含書名，方便前端顯示）。新增與更新書的狀態應在同一筆交易內完成，避免同時被多人借出。
func CreateBorrowal(c *gin.Context) {
	var input models.Borrowal
	if err := c.ShouldBindJSON(&input); err != nil {
		status, body := formatValidationError(err)
		c.JSON(status, body)
		return
	}
	if !ValidateBookAvailable(c, input.BookID) {
		return
	}

	//1.開啟資料庫交易
	tx, err := database.DB.Begin()
	if err != nil {
		respondError(c, err)
		return
	}
	//設下安全網：如果中途發生錯誤，回滾交易
	defer tx.Rollback()

	var newBorrowal models.BorrowalWithBookTitle

	//2.新增借閱紀錄
	insertQuery := "INSERT INTO borrowals (book_id, user_name) VALUES ($1, $2) RETURNING id, book_id, user_name, borrowed_at, returned_at"
	err = tx.QueryRow(insertQuery, input.BookID, input.UserName).Scan(&newBorrowal.ID, &newBorrowal.BookID, &newBorrowal.UserName, &newBorrowal.BorrowedAt, &newBorrowal.ReturnedAt)
	if err != nil {
		respondError(c, err)
		return
	}

	//3.把該書改為已借出
	updateBookQuery := "UPDATE books SET available = false WHERE id = $1 RETURNING title"
	err = tx.QueryRow(updateBookQuery, input.BookID).Scan(&newBorrowal.BookTitle)
	if err != nil {
		respondError(c, err)
		return
	}

	//4.提交交易
	if err := tx.Commit(); err != nil {
		respondError(c, err)
		return
	}

	//5.回傳 201 與剛建立的借閱紀錄（含書名）
	c.JSON(http.StatusCreated, newBorrowal)
}

// ReturnBorrowal 依網址上的 id 將一筆借閱紀錄標記為「已歸還」（還書）（此 API 需帶 token）。
// 只有「尚未歸還」的借閱（歸還時間為空）可以還書。若該筆借閱已經有歸還時間，表示已還過，回傳 400，錯誤訊息為「此筆借閱已歸還」。
// 若該 id 沒有對應的借閱紀錄，回傳 404。
// 若該筆借閱尚未歸還，則：寫入該筆借閱的歸還時間，並把對應的那本書的狀態改回「可借閱」，回傳 200；body 可回傳更新後的借閱紀錄或簡單成功訊息（例如 "returned"）。
func ReturnBorrowal(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if !ValidateBorrowalCanReturn(c, id) {
		return
	}

	//1.開啟資料庫交易
	tx, err := database.DB.Begin()
	if err != nil {
		respondError(c, err)
		return
	}

	//設下安全網：如果中途發生錯誤，回滾交易
	defer tx.Rollback()


	//2.寫入該筆借閱的歸還時間
	updateQuery := `
		UPDATE borrowals 
		SET returned_at = NOW()
		WHERE id = $1 
		RETURNING book_id
	`
	var bookID int
	err = tx.QueryRow(updateQuery, id).Scan(&bookID)
	if err != nil {
		respondError(c, err)
		return
	}
	
	//3.把對應的書改回可借閱
	updateBookQuery := "UPDATE books SET available = true WHERE id = $1"
	_, err = tx.Exec(updateBookQuery, bookID)
	if err != nil {
		respondError(c, err)
		return
	}

	//4.提交交易
	if err := tx.Commit(); err != nil {
		respondError(c, err)
		return
	}

	//5.回傳 200 與簡單成功訊息
	c.JSON(http.StatusOK, gin.H{"message": "returned"})
}

// ValidateBookAvailable 檢查該書是否存在且狀態為「可借閱」。
// 若書不存在或已被借出則已寫入 c.JSON(400, gin.H{"error": "此書目前已被借出或不存在"}) 並回傳 false；可借閱則回傳 true。
// 使用於 CreateBorrowal：通過後才新增借閱紀錄並將書改為已借出。
func ValidateBookAvailable(c *gin.Context, bookID int) bool {
	var available bool
	
	//去資料庫查這本書的available狀態
	query := "SELECT available FROM books WHERE id = $1"
	err := database.DB.QueryRow(query, bookID).Scan(&available)
	if err != nil {
		if err == sql.ErrNoRows { //書不存在
			c.JSON(http.StatusBadRequest, gin.H{"error": "此書目前已被借出或不存在"})
			return false
		}
		respondError(c, err)
		return false
	}

	if !available { //書存在但不可借閱
		c.JSON(http.StatusBadRequest, gin.H{"error": "此書目前已被借出或不存在"})
		return false
	}
	return true
}

// ValidateBorrowalCanReturn 檢查該筆借閱是否存在且「尚未歸還」（歸還時間為空），才允許還書。
// 若借閱不存在則已寫入 404 並回傳 false；若已歸還則已寫入 400「此筆借閱已歸還」並回傳 false；可還書則回傳 true。
// 使用於 ReturnBorrowal：通過後才寫入歸還時間並將書改回可借閱。
func ValidateBorrowalCanReturn(c *gin.Context, borrowalID int) bool {
	var returnedAt *time.Time //關鍵：使用指標來接資料，這樣遇到資料庫的 NULL 才會變成 nil，程式才不會崩潰
	
	query := "SELECT returned_at FROM borrowals WHERE id = $1"
	err := database.DB.QueryRow(query, borrowalID).Scan(&returnedAt)
	if err != nil {
		if err == sql.ErrNoRows { //借閱紀錄不存在
			c.JSON(http.StatusNotFound, gin.H{"error": "借閱紀錄不存在"})
			return false
		}
		respondError(c, err)
		return false
	}

	if returnedAt != nil { //已經有歸還時間，表示已還過了
		c.JSON(http.StatusBadRequest, gin.H{"error": "此筆借閱已歸還"})
		return false
	}
	return true
}
