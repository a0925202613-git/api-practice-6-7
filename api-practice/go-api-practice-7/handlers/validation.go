package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// formatValidationError 把「綁定／驗證請求 body 時產生的錯誤」整理成 API 要回傳的格式。
// 呼叫方會用回傳的狀態碼和 body 來寫 c.JSON(status, body)。
// - 若 err 為 nil：回傳 0, nil，表示沒有錯誤，呼叫方不必回傳錯誤。
// - 若是 validator 的欄位驗證錯誤：回傳 400，以及一個包含「資料驗證失敗」總訊息與「每個欄位」的錯誤明細（field + message），方便前端顯示在哪一欄出錯。
// - 若是其他類型的錯誤（例如 JSON 格式錯誤）：回傳 400，以及單純的錯誤訊息字串。
func formatValidationError(err error) (int, interface{}) {
	if err == nil {
		return 0, nil
	}
	if errs, ok := err.(validator.ValidationErrors); ok {
		details := make([]gin.H, 0, len(errs))
		for _, e := range errs {
			// 欄位名稱改成小寫開頭（例如 BookId -> bookId），與 JSON 慣例一致，前端比較好對應
			field := e.Field()
			if len(field) > 0 {
				field = strings.ToLower(field[:1]) + field[1:]
			}
			details = append(details, gin.H{
				"field":   field,
				"message": validationMessage(e.Tag(), e.Param()),
			})
		}
		return http.StatusBadRequest, gin.H{
			"error":   "資料驗證失敗",
			"details": details,
		}
	}
	// 不是 validator 的 ValidationErrors（例如 binding 錯誤），直接回傳該錯誤的字串
	return http.StatusBadRequest, gin.H{"error": err.Error()}
}

// validationMessage 根據「驗證規則代碼」（例如 required、max、min）和該規則的參數（例如 max 的長度數字），
// 回傳一句給使用者看的錯誤說明，例如「此欄位為必填」「超過最大長度 100」。
func validationMessage(tag, param string) string {
	switch tag {
	case "required":
		return "此欄位為必填"
	case "max":
		return "超過最大長度 " + param
	case "min":
		return "低於最小值 " + param
	case "gte":
		return "須大於等於 " + param
	case "lte":
		return "須小於等於 " + param
	default:
		return "不符合規則: " + tag
	}
}
