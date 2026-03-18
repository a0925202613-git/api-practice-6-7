package handlers

import (
	"net/http"

	"go-api-practice-7/models"

	"github.com/gin-gonic/gin"
)

// GetBorrowals 取得所有借閱紀錄的列表。
// 每筆借閱紀錄要帶出「對應的書名」（借閱只存 book_id，回傳時要讓前端看到書名），列表依借閱 id 排序。
// 成功時回傳 200 與借閱紀錄陣列（含書名）；若沒有資料就回傳空陣列。注意：歸還時間可能尚未填寫（未還書），需能正確處理空值。
func GetBorrowals(c *gin.Context) {
	// TODO: 查詢所有借閱紀錄，一併取得每筆對應的書名，結果依借閱 id 排序
	// TODO: 組出列表（歸還時間可能為空），回傳 200（無資料就空陣列）
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 GetBorrowals"})
}

// GetBorrowalByID 依網址上的 id 取得單一筆借閱紀錄。
// 回傳的資料要包含「該筆借閱對應的書名」（不只 book_id，要有書名給前端顯示）。若該 id 沒有對應的借閱紀錄，回傳 404；找到就回傳 200 與該筆資料（含書名）。歸還時間可能為空（尚未還書）。
func GetBorrowalByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	_ = id
	// TODO: 用 id 查出一筆借閱紀錄，並一併取得該筆對應的書名
	// TODO: 查不到就回 404；查到就回 200 與該筆借閱資料（含書名）
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 GetBorrowalByID"})
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
	// TODO: 在同一筆交易內新增借閱紀錄，並把該書改為已借出，取得建立後的借閱資料（含書名），回傳 201
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 CreateBorrowal（需檢查書可借閱）"})
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
	// TODO: 寫入該筆的歸還時間，並把對應的書改回可借閱，回傳 200（可帶更新後資料或成功訊息）
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 ReturnBorrowal（僅未歸還者可還書）"})
}

// ValidateBookAvailable 檢查該書是否存在且狀態為「可借閱」。
// 若書不存在或已被借出則已寫入 c.JSON(400, gin.H{"error": "此書目前已被借出或不存在"}) 並回傳 false；可借閱則回傳 true。
// 使用於 CreateBorrowal：通過後才新增借閱紀錄並將書改為已借出。
func ValidateBookAvailable(c *gin.Context, bookID int) bool {
	// TODO: 實作
	return false
}

// ValidateBorrowalCanReturn 檢查該筆借閱是否存在且「尚未歸還」（歸還時間為空），才允許還書。
// 若借閱不存在則已寫入 404 並回傳 false；若已歸還則已寫入 400「此筆借閱已歸還」並回傳 false；可還書則回傳 true。
// 使用於 ReturnBorrowal：通過後才寫入歸還時間並將書改回可借閱。
func ValidateBorrowalCanReturn(c *gin.Context, borrowalID int) bool {
	// TODO: 實作
	return false
}
