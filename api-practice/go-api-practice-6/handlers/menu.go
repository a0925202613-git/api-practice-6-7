package handlers

import (
	"net/http"

	"go-api-practice-6/database"
	"go-api-practice-6/models"

	"github.com/gin-gonic/gin"
)

// GetMenus 取得「所有」菜單品項的列表。
// 不帶任何篩選條件，就是把目前資料庫裡所有菜單一筆一筆列出來，依菜單 id 由小到大排序後回傳。
// 若目前沒有任何菜單，回傳空陣列 []；有資料就回傳 200 與該陣列。
func GetMenus(c *gin.Context) {

	// TODO: 查詢所有菜單，結果依 id 排序
	query := "SELECT id, name, category, price, created_at, updated_at FROM menus ORDER BY id ASC"
	// 執行 SQL 查詢，取得所有菜單資料的結果集（rows）和可能的錯誤（err）
	rows, err := database.DB.Query(query) 
	if err != nil {
		respondError(c, err)
		return
	}
	defer rows.Close() // 確保在函式結束前關閉資料庫連線，避免資源洩漏
	
	 menus := make([]models.Menu, 0) // 建立一個空的菜單切片，用來存放查詢結果	

	for rows.Next() { // 迭代查詢結果的每一行
		var m models.Menu // 存放當前行的資料
		if err := rows.Scan(&m.ID, &m.Name, &m.Category, &m.Price, &m.CreatedAt, &m.UpdatedAt); err != nil {
			respondError(c, err)
			return
		}
		menus = append(menus, m) // 將當前行的菜單資料加入到菜單切片中
	}
	if err := rows.Err(); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, menus)
}

// GetMenuByID 依網址上的 id 取得「單一筆」菜單。
// 若該 id 沒有對應的菜單（例如 id 不存在或已被刪除），回傳 404；找到就回傳 200 與該筆菜單的完整資料（id、name、category、price、時間等）。
func GetMenuByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	_ = id
	// TODO: 用 id 查出一筆菜單
	// TODO: 查不到就回 404；查到就回 200 與該筆菜單的完整資料
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 GetMenuByID"})
}

// CreateMenu 新增一筆菜單（此 API 需帶 token）。
// 請求 body 需提供 name、category、price 等欄位，由 ShouldBindJSON 綁定到 models.Menu；若有必填未填或格式不符，驗證失敗，用 formatValidationError 回傳 400。
// 驗證通過後，在資料庫新增一筆菜單（欄位值依 request），並取得「建立後」的完整一筆資料（含自動產生的 id、created_at、updated_at 等），回傳 201 與該筆資料。
func CreateMenu(c *gin.Context) {
	var input models.Menu
	if err := c.ShouldBindJSON(&input); err != nil {
		status, body := formatValidationError(err)
		c.JSON(status, body)
		return
	}
	_ = input
	// TODO: 依 request 內容新增一筆菜單，並取得建立後的那一筆完整資料（含 id、時間等）
	// TODO: 回傳 201 與該筆資料
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 CreateMenu"})
}

// UpdateMenu 依網址上的 id 更新一筆菜單（此 API 需帶 token）。
// 請求 body 提供要更新的 name、category、price，由 ShouldBindJSON 綁定；驗證失敗就回傳 400。
// 若該 id 在資料庫裡不存在（沒有對應的菜單），回傳 404；若存在，則用 request 的內容更新該筆菜單的 name、category、price（以及更新時間），並回傳 200 與「更新後」的該筆完整資料。
func UpdateMenu(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var input models.Menu
	if err := c.ShouldBindJSON(&input); err != nil {
		status, body := formatValidationError(err)
		c.JSON(status, body)
		return
	}
	_, _ = id, input
	// TODO: 用 id 找到該筆菜單，並用 request 的 name、category、price 更新
	// TODO: 若沒有該 id（更新影響筆數為 0），回傳 404；有更新到就回 200 與更新後的該筆完整資料
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 UpdateMenu"})
}

// DeleteMenu 依網址上的 id 刪除一筆菜單（此 API 需帶 token）。
// 若該 id 在資料庫裡不存在（沒有對應的菜單可刪），回傳 404；若存在並成功刪除，回傳 200，body 可帶簡單成功訊息（例如 "message": "deleted"）。
func DeleteMenu(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	_ = id
	// TODO: 用 id 刪除該筆菜單
	// TODO: 若沒有該 id（刪除影響筆數為 0），回傳 404；有刪到就回 200 與成功訊息
	c.JSON(http.StatusNotImplemented, gin.H{"error": "請實作 DeleteMenu"})
}
