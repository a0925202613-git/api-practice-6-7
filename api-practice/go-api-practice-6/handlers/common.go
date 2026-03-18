package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 自定義錯誤變數
var ErrNotFound = errors.New("not found")
// errors.New：建立一個新的錯誤物件，內容是 "not_found"

// 如果錯誤是 ErrNotFound，回傳 404 Not Found；否則回傳 500 Internal Server Error
func respondError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// 將字串參數轉換為整數，如果轉換失敗，回傳 400 Bad Request
func parseID(c *gin.Context, param string) (int, bool) {
	// 回傳兩個值：轉換後的整數和一個布林值，表示轉換是否成功
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true

	//gin.H 是 Gin 框架提供的一個捷徑（Shortcut），它本質上就是 map[string]interface{}
}
