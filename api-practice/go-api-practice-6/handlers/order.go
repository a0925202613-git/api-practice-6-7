package handlers

import (
	"database/sql"
	"net/http"

	"go-api-practice-6/database"
	"go-api-practice-6/models"

	"github.com/gin-gonic/gin"
)

// GetOrders 取得訂單列表。
// 從網址 query 讀取 status（pending / completed / cancelled 三種之一）。若有傳且是這三種之一，就只回傳該狀態的訂單；沒傳或不是有效值，就回傳全部訂單。
// 每筆訂單都要帶出「對應的菜單名稱」（訂單只存 menu_id，回傳時要讓前端看到菜單名稱），列表依訂單 id 排序。
// 成功時回傳 200 與訂單陣列（含菜單名稱）；沒有符合的訂單就回傳空陣列。
func GetOrders(c *gin.Context) {
	status := c.Query("status") // pending | completed | cancelled，空字串表示全部

	//基礎的查詢
	query := "SELECT o.id, o.menu_id, m.name, o.quantity, o.status, o.created_at, o.updated_at FROM orders o JOIN menus m ON o.menu_id = m.id"

	// 準備一個百寶袋 args 來裝參數 (裡面裝什麼型態都可以，所以用 any 或 interface{})
	var args []interface{}

	//檢查有沒有傳入有效的 status，如果有就加上 WHERE 條件
	if status == "pending" || status == "completed" || status == "cancelled" {
		//把 WHERE 條件拼接到 query 字串後面
		query += " WHERE o.status = $1"
		//把 status 這個變數放進百寶袋 args 裡面
		args = append(args, status)
	}
	// 最後都要加上排序
	query += " ORDER BY o.id ASC"

	rows, err := database.DB.Query(query, args...) // args... 是把 args 這個切片裡的元素一個一個拿出來當作參數傳給 Query
	if err != nil {
		respondError(c, err)
		return
	}
	defer rows.Close()

	orders := make([]models.OrderWithMenuName, 0)
	for rows.Next() {
		var o models.OrderWithMenuName
		if err := rows.Scan(&o.ID, &o.MenuID, &o.MenuName, &o.Quantity, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			respondError(c, err)
			return
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, orders)

}

// GetOrderByID 依網址上的 id 取得單一筆訂單。
// 回傳的訂單資料要包含「該訂單對應的菜單名稱」（不只 menu_id，要有名稱給前端顯示）。
// 若該 id 沒有對應的訂單（例如 id 不存在或已被刪除），回傳 404；找到就回傳 200 與該筆訂單（含菜單名稱）。
func GetOrderByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	query := `SELECT o.id, o.menu_id, m.name, o.quantity, o.status, o.created_at, o.updated_at 
	 FROM orders o JOIN menus m ON o.menu_id = m.id 
	 WHERE o.id = $1`

	//準備空箱子來裝查詢結果
	var o models.OrderWithMenuName

	//執行查詢，並把結果掃描到剛剛準備的空箱子裡面
	if err := database.DB.QueryRow(query, id).Scan(
		&o.ID,
		&o.MenuID,
		&o.MenuName,
		&o.Quantity,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows { //找不到對應的訂單
			c.JSON(http.StatusNotFound, gin.H{"error": "訂單不存在"})
		} else {
			respondError(c, err) //其他查詢錯誤
		}
		return
	}

	c.JSON(http.StatusOK, o)
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

	// 檢查 menu_id 是否存在
	if !ValidateMenuExists(c, input.MenuID) {
		return
	}

	query := "INSERT INTO orders (menu_id, quantity, status) VALUES ($1, $2, $3) RETURNING id, menu_id, quantity, status, created_at, updated_at"

	//準備空箱子來裝新增後的訂單資料
	var o models.Order

	//執行新增，並把新增後的訂單資料掃描到剛剛準備的空箱子裡面
	if err := database.DB.QueryRow(query, input.MenuID, input.Quantity, input.Status).Scan(
		&o.ID,
		&o.MenuID,
		&o.Quantity,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		respondError(c, err)
		return
	}

	//新增成功，回傳 201 與剛建立的那筆訂單資料
	c.JSON(http.StatusCreated, o)
}

func UpdateOrder(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var input models.Order
	if err := c.ShouldBindJSON(&input); err != nil {
		status, body := formatValidationError(err)
		c.JSON(status, body)
		return
	}

	if !ValidateMenuExists(c, input.MenuID) {
		return
	}

	query := "UPDATE orders SET menu_id = $1, quantity = $2, status = $3, updated_at = NOW() WHERE id = $4 RETURNING id, menu_id, quantity, status, created_at, updated_at"

	//準備空箱子來裝新增後的訂單資料
	var o models.Order

	//執行新增，並把新增後的訂單資料掃描到剛剛準備的空箱子裡面
	if err := database.DB.QueryRow(query, input.MenuID, input.Quantity, input.Status, id).Scan(
		&o.ID,
		&o.MenuID,
		&o.Quantity,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		respondError(c, err)
		return
	}

	//新增成功，回傳 201 與剛建立的那筆訂單資料
	c.JSON(http.StatusCreated, o)
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

	query := "UPDATE orders SET status = 'cancelled', updated_at = NOW() WHERE id = $1 RETURNING id, menu_id, quantity, status, created_at, updated_at"

	//準備空箱子來裝更新後的訂單資料
	var o models.Order

	//執行更新，並把更新後的訂單資料掃描到剛剛準備的空箱子裡面
	if err := database.DB.QueryRow(query, id).Scan(
		&o.ID,
		&o.MenuID,
		&o.Quantity,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows { //找不到對應的訂單
			c.JSON(http.StatusNotFound, gin.H{"error": "訂單不存在"})
		} else {
			respondError(c, err) //其他查詢錯誤
		}
		return
	}

	//取消成功，回傳 200 與更新後的訂單資料
	c.JSON(http.StatusOK, o)
}

// ValidateMenuExists 檢查 menuID 是否在菜單表內存在。
// 若不存在則已寫入 c.JSON(400, gin.H{"error": "menu_id 不存在"}) 並回傳 false；存在則回傳 true。
// 使用於 CreateOrder：通過後才新增訂單。
func ValidateMenuExists(c *gin.Context, menuID int) bool {
	var id int
	query := "SELECT id FROM menus WHERE id = $1"
	if err := database.DB.QueryRow(query, menuID).Scan(
		&id,
	); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "menu_id 不存在"})
		} else {
			respondError(c, err)
		}
		return false
	}

	return true
}

// ValidateOrderCanCancel 檢查該訂單是否存在且狀態為「待處理」，才允許取消。
// 若訂單不存在則已寫入 404 並回傳 false；若存在但狀態不是待處理則已寫入 400「僅能取消待處理中的訂單」並回傳 false；可取消則回傳 true。
// 使用於 CancelOrder：通過後才將訂單狀態改為已取消。
func ValidateOrderCanCancel(c *gin.Context, orderID int) bool {
	var s string
	query := "SELECT status FROM orders WHERE id = $1"
	if err := database.DB.QueryRow(query, orderID).Scan(
		&s,
	); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "訂單不存在"})
		} else {
			respondError(c, err)
		}
		return false
	} else {
		if s != "pending" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "僅能取消待處理中的訂單"})
			return false
		}
	}
	return true
}
