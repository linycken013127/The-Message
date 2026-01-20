# Task 1.1: 玩家管理

## 概述

實作玩家註冊 API，允許玩家使用唯一名稱註冊帳號。

**對應 Feature:** `features/01-player-management.feature`

**API:** `POST /api/v1/players`

## 完成的 Scenarios

| Scenario | 狀態 |
|----------|------|
| 玩家名稱沒有重複 - 成功新增玩家 | ✅ 完成 |
| 玩家無法使用重複名稱 | ✅ 完成 |
| 玩家名稱最大長度為 10 字元 | ✅ 完成 |

## API 規格

### POST /api/v1/players

**請求:**
```json
{
  "playerName": "Alice"
}
```

**成功回應 (200):**
```json
{
  "playerId": 1,
  "playerName": "Alice"
}
```

**失敗回應 (400):**
```json
{
  "message": "玩家名稱已存在"
}
```

### 驗證規則

| 規則 | 錯誤訊息 |
|------|----------|
| 名稱不可為空 | 玩家名稱不可為空 |
| 名稱最大 10 字元 | 玩家名稱最大長度為 10 字元 |
| 名稱不可重複（不區分大小寫） | 玩家名稱已存在 |

## 實作架構

依照 Clean Architecture 分層實作：

```
┌─────────────────────────────────────────────────────────────┐
│                    Adapter Layer                             │
│  handler/account_handler.go                                  │
│  request/request.go (RegisterPlayerRequest)                  │
│  response/response.go (RegisterPlayerResponse)               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Use Case Layer                            │
│  usecase/account_usecase.go                                  │
│  - RegisterAccount(ctx, name) (*Account, error)              │
│  - GetAccountById(ctx, id) (*Account, error)                 │
│  - GetAccountByName(ctx, name) (*Account, error)             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                              │
│  entity/account.go                                           │
│  - Account struct                                            │
│  - NewAccount(name) factory method                           │
│  - ValidateAccountName(name) validation                      │
│  repository/account_repository.go (interface)                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                Infrastructure Layer                          │
│  persistence/mysql/account_repository.go                     │
│  persistence/mysql/model/account_model.go                    │
│  database/migrations/000008_create_accounts_table.up.sql     │
└─────────────────────────────────────────────────────────────┘
```

## 新增檔案

### Domain Layer

1. **`internal/domain/entity/account.go`**
   - `Account` 結構體
   - `NewAccount()` 工廠方法
   - `ValidateAccountName()` 驗證函式
   - 錯誤常數定義

2. **`internal/domain/repository/account_repository.go`**
   - `AccountRepository` 介面定義

### Infrastructure Layer

3. **`database/migrations/000008_create_accounts_table.up.sql`**
   ```sql
   CREATE TABLE accounts (
     id         INT AUTO_INCREMENT PRIMARY KEY,
     name       VARCHAR(50) NOT NULL,
     created_at DATETIME NOT NULL,
     updated_at DATETIME NOT NULL,
     deleted_at DATETIME,
     UNIQUE KEY uk_accounts_name (name)
   );
   ```

4. **`database/migrations/000008_create_accounts_table.down.sql`**

5. **`internal/infrastructure/persistence/mysql/model/account_model.go`**
   - GORM Model 定義
   - `ToEntity()` 轉換方法
   - `AccountModelFromEntity()` 轉換方法

6. **`internal/infrastructure/persistence/mysql/account_repository.go`**
   - `CreateAccount()` 實作
   - `GetAccountById()` 實作
   - `GetAccountByName()` 實作
   - `ExistsByName()` 實作（不區分大小寫）

### Use Case Layer

7. **`internal/usecase/account_usecase.go`**
   - `AccountUseCase` 介面
   - `RegisterAccount()` 業務邏輯

### Adapter Layer

8. **`internal/adapter/http/handler/account_handler.go`**
   - `RegisterPlayer` HTTP handler

9. **`internal/adapter/http/request/request.go`** (修改)
   - 新增 `RegisterPlayerRequest`

10. **`internal/adapter/http/response/response.go`** (修改)
    - 新增 `RegisterPlayerResponse`

## E2E 測試

**測試檔案:** `tests/e2e/player_management_clean_test.go`

| 測試案例 | 說明 |
|----------|------|
| `TestRegisterPlayer_Success` | 成功註冊玩家 |
| `TestRegisterPlayer_DuplicateName` | 名稱重複失敗 |
| `TestRegisterPlayer_NameTooLong` | 名稱超過 10 字元失敗 |
| `TestRegisterPlayer_NameExactly10Chars` | 名稱剛好 10 字元成功 |
| `TestRegisterPlayer_EmptyName` | 空名稱失敗 |
| `TestRegisterPlayer_NameWithSpaces` | 名稱含空格成功 |
| `TestRegisterPlayer_DuplicateNameCaseInsensitive` | 不區分大小寫的重複檢查 |
| `TestRegisterPlayer_TrimWhitespace` | 名稱前後空白自動 trim |

## 執行測試

```bash
# 執行所有 e2e 測試
go test -v ./tests/e2e/...

# 只執行玩家管理測試
go test -v ./tests/e2e/... -run "TestRegisterPlayer"
```

## 設計決策

### 為什麼新增 Account Entity？

原有的 `Player` entity 是綁定在遊戲 (`Game`) 內的玩家實例，需要 `game_id` 外鍵。但玩家管理功能需要獨立於遊戲存在的帳號概念。

- **Account**: 代表註冊的使用者帳號，獨立於遊戲
- **Player**: 代表在特定遊戲中的玩家實例

### 名稱驗證規則

1. **長度限制**: 使用 `[]rune` 計算字元數，支援中文等 Unicode 字元
2. **唯一性**: 使用 SQL `LOWER()` 函式進行不區分大小寫比較
3. **空白處理**: 自動 trim 前後空白

## 後續擴展

此 Account entity 將作為後續功能的基礎：

- Task 1.2: 遊戲房管理 - 使用 Account ID 作為認證身份
- Task 1.3: 認證機制 - 基於 Account 實作 JWT 或其他認證方式
