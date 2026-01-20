# Task 1.3: 認證機制

## 功能概述

完善 API 認證機制，確保需要認證的 API 會驗證：
1. X-Account-ID header 是否存在
2. X-Account-ID 格式是否正確（數字）
3. X-Account-ID 對應的帳號是否存在於資料庫

### 認證方式
使用 HTTP Header `X-Account-ID` 進行簡易認證。

### 需要認證的 API
| Method | Path | 說明 |
|--------|------|------|
| POST | `/api/v1/games` | 建立遊戲房 |
| POST | `/api/v1/games/:gameId/join` | 加入遊戲房 |
| POST | `/api/v1/games/:gameId/start` | 開始遊戲 |

## 架構設計

### Clean Architecture 分層

```
┌─────────────────────────────────────────────────────────────┐
│                    Adapter Layer                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ middleware/auth_middleware.go                        │   │
│  │ - AuthMiddleware(accountRepo) gin.HandlerFunc        │   │
│  │ - GetAccountID(c) int                               │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ AccountRepository.GetAccountById()                   │   │
│  │ - 用於驗證帳號是否存在                               │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 實作細節

### AuthMiddleware (`internal/adapter/http/middleware/auth_middleware.go`)

```go
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
```

### Handler 整合

GameRoomHandler 使用 AuthMiddleware：

```go
// GameRoomHandlerOptions 遊戲房處理器選項
type GameRoomHandlerOptions struct {
    Engine          *gin.Engine
    GameRoomUseCase usecase.GameRoomUseCase
    AccountRepo     repository.AccountRepository  // 新增：用於認證
}

// RegisterGameRoomHandler 註冊遊戲房處理器
func RegisterGameRoomHandler(opts *GameRoomHandlerOptions) {
    handler := &GameRoomHandler{
        gameRoomUseCase: opts.GameRoomUseCase,
    }

    // 需要認證的路由
    authGroup := opts.Engine.Group("/api/v1")
    authGroup.Use(middleware.AuthMiddleware(opts.AccountRepo))
    {
        authGroup.POST("/games", handler.CreateGameRoom)
        authGroup.POST("/games/:gameId/join", handler.JoinGameRoom)
        authGroup.POST("/games/:gameId/start", handler.StartGame)
    }
}
```

## E2E 測試案例

### 測試檔案: `tests/e2e/auth_clean_test.go`

| 測試名稱 | 說明 | 預期結果 |
|----------|------|----------|
| TestAuth_ValidAccountID | 使用有效帳號 ID | 200 OK |
| TestAuth_MissingHeader | 沒有 X-Account-ID header | 401 Unauthorized |
| TestAuth_InvalidFormat | X-Account-ID 格式不正確 | 401 Unauthorized |
| TestAuth_NonExistentAccount | 帳號不存在 | 401 Unauthorized |
| TestAuth_NegativeAccountID | 負數帳號 ID | 401 Unauthorized |
| TestAuth_ZeroAccountID | 帳號 ID 為 0 | 401 Unauthorized |
| TestAuth_JoinGameRoom_RequiresAuth | 未認證加入遊戲房 | 401 Unauthorized |
| TestAuth_StartGame_RequiresAuth | 未認證開始遊戲 | 401 Unauthorized |

## 錯誤訊息

| 情境 | HTTP 狀態碼 | 錯誤訊息 |
|------|-------------|----------|
| 沒有 X-Account-ID header | 401 | 未提供認證資訊 |
| X-Account-ID 格式不正確 | 401 | 無效的認證資訊 |
| 帳號不存在 | 401 | 帳號不存在 |

## 相關檔案清單

### 新增檔案
- `internal/adapter/http/middleware/auth_middleware.go` - 認證中介軟體
- `tests/e2e/auth_clean_test.go` - 8 個 E2E 測試案例

### 修改檔案
- `internal/adapter/http/handler/game_room_handler.go`
  - 移除舊的 AuthMiddleware 函數
  - 使用新的 middleware.AuthMiddleware
  - GameRoomHandlerOptions 新增 AccountRepo 欄位
- `cmd/app/main.go` - 傳遞 AccountRepo 給 GameRoomHandler
- `tests/e2e/suite_clean_test.go` - 傳遞 AccountRepo 給 GameRoomHandler

## 與 Task 1.1, 1.2 的關聯

- **Task 1.1**: 提供 Account 實體和 AccountRepository
- **Task 1.2**: 使用 AuthMiddleware 保護遊戲房 API
- **Task 1.3**: 完善 AuthMiddleware，增加帳號存在性驗證

## 認證流程

```
Request with X-Account-ID header
           │
           ▼
    ┌──────────────┐
    │ Header 存在？│──No──▶ 401: 未提供認證資訊
    └──────┬───────┘
           │ Yes
           ▼
    ┌──────────────┐
    │ 格式正確？   │──No──▶ 401: 無效的認證資訊
    └──────┬───────┘
           │ Yes
           ▼
    ┌──────────────┐
    │ ID > 0？     │──No──▶ 401: 帳號不存在
    └──────┬───────┘
           │ Yes
           ▼
    ┌──────────────┐
    │ 帳號存在？   │──No──▶ 401: 帳號不存在
    └──────┬───────┘
           │ Yes
           ▼
    Set accountID in context
           │
           ▼
      Continue to Handler
```

## 測試執行結果

```bash
go test -v ./tests/e2e/... -run "TestCleanArchTestSuite/TestAuth"
```

所有 8 個認證測試案例通過：
- TestAuth_ValidAccountID ✓
- TestAuth_MissingHeader ✓
- TestAuth_InvalidFormat ✓
- TestAuth_NonExistentAccount ✓
- TestAuth_NegativeAccountID ✓
- TestAuth_ZeroAccountID ✓
- TestAuth_JoinGameRoom_RequiresAuth ✓
- TestAuth_StartGame_RequiresAuth ✓

## 安全性考量

### 目前實作
- 使用 X-Account-ID header 進行簡易認證
- 適用於開發和測試環境

### 未來改進方向
1. **Token-based 認證**: 使用 JWT 或 Session Token
2. **密碼驗證**: 結合帳號密碼驗證
3. **OAuth 整合**: 支援第三方登入
4. **Rate Limiting**: 限制 API 請求頻率
5. **HTTPS**: 強制使用 HTTPS 傳輸
