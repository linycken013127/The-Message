package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/response"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 認證中介軟體
// 從 Header 讀取 X-Account-ID 並驗證帳號存在
func AuthMiddleware(accountRepo repository.AccountRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 檢查 header 是否存在
		accountIDStr := c.GetHeader("X-Account-ID")
		if accountIDStr == "" {
			response.Error(c, http.StatusUnauthorized, "未提供認證資訊")
			c.Abort()
			return
		}

		// 2. 檢查格式是否正確
		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "無效的認證資訊")
			c.Abort()
			return
		}

		// 3. 檢查帳號 ID 是否合法
		if accountID <= 0 {
			response.Error(c, http.StatusUnauthorized, "帳號不存在")
			c.Abort()
			return
		}

		// 4. 檢查帳號是否存在於資料庫
		account, err := accountRepo.GetAccountById(context.Background(), accountID)
		if err != nil || account == nil {
			response.Error(c, http.StatusUnauthorized, "帳號不存在")
			c.Abort()
			return
		}

		// 5. 將帳號資訊存入 context
		c.Set("accountID", accountID)
		c.Set("account", account)
		c.Next()
	}
}

// GetAccountID 從 context 取得帳號 ID
func GetAccountID(c *gin.Context) int {
	accountID, exists := c.Get("accountID")
	if !exists {
		return 0
	}
	return accountID.(int)
}
