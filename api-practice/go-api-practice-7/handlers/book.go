package handlers

import (
	"net/http"

	"go-api-practice-7/models"

	"github.com/gin-gonic/gin"
)

// GetBooks 取得書籍列表。
// 從網址 query 讀取 available（"true" 或 "false"）。若為 "true" 只回傳「可借閱」的書；若為 "false" 只回傳「已借出」的書；沒傳或不是這兩個值則回傳全部書籍。
// 列表依書籍 id 排序。成功時回傳 200 與書籍陣列；沒有符合的書就回傳空陣列。
func GetBooks(c *gin.Context) {
	available := c.Query("available") // "true" | "false" | ""
	_ = available
	// TODO: 從 query 讀取 available，判斷是否為 "true" 或 "false"，決定查詢時要不要只篩選「可借閱」或「已借出」
	// TODO: 查詢書籍，結果依 id 排序，組出列表
	// TODO: 回傳 200 與列表（無資料就空陣列）
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 GetBooks（含 available 篩選）"})
}

// GetBookByID 依網址上的 id 取得單一筆書籍。
// 若該 id 沒有對應的書籍（例如 id 不存在或已被刪除），回傳 404；找到就回傳 200 與該筆書籍的完整資料（id、title、isbn、available、時間等）。
func GetBookByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	_ = id
	// TODO: 用 id 查出一筆書籍
	// TODO: 查不到就回 404；查到就回 200 與該筆書籍的完整資料
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 GetBookByID"})
}

// CreateBook 新增一筆書籍（此 API 需帶 token）。
// 請求 body 需提供 title（必填）、isbn 等，由 ShouldBindJSON 綁定；驗證失敗時用 formatValidationError 回傳 400。
// 驗證通過後，新增一筆書籍（available 預設為「可借閱」），並取得建立後的完整一筆資料（含 id、created_at、updated_at 等），回傳 201 與該筆資料。
func CreateBook(c *gin.Context) {
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		status, body := formatValidationError(err)
		c.JSON(status, body)
		return
	}
	_ = input
	// TODO: 依 request 新增一筆書籍（可借閱狀態預設為 true），取得建立後的那一筆完整資料
	// TODO: 回傳 201 與該筆資料
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 CreateBook"})
}

// UpdateBook 依網址上的 id 更新一筆書籍（此 API 需帶 token）。
// 請求 body 提供要更新的 title、isbn 等，由 ShouldBindJSON 綁定；驗證失敗就回傳 400。
// 若該 id 在資料庫裡不存在，回傳 404。若存在，則用 request 的內容更新該筆書籍的 title、isbn（以及更新時間）；「是否可借閱」通常由借書／還書流程更新，此 API 可不改。
// 更新成功後回傳 200 與更新後的該筆完整資料。
func UpdateBook(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		status, body := formatValidationError(err)
		c.JSON(status, body)
		return
	}
	_, _ = id, input
	// TODO: 用 id 找到該筆書籍，用 request 的 title、isbn 更新（可借閱狀態可不在此改）
	// TODO: 若沒有該 id（更新影響筆數為 0），回傳 404；有更新到就回 200 與更新後的該筆完整資料
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 UpdateBook"})
}

// DeleteBook 依網址上的 id 刪除一筆書籍（此 API 需帶 token）。
// 若該 id 在資料庫裡不存在，回傳 404；若存在並成功刪除，回傳 200，body 可帶簡單成功訊息（例如 "message": "deleted"）。
func DeleteBook(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	_ = id
	// TODO: 用 id 刪除該筆書籍
	// TODO: 若沒有該 id（刪除影響筆數為 0），回傳 404；有刪到就回 200 與成功訊息
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 DeleteBook"})
}
