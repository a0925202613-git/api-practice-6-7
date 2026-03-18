package handlers

import (
	"net/http"

	"go-api-practice-6/models"

	"github.com/gin-gonic/gin"
)

// GetOrders 取得訂單列表。
// 從網址 query 讀取 status（pending / completed / cancelled 三種之一）。若有傳且是這三種之一，就只回傳該狀態的訂單；沒傳或不是有效值，就回傳全部訂單。
// 每筆訂單都要帶出「對應的菜單名稱」（訂單只存 menu_id，回傳時要讓前端看到菜單名稱），列表依訂單 id 排序。
// 成功時回傳 200 與訂單陣列（含菜單名稱）；沒有符合的訂單就回傳空陣列。
func GetOrders(c *gin.Context) {
	status := c.Query("status") // pending | completed | cancelled，空字串表示全部
	_ = status
	// TODO: 從 query 讀取 status，判斷是否為有效值（pending/completed/cancelled），決定查詢時要不要加上「只查該狀態」
	// TODO: 查詢訂單時要一併取得每筆訂單對應的菜單名稱，結果依訂單 id 排序
	// TODO: 把結果組成訂單列表（含菜單名稱），回傳 200
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 GetOrders（含 status 篩選）"})
}

// GetOrderByID 依網址上的 id 取得單一筆訂單。
// 回傳的訂單資料要包含「該訂單對應的菜單名稱」（不只 menu_id，要有名稱給前端顯示）。
// 若該 id 沒有對應的訂單（例如 id 不存在或已被刪除），回傳 404；找到就回傳 200 與該筆訂單（含菜單名稱）。
func GetOrderByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	_ = id
	// TODO: 用 id 查出一筆訂單，並一併取得該訂單對應的菜單名稱
	// TODO: 查不到就回 404（例如 id 不存在）；查到就回 200 與該筆訂單資料（含菜單名稱）
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 GetOrderByID"})
}

// CreateOrder 新增一筆訂單（此 API 需帶 token）。
// 請求 body 需提供 menu_id、quantity（必填），由 ShouldBindJSON 綁定；驗證失敗時用 formatValidationError 回傳 400。
// 實作時要先確認 request 裡的 menu_id 在菜單表裡真的存在；若不存在，回傳 400，錯誤訊息為「menu_id 不存在」。
// 若 menu_id 存在，則新增一筆訂單（menu_id、quantity 照傳入值，狀態設為「待處理」），並回傳 201 與「剛建立的那筆訂單」的完整資料（含 id、時間等）。
func CreateOrder(c *gin.Context) {
	var input models.Order
	if err := c.ShouldBindJSON(&input); err != nil {
		status, body := formatValidationError(err)
		c.JSON(status, body)
		return
	}
	if !ValidateMenuExists(c, input.MenuID) {
		return
	}
	// TODO: 新增一筆訂單（狀態為待處理），取得建立後的完整一筆資料，回傳 201 與該筆資料
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 CreateOrder（含 menu 存在檢查）"})
}

// CancelOrder 依網址上的 id 取消一筆訂單（此 API 需帶 token）。
// 只有狀態為「待處理」的訂單可以取消。若該筆訂單目前是「已完成」或「已取消」，則回傳 400，錯誤訊息為「僅能取消待處理中的訂單」。
// 若該 id 根本沒有對應的訂單，應回傳 404。
// 若訂單是待處理，則將該筆訂單的狀態改為「已取消」，並回傳 200；body 可回傳更新後的訂單或簡單的成功訊息（例如 "cancelled"）。
func CancelOrder(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if !ValidateOrderCanCancel(c, id) {
		return
	}
	// TODO: 將該筆訂單狀態改為「已取消」，回傳 200（可帶更新後資料或成功訊息）
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 CancelOrder（僅 pending 可取消）"})
}

// ValidateMenuExists 檢查 menuID 是否在菜單表內存在。
// 若不存在則已寫入 c.JSON(400, gin.H{"error": "menu_id 不存在"}) 並回傳 false；存在則回傳 true。
// 使用於 CreateOrder：通過後才新增訂單。
func ValidateMenuExists(c *gin.Context, menuID int) bool {
	// TODO: 實作
	return false
}

// ValidateOrderCanCancel 檢查該訂單是否存在且狀態為「待處理」，才允許取消。
// 若訂單不存在則已寫入 404 並回傳 false；若存在但狀態不是待處理則已寫入 400「僅能取消待處理中的訂單」並回傳 false；可取消則回傳 true。
// 使用於 CancelOrder：通過後才將訂單狀態改為已取消。
func ValidateOrderCanCancel(c *gin.Context, orderID int) bool {
	// TODO: 實作
	return false
}
