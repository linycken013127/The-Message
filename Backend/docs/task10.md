# Task 4.1 - 遊戲查詢實作

## 功能概述

遊戲查詢 API 提供玩家查詢正在進行的遊戲資訊，包含遊戲狀態、階段、當前玩家、情報傳遞狀態、個人資訊和所有玩家的公開資訊。

## API 端點

### GET /api/v1/games/:gameId

查詢遊戲資訊（需要認證）

**Header:**
- `X-Account-ID`: 帳號 ID（必填）

**Path Parameters:**
- `gameId`: 遊戲 ID

**成功回應 (200):**
```json
{
  "gameId": 1,
  "status": "PLAYING",
  "phase": "ACTION",
  "currentPlayerId": 1,
  "hasIntelInTransit": false,
  "myFaction": "潛伏戰線",
  "myRedCount": 1,
  "myBlueCount": 0,
  "myBlackCount": 1,
  "myHandCardCount": 3,
  "players": [
    {
      "id": 1,
      "name": "Player1",
      "alive": true,
      "intelCount": 2
    },
    {
      "id": 2,
      "name": "Player2",
      "alive": true,
      "intelCount": 0
    },
    {
      "id": 3,
      "name": "Player3",
      "alive": false,
      "intelCount": 3
    }
  ]
}
```

**錯誤回應:**
- `401 Unauthorized`: 未提供認證或認證無效
- `404 Not Found`: 遊戲不存在

## 架構設計

### Domain Layer

#### Repository 更新
```
internal/domain/repository/game_player_repository.go
```

新增方法：
- `GetGamePlayerByGameIDAndAccountID`: 根據遊戲 ID 和帳號 ID 取得玩家關聯

### Use Case Layer

#### GameQueryUseCase
```
internal/usecase/game_query_usecase.go
```

**介面:**
```go
type GameQueryUseCase interface {
    GetGameInfo(ctx context.Context, gameID int, accountID int) (*GameInfoResult, error)
}
```

**結果結構:**
```go
type GameInfoResult struct {
    GameID            int                `json:"gameId"`
    Status            string             `json:"status"`
    Phase             string             `json:"phase"`
    CurrentPlayerID   int                `json:"currentPlayerId"`
    HasIntelInTransit bool               `json:"hasIntelInTransit"`
    MyFaction         string             `json:"myFaction"`
    MyRedCount        int                `json:"myRedCount"`
    MyBlueCount       int                `json:"myBlueCount"`
    MyBlackCount      int                `json:"myBlackCount"`
    MyHandCardCount   int                `json:"myHandCardCount"`
    Players           []PlayerPublicInfo `json:"players"`
}

type PlayerPublicInfo struct {
    ID         int    `json:"id"`
    Name       string `json:"name"`
    Alive      bool   `json:"alive"`
    IntelCount int    `json:"intelCount"`
}
```

### Adapter Layer

#### GameQueryHandler
```
internal/adapter/http/handler/game_query_handler.go
```

處理 GET /api/v1/games/:gameId 請求

### Infrastructure Layer

#### GamePlayerRepository 更新
```
internal/infrastructure/persistence/mysql/game_player_repository.go
```

新增 `GetGamePlayerByGameIDAndAccountID` 方法實作

## 測試案例

### E2E 測試
```
tests/e2e/game_query_clean_test.go
```

1. **TestGameQuery_Success**: 成功查詢遊戲資訊
2. **TestGameQuery_WithPersonalInfo**: 查詢結果包含玩家個人資訊
3. **TestGameQuery_Unauthorized**: 未認證無法查詢
4. **TestGameQuery_WithIntelligenceInTransit**: 查詢結果包含情報傳遞狀態
5. **TestGameQuery_WithPlayersInfo**: 查詢結果包含所有玩家公開資訊
6. **TestGameQuery_NotInGame**: 查詢不存在的遊戲

## 資訊可見性規則

### 公開資訊（所有玩家可見）
- 遊戲 ID
- 遊戲狀態（WAITING/PLAYING/ENDED）
- 遊戲階段（ACTION/INTELLIGENCE）
- 當前行動玩家 ID
- 是否有情報正在傳遞
- 所有玩家的：
  - 玩家 ID
  - 玩家名稱
  - 存活狀態
  - 情報區卡片數量

### 私人資訊（僅自己可見）
- 自己的身份（潛伏戰線/軍情處/打醬油）
- 自己的紅色情報數量
- 自己的藍色情報數量
- 自己的黑色情報數量
- 自己的手牌數量

## 查詢流程

```
玩家發送查詢請求
    ↓
驗證認證（X-Account-ID）
    ↓
┌─────────────────────────────────────┐
│ 認證有效?                           │
├─────────────────────────────────────┤
│ 否: → 回傳 401 Unauthorized         │
├─────────────────────────────────────┤
│ 是: → 查詢遊戲                      │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ 遊戲存在?                           │
├─────────────────────────────────────┤
│ 否: → 回傳 404 Not Found            │
├─────────────────────────────────────┤
│ 是: → 組裝回應                      │
└─────────────────────────────────────┘
    ↓
查詢情報傳遞狀態
    ↓
取得所有玩家
    ↓
找出請求者對應的玩家
    ↓
計算每位玩家的情報數量
    ↓
組裝公開資訊 + 個人資訊
    ↓
回傳 200 OK
```

## 相關文件

- Feature File: features/09-game-query.feature
- Task 2.3: 情報階段 - docs/task6.md (情報傳遞狀態)
- Task 3.1: 死亡條件 - docs/task8.md (玩家存活狀態)
- Task 3.2: 勝利條件 - docs/task9.md (遊戲狀態)

## 新增檔案

1. `internal/usecase/game_query_usecase.go` - 遊戲查詢用例
2. `internal/adapter/http/handler/game_query_handler.go` - 遊戲查詢處理器
3. `tests/e2e/game_query_clean_test.go` - E2E 測試

## 修改檔案

1. `internal/domain/repository/game_player_repository.go` - 新增 GetGamePlayerByGameIDAndAccountID 方法
2. `internal/infrastructure/persistence/mysql/game_player_repository.go` - 實作新方法
3. `tests/e2e/suite_clean_test.go` - 註冊新的 Use Case 和 Handler

## 注意事項

1. 認證機制使用 `X-Account-ID` header
2. 個人資訊僅顯示給對應的玩家
3. 情報數量計算包含所有類型（紅/藍/黑）
4. 手牌數量計算僅包含 type 為 "hand" 的卡片
