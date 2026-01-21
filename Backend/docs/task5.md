# Task 2.2: 行動階段

## 功能概述

實作遊戲行動階段的三個主要 API：
1. **抽牌** (Draw Cards) - 當前玩家從牌堆抽 2 張牌
2. **出功能牌** (Play Function Card) - 當前玩家出一張功能牌
3. **跳過** (Pass) - 當前玩家跳過行動

當所有玩家都跳過後，遊戲進入情報階段。

## API 端點

### 1. 抽牌 API

**端點:** `POST /api/v1/games/:gameId/actions/draw`

**請求:**
```json
{
  "player_id": 1
}
```

**成功回應:**
```json
{
  "success": true,
  "card_ids": [1, 2],
  "count": 2
}
```

**規則:**
- 只有當前玩家可以抽牌
- 每次抽 2 張，牌堆不足時抽剩餘的
- 牌堆為空時回傳錯誤

### 2. 出功能牌 API

**端點:** `POST /api/v1/games/:gameId/actions/play-card`

**請求:**
```json
{
  "player_id": 1,
  "card_id": 5
}
```

**成功回應:**
```json
{
  "success": true,
  "card_id": 5
}
```

**規則:**
- 只有當前玩家可以出牌
- 卡片必須在手牌中
- 不能打出最後一張手牌

### 3. 跳過 API

**端點:** `POST /api/v1/games/:gameId/actions/pass`

**請求:**
```json
{
  "player_id": 1
}
```

**成功回應:**
```json
{
  "success": true,
  "all_passed": false,
  "phase_changed": false,
  "new_phase": "ACTION"
}
```

**當所有玩家都跳過:**
```json
{
  "success": true,
  "all_passed": true,
  "phase_changed": true,
  "new_phase": "INTELLIGENCE"
}
```

## 遊戲階段

新增 `Phase` 欄位到 Game Entity：

| 階段 | 值 | 說明 |
|------|-----|------|
| 行動階段 | `ACTION` | 玩家可以抽牌、出牌、跳過 |
| 情報階段 | `INTELLIGENCE` | 玩家傳遞情報 |

## 實作細節

### 1. Domain Layer 新增

**Game Entity 新增欄位 (`game.go`):**
```go
type Game struct {
    // ... 原有欄位
    Phase string  // 新增：遊戲階段
}

const (
    GamePhaseAction       = "ACTION"       // 行動階段
    GamePhaseIntelligence = "INTELLIGENCE" // 情報階段
)
```

**新增 ActionPass Entity (`action_pass.go`):**
```go
type ActionPass struct {
    ID        int
    GameID    int
    PlayerID  int
    Round     int
    CreatedAt time.Time
}
```

**新增錯誤定義:**
```go
var (
    ErrNotYourTurn          = errors.New("還沒輪到你")
    ErrNotInActionPhase     = errors.New("不在行動階段")
    ErrDeckNotEnough        = errors.New("牌堆不足")
    ErrCannotPlayLastCard   = errors.New("不能打出最後一張手牌")
    ErrCardNotInHand        = errors.New("這張牌不在你的手牌中")
)
```

### 2. Infrastructure Layer

**新增 Migration (`000010_add_phase_and_action_passes.up.sql`):**
```sql
-- 為 games 表新增 phase 欄位
ALTER TABLE games
    ADD COLUMN phase VARCHAR(20) NOT NULL DEFAULT 'ACTION' AFTER status;

-- 建立 action_passes 表
CREATE TABLE action_passes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    game_id INT NOT NULL,
    player_id INT NOT NULL,
    round INT NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players (id) ON DELETE CASCADE,
    UNIQUE KEY uk_game_player_round (game_id, player_id, round)
);
```

**新增 ActionPassRepository:**
- `CreateActionPass` - 建立跳過記錄
- `CountActionPassesByGameIDAndRound` - 計算跳過次數
- `DeleteActionPassesByGameIDAndRound` - 清除跳過記錄
- `ExistsByGameIDPlayerIDAndRound` - 檢查是否已跳過

### 3. Use Case Layer

**ActionPhaseUseCase (`action_phase_usecase.go`):**

```go
type ActionPhaseUseCase interface {
    // DrawCards 當前玩家抽牌（抽 2 張）
    DrawCards(ctx context.Context, gameID int, playerID int) ([]*entity.Card, error)
    // PlayFunctionCard 出功能牌
    PlayFunctionCard(ctx context.Context, gameID int, playerID int, cardID int) (*entity.Card, error)
    // Pass 跳過行動
    Pass(ctx context.Context, gameID int, playerID int) (bool, error)
}
```

**抽牌邏輯:**
1. 驗證遊戲狀態（進行中、行動階段）
2. 驗證是否為當前玩家
3. 從牌堆抽 2 張（不足時抽剩餘的）
4. 清除該回合的 pass 記錄

**出牌邏輯:**
1. 驗證遊戲狀態
2. 驗證是否為當前玩家
3. 檢查卡片是否在手牌中
4. 檢查手牌數量 > 1（不能打最後一張）
5. 刪除手牌
6. 清除 pass 記錄

**跳過邏輯:**
1. 驗證遊戲狀態
2. 驗證是否為當前玩家
3. 記錄 pass
4. 檢查是否所有存活玩家都已 pass
5. 若全部 pass，切換到情報階段
6. 若未全部 pass，切換到下一位玩家

### 4. Adapter Layer

**ActionPhaseHandler (`action_phase_handler.go`):**
- `DrawCards` - 處理抽牌請求
- `PlayFunctionCard` - 處理出牌請求
- `Pass` - 處理跳過請求

## 測試案例

新增 11 個 E2E 測試案例於 `tests/e2e/action_phase_clean_test.go`：

### 抽牌測試
| 測試案例 | 說明 |
|----------|------|
| `TestDrawCards_Success` | 成功抽 2 張牌 |
| `TestDrawCards_NotCurrentPlayer` | 非當前玩家抽牌失敗 |
| `TestDrawCards_DeckNotEnough` | 牌堆只剩 1 張時抽 1 張 |
| `TestDrawCards_EmptyDeck` | 牌堆為空時失敗 |

### 出牌測試
| 測試案例 | 說明 |
|----------|------|
| `TestPlayFunctionCard_Success` | 成功出牌 |
| `TestPlayFunctionCard_LastCard` | 不能打最後一張手牌 |
| `TestPlayFunctionCard_CardNotInHand` | 卡片不在手牌中 |
| `TestPlayFunctionCard_NotCurrentPlayer` | 非當前玩家出牌失敗 |

### 跳過測試
| 測試案例 | 說明 |
|----------|------|
| `TestPass_Success` | 成功跳過 |
| `TestPass_AllPlayersPass` | 所有人跳過後進入情報階段 |
| `TestPass_NotCurrentPlayer` | 非當前玩家跳過失敗 |

## 架構關係

```
ActionPhaseHandler
       ↓
ActionPhaseUseCase
       │
       ├──→ GameRepository (GetGameById, UpdateGame)
       │
       ├──→ PlayerRepository (GetPlayerById)
       │
       ├──→ PlayerCardRepository (CountHandCards, DeletePlayerCard)
       │
       ├──→ DeckUseCase (GetDecksByGameId, DeleteDeckFromGame)
       │
       └──→ ActionPassRepository (CreateActionPass, CountByRound)
```

## 異動檔案

### 新增
- `internal/domain/entity/action_pass.go` - ActionPass Entity
- `internal/domain/repository/action_pass_repository.go` - Repository 介面
- `internal/infrastructure/persistence/mysql/action_pass_repository.go` - MySQL 實作
- `internal/infrastructure/persistence/mysql/model/action_pass_model.go` - GORM Model
- `internal/usecase/action_phase_usecase.go` - 行動階段 UseCase
- `internal/adapter/http/handler/action_phase_handler.go` - HTTP Handler
- `database/migrations/000010_add_phase_and_action_passes.up.sql` - Migration
- `database/migrations/000010_add_phase_and_action_passes.down.sql` - Rollback
- `tests/e2e/action_phase_clean_test.go` - E2E 測試

### 修改
- `internal/domain/entity/game.go` - 新增 Phase 欄位與階段常數
- `internal/infrastructure/persistence/mysql/model/game_model.go` - 新增 Phase 欄位
- `internal/domain/repository/player_card_repository.go` - 新增 CountHandCardsByPlayerID
- `internal/infrastructure/persistence/mysql/player_card_repository.go` - 實作 CountHandCardsByPlayerID
- `internal/usecase/game_room_usecase.go` - StartGame 設定 Phase 為 ACTION
- `cmd/app/main.go` - 註冊新 Handler
- `tests/e2e/suite_clean_test.go` - 新增依賴注入
- `tests/e2e/player_clean_test.go` - 修正舊測試使用固定 3 人

## 測試結果

```
=== RUN   TestCleanArchTestSuite
--- PASS: TestCleanArchTestSuite (34.78s)
    --- PASS: TestCleanArchTestSuite/TestDrawCards_DeckNotEnough (0.88s)
    --- PASS: TestCleanArchTestSuite/TestDrawCards_EmptyDeck (0.88s)
    --- PASS: TestCleanArchTestSuite/TestDrawCards_NotCurrentPlayer (0.74s)
    --- PASS: TestCleanArchTestSuite/TestDrawCards_Success (0.75s)
    --- PASS: TestCleanArchTestSuite/TestPass_AllPlayersPass (0.71s)
    --- PASS: TestCleanArchTestSuite/TestPass_NotCurrentPlayer (0.69s)
    --- PASS: TestCleanArchTestSuite/TestPass_Success (0.69s)
    --- PASS: TestCleanArchTestSuite/TestPlayFunctionCard_CardNotInHand (0.71s)
    --- PASS: TestCleanArchTestSuite/TestPlayFunctionCard_LastCard (0.69s)
    --- PASS: TestCleanArchTestSuite/TestPlayFunctionCard_NotCurrentPlayer (0.82s)
    --- PASS: TestCleanArchTestSuite/TestPlayFunctionCard_Success (0.74s)
    ... (其他 40 個測試也通過)
PASS
ok      github.com/Game-as-a-Service/The-Message/tests/e2e
```

全部 51 個測試通過。
