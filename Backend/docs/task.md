# The Message - 功能實作規劃文件

## 目錄

1. [概述](#概述)
2. [現況分析](#現況分析)
3. [Feature 清單與實作狀態](#feature-清單與實作狀態)
4. [實作優先順序](#實作優先順序)
5. [詳細任務規劃](#詳細任務規劃)

---

## 概述

本文件根據 `features/` 目錄下的 BDD Feature Files 分析目前 Golang 後端實作狀況，並規劃完成所有功能的步驟。

**開發原則：**
- 先寫 E2E 測試，再實作功能 (TDD/BDD)
- 遵循 Clean Architecture 架構
- 每完成一個功能就 commit

---

## 現況分析

### 已實作功能

| 功能 | API | 狀態 | 說明 |
|------|-----|------|------|
| 開始遊戲 | POST /api/v1/games | ⚠️ 部分完成 | 無認證、無遊戲房概念 |
| 初始化玩家 | - | ⚠️ 部分完成 | 只支援 3 人身份配置 |
| 初始化牌組 | - | ✅ 完成 | 支援洗牌、發牌 |
| 抽牌 | - | ✅ 完成 | 每位玩家抽 3 張 |
| 出牌 | POST /api/v1/players/:playerId/player-cards | ⚠️ 部分完成 | 缺少手牌數檢查 |
| 傳情報 | POST /api/v1/player/:playerId/transmit-intelligence | ⚠️ 部分完成 | 缺少卡片類型處理 |
| 接收卡片 | POST /api/v1/players/:playerId/accept | ⚠️ 部分完成 | 缺少拒絕邏輯 |
| 勝利檢查 | - | ⚠️ 部分完成 | 邏輯需修正 |
| 取得玩家手牌 | GET /api/v1/player/:playerId/player-cards/ | ✅ 完成 | - |
| SSE 事件 | GET /api/v1/games/:gameId/events | ✅ 完成 | - |
| 健康檢查 | GET /api/v1/heartbeat | ✅ 完成 | - |

### 未實作功能

| Feature File | 主要缺失 |
|--------------|----------|
| 01-player-management | 完全未實作玩家註冊 API |
| 02-game-room | 缺少遊戲房概念 (建立、加入、開始) |
| 03-game-initialization | 身份牌需支援 3-9 人配置 |
| 04-action-phase | 缺少抽牌、pass API |
| 05-intelligence-phase | 缺少卡片類型傳遞邏輯 |
| 06-intelligence-reception | 缺少拒絕、自動接收邏輯 |
| 07-victory-conditions | 勝利邏輯需修正 |
| 08-death-conditions | 完全未實作死亡機制 |
| 09-game-query | 缺少完整的遊戲查詢 API |

---

## Feature 清單與實作狀態

### 01-player-management.feature

**API:** `POST /api/v1/players`

| Scenario | 實作狀態 |
|----------|----------|
| 玩家名稱沒有重複 - 成功新增玩家 | ❌ 未實作 |
| 玩家無法使用重複名稱 | ❌ 未實作 |
| 玩家名稱最大長度為 10 字元 | ❌ 未實作 |

**需新增：**
- Entity: 獨立的 Player 管理 (不依賴 Game)
- Repository: PlayerRepository.FindByName()
- UseCase: CreatePlayer with validation
- Handler: POST /api/v1/players

---

### 02-game-room.feature

**API:**
- `POST /api/v1/games` (建立遊戲房)
- `POST /api/v1/games/:gameId/join` (加入遊戲房)
- `POST /api/v1/games/:gameId/start` (開始遊戲)

| Scenario | 實作狀態 |
|----------|----------|
| 成功建立遊戲房 (需認證) | ❌ 未實作 |
| 未認證無法建立遊戲房 | ❌ 未實作 |
| 成功加入遊戲房 | ❌ 未實作 |
| 房主可以開始遊戲 | ❌ 未實作 |
| 非房主無法開始遊戲 | ❌ 未實作 |
| 人數少於 3 人時無法開始 | ❌ 未實作 |
| 第 10 位玩家無法加入 | ❌ 未實作 |
| 嘗試重複加入同一遊戲房 | ❌ 未實作 |
| 嘗試加入進行中的遊戲 | ❌ 未實作 |

**需新增/修改：**
- Entity: Game 新增 HostPlayerID, MaxPlayers, Status (WAITING/PLAYING/ENDED)
- Entity: GamePlayer 關聯表
- Repository: 遊戲房相關查詢
- UseCase: CreateGameRoom, JoinGameRoom, StartGame
- Handler: 三個新 API
- Middleware: 認證機制

---

### 03-game-initialization.feature

**API:** `POST /api/v1/games/:gameId/start`

| Scenario | 實作狀態 |
|----------|----------|
| 不同人數的身份配置 (3人/5人) | ⚠️ 只支援 3 人 |
| 發放初始手牌 (每人 3 張) | ✅ 已實作 |
| 建立完整牌堆 (45 張) | ✅ 已實作 |
| 設定第一位行動玩家 | ✅ 已實作 |

**需修改：**
- UseCase: InitIdentityCards 支援 3-9 人配置

**身份配置表：**
```
3人: 潛伏戰線 1, 軍情處 1, 打醬油的 1
5人: 潛伏戰線 2, 軍情處 2, 打醬油的 1
6人: 潛伏戰線 2, 軍情處 2, 打醬油的 2
7人: 潛伏戰線 2, 軍情處 2, 打醬油的 3
8人: 潛伏戰線 3, 軍情處 3, 打醬油的 2
9人: 潛伏戰線 3, 軍情處 3, 打醬油的 3
```

---

### 04-action-phase.feature

**API:**
- `POST /api/v1/games/:gameId/actions/draw` (抽牌)
- `POST /api/v1/games/:gameId/actions/play-card` (出功能牌)
- `POST /api/v1/games/:gameId/actions/pass` (跳過)

| Scenario | 實作狀態 |
|----------|----------|
| 當前玩家抽 2 張牌 | ❌ 未實作 |
| 牌堆只剩 1 張牌 | ❌ 未實作 |
| 玩家出功能牌 | ⚠️ 部分完成 |
| 當前行動玩家不能出掉最後一張手牌 | ❌ 未實作 |
| 嘗試出不在手牌中的卡牌 | ⚠️ 部分完成 |
| 所有玩家 pass 後進入情報階段 | ❌ 未實作 |

**需新增/修改：**
- Entity: Game 新增 Phase 欄位 (ACTION/INTELLIGENCE)
- UseCase: DrawCards, PlayFunctionCard, Pass
- Handler: 三個新 API
- 業務邏輯: Pass 計數、階段切換

---

### 05-intelligence-phase.feature

**API:**
- `POST /api/v1/games/:gameId/intelligence/pass-card` (傳遞情報牌)

| Scenario | 實作狀態 |
|----------|----------|
| 選擇手牌作為情報 | ⚠️ 部分完成 |
| 傳遞密電卡牌 (蓋牌向右) | ❌ 未實作 |
| 傳遞文本卡牌 (明牌向右) | ❌ 未實作 |
| 傳遞直達卡牌給指定玩家 | ❌ 未實作 |
| 直達卡牌不能指定自己 | ❌ 未實作 |
| 直達卡牌只能指定存活玩家 | ❌ 未實作 |
| 只有當前行動玩家可以傳遞情報 | ⚠️ 部分完成 |
| 卡牌必須在手牌中才能傳遞 | ⚠️ 部分完成 |

**需新增/修改：**
- Entity: IntelligenceTransfer (情報傳遞狀態)
  - ID, GameID, CardID, SenderPlayerID
  - CurrentTargetPlayerID, OriginalTargetPlayerID
  - FaceUp, Status (IN_TRANSIT/ACCEPTED/REJECTED)
- Repository: IntelligenceTransferRepository
- UseCase: PassIntelligenceCard (依卡片類型處理)
- Handler: 新 API

---

### 06-intelligence-reception.feature

**API:**
- `POST /api/v1/games/:gameId/intelligence/accept` (接收情報)
- `POST /api/v1/games/:gameId/intelligence/reject` (拒絕情報)

| Scenario | 實作狀態 |
|----------|----------|
| 玩家接收情報 | ⚠️ 部分完成 |
| 玩家拒絕接收密電情報 | ❌ 未實作 |
| 玩家拒絕接收文本情報 | ❌ 未實作 |
| 密電情報回到發送者自動接收 | ❌ 未實作 |
| 直達情報被拒絕後自動歸屬發送者 | ❌ 未實作 |
| 非目標玩家嘗試接收情報 | ❌ 未實作 |
| 進入下一位玩家的回合 | ⚠️ 部分完成 |

**需新增/修改：**
- UseCase: AcceptIntelligence, RejectIntelligence
- 業務邏輯: 傳遞方向計算、自動接收檢查
- Handler: 兩個新 API

---

### 07-victory-conditions.feature

**系統自動觸發**

| Scenario | 實作狀態 |
|----------|----------|
| 潛伏戰線收集 3 張紅色情報獲勝 | ⚠️ 邏輯錯誤 |
| 軍情處收集 3 張藍色情報獲勝 | ⚠️ 邏輯錯誤 |
| 打醬油的跟隨獲勝陣營 | ❌ 未實作 |
| 未達成勝利條件時遊戲繼續 | ⚠️ 部分完成 |

**需修改：**
- UseCase: CheckWin 邏輯修正
  - 潛伏戰線: 紅色情報 >= 3
  - 軍情處: 藍色情報 >= 3
  - 打醬油的: 不能單獨觸發勝利
- Entity: Game 新增 Winner 欄位

---

### 08-death-conditions.feature

**系統自動觸發**

| Scenario | 實作狀態 |
|----------|----------|
| 玩家收集 3 張黑色情報死亡 | ❌ 未實作 |
| 玩家收集超過 3 張黑色情報死亡 | ❌ 未實作 |
| 死亡玩家被跳過 | ❌ 未實作 |
| 情報跳過死亡玩家 | ❌ 未實作 |
| 黑色情報少於 3 張時玩家保持存活 | ❌ 未實作 |
| 死亡條件優先於勝利條件 | ❌ 未實作 |
| 所有玩家死亡時遊戲結束為平局 | ❌ 未實作 |

**需新增：**
- UseCase: CheckDeath, KillPlayer
- 業務邏輯: 死亡優先於勝利、平局判定
- NextPlayer 邏輯需跳過死亡玩家

---

### 09-game-query.feature

**API:** `GET /api/v1/games/:gameId`

| Scenario | 實作狀態 |
|----------|----------|
| 成功查詢遊戲資訊 | ⚠️ 部分完成 |
| 查詢結果包含可見資訊 | ❌ 未實作 |
| 未提供認證 | ❌ 未實作 |
| 查詢包含完整遊戲狀態 | ❌ 未實作 |

**需新增/修改：**
- Handler: GetGameInfo 重構
- 回應需包含:
  - 遊戲狀態 (status, phase, currentPlayerId, hasIntelInTransit)
  - 玩家個人資訊 (myFaction, myRedCount, myBlueCount, myBlackCount, myHandCardCount)
  - 所有玩家公開資訊

---

## 實作優先順序

基於依賴關係，建議按以下順序實作：

### Phase 1: 基礎架構 (Foundation)

1. **Task 1.1: 玩家管理**
   - E2E Test: `player_management_test.go`
   - 實作 POST /api/v1/players
   - 驗證名稱唯一性、長度限制

2. **Task 1.2: 遊戲房管理**
   - E2E Test: `game_room_test.go`
   - 修改 Game Entity (加入 HostPlayerID, Status, MaxPlayers)
   - 實作 POST /api/v1/games (建立遊戲房)
   - 實作 POST /api/v1/games/:gameId/join
   - 實作 POST /api/v1/games/:gameId/start

3. **Task 1.3: 認證機制**
   - 實作簡易認證 Middleware
   - 整合到遊戲房相關 API

### Phase 2: 遊戲流程 (Game Flow)

4. **Task 2.1: 遊戲初始化完善**
   - E2E Test: `game_initialization_test.go`
   - 完善 3-9 人身份配置
   - 確認牌堆 45 張

5. **Task 2.2: 行動階段**
   - E2E Test: `action_phase_test.go`
   - 實作 POST /api/v1/games/:gameId/actions/draw
   - 實作 POST /api/v1/games/:gameId/actions/play-card
   - 實作 POST /api/v1/games/:gameId/actions/pass
   - 實作 Pass 計數與階段切換

6. **Task 2.3: 情報階段**
   - E2E Test: `intelligence_phase_test.go`
   - 新增 IntelligenceTransfer Entity
   - 實作 POST /api/v1/games/:gameId/intelligence/pass-card
   - 處理密電/文本/直達三種類型

7. **Task 2.4: 情報接收**
   - E2E Test: `intelligence_reception_test.go`
   - 實作 POST /api/v1/games/:gameId/intelligence/accept
   - 實作 POST /api/v1/games/:gameId/intelligence/reject
   - 實作自動接收邏輯

### Phase 3: 勝負判定 (Win/Lose)

8. **Task 3.1: 死亡條件**
   - E2E Test: `death_conditions_test.go`
   - 實作 CheckDeath UseCase
   - 修改 NextPlayer 跳過死亡玩家
   - 修改情報傳遞跳過死亡玩家

9. **Task 3.2: 勝利條件**
   - E2E Test: `victory_conditions_test.go`
   - 修正 CheckWin UseCase
   - 實作平局判定

### Phase 4: 查詢與優化 (Query & Polish)

10. **Task 4.1: 遊戲查詢**
    - E2E Test: `game_query_test.go`
    - 重構 GET /api/v1/games/:gameId
    - 實作資訊可見性控制

---

## 詳細任務規劃

### Task 1.1: 玩家管理

**E2E 測試檔案:** `tests/e2e/player_management_clean_test.go`

```go
// 測試案例
func TestCreatePlayer_Success()
func TestCreatePlayer_DuplicateName()
func TestCreatePlayer_NameTooLong()
```

**實作步驟:**

1. **Domain Layer**
   - 修改 `internal/domain/entity/player.go`
     - 新增 ExternalID 欄位 (用於外部識別)
   - 修改 `internal/domain/repository/player_repository.go`
     - 新增 `FindByName(ctx, name) (*Player, error)`
     - 新增 `ExistsByName(ctx, name) (bool, error)`

2. **Infrastructure Layer**
   - 修改 `internal/infrastructure/persistence/mysql/player_repository.go`
     - 實作 FindByName, ExistsByName

3. **Use Case Layer**
   - 修改 `internal/usecase/player_usecase.go`
     - 新增 `RegisterPlayer(ctx, name) (*Player, error)`
     - 加入名稱驗證邏輯

4. **Adapter Layer**
   - 修改 `internal/adapter/http/handler/player_handler.go`
     - 新增 `RegisterPlayer` handler
   - 修改 `internal/adapter/http/request/request.go`
     - 新增 `RegisterPlayerRequest`

5. **Routes**
   - 在 main.go 註冊 `POST /api/v1/players`

---

### Task 1.2: 遊戲房管理

**E2E 測試檔案:** `tests/e2e/game_room_clean_test.go`

```go
// 測試案例
func TestCreateGameRoom_Success()
func TestCreateGameRoom_Unauthorized()
func TestJoinGameRoom_Success()
func TestJoinGameRoom_Full()
func TestJoinGameRoom_AlreadyJoined()
func TestJoinGameRoom_GameStarted()
func TestStartGame_AsHost()
func TestStartGame_NotHost()
func TestStartGame_NotEnoughPlayers()
```

**實作步驟:**

1. **Domain Layer**
   - 修改 `internal/domain/entity/game.go`
     ```go
     type Game struct {
         // 新增欄位
         HostPlayerID   int
         MaxPlayers     int
         CurrentPlayers int
     }

     // 新增常數
     const (
         GameStatusWaiting = "WAITING"
         GameStatusPlaying = "PLAYING"
         GameStatusEnded   = "ENDED"
     )
     ```
   - 新增 `internal/domain/entity/game_player.go`
     ```go
     type GamePlayer struct {
         ID       int
         GameID   int
         PlayerID int
         JoinedAt time.Time
     }
     ```

2. **Infrastructure Layer**
   - 新增 `game_player_repository.go`
   - 新增資料庫 migration

3. **Use Case Layer**
   - 修改 `game_usecase.go`
     - `CreateGameRoom(ctx, hostPlayerID, maxPlayers) (*Game, error)`
     - `JoinGameRoom(ctx, gameID, playerID) error`
     - `StartGame(ctx, gameID, playerID) error`

4. **Adapter Layer**
   - 修改 `game_handler.go`
     - `CreateGameRoom` handler
     - `JoinGameRoom` handler
     - `StartGame` handler

---

### Task 2.2: 行動階段

**E2E 測試檔案:** `tests/e2e/action_phase_clean_test.go`

```go
// 測試案例
func TestDrawCards_Success()
func TestDrawCards_DeckNotEnough()
func TestPlayFunctionCard_Success()
func TestPlayFunctionCard_LastCard()
func TestPlayFunctionCard_CardNotInHand()
func TestPass_AllPlayersPass()
```

**實作步驟:**

1. **Domain Layer**
   - 修改 `internal/domain/entity/game.go`
     ```go
     // 新增欄位
     Phase string // ACTION, INTELLIGENCE

     // 新增常數
     const (
         GamePhaseAction       = "ACTION"
         GamePhaseIntelligence = "INTELLIGENCE"
     )
     ```
   - 新增 `internal/domain/entity/action_pass.go`
     ```go
     type ActionPass struct {
         ID       int
         GameID   int
         PlayerID int
         Round    int
     }
     ```

2. **Use Case Layer**
   - 修改 `game_usecase.go`
     - `DrawCardsForCurrentPlayer(ctx, gameID) ([]Card, error)`
     - `Pass(ctx, gameID, playerID) error`
     - `CheckAllPassed(ctx, gameID) (bool, error)`

3. **Adapter Layer**
   - 新增 handler 方法:
     - `DrawCards`
     - `PlayFunctionCard`
     - `Pass`

---

### Task 2.3: 情報階段

**E2E 測試檔案:** `tests/e2e/intelligence_phase_clean_test.go`

```go
// 測試案例
func TestPassIntelligence_Secret()
func TestPassIntelligence_Text()
func TestPassIntelligence_Direct()
func TestPassIntelligence_DirectToSelf()
func TestPassIntelligence_DirectToDeadPlayer()
func TestPassIntelligence_NotCurrentPlayer()
func TestPassIntelligence_CardNotInHand()
```

**實作步驟:**

1. **Domain Layer**
   - 新增 `internal/domain/entity/intelligence_transfer.go`
     ```go
     type IntelligenceTransfer struct {
         ID                     int
         GameID                 int
         CardID                 int
         SenderPlayerID         int
         CurrentTargetPlayerID  int
         OriginalTargetPlayerID int
         FaceUp                 bool
         Status                 string // IN_TRANSIT, COMPLETED
         CreatedAt              time.Time
     }

     const (
         IntelligenceStatusInTransit = "IN_TRANSIT"
         IntelligenceStatusCompleted = "COMPLETED"
     )
     ```

2. **Infrastructure Layer**
   - 新增 `intelligence_transfer_repository.go`
   - 新增資料庫 migration

3. **Use Case Layer**
   - 新增 `intelligence_usecase.go`
     - `PassIntelligenceCard(ctx, gameID, playerID, cardID, targetPlayerID) error`
     - 內部處理三種卡片類型

---

### Task 2.4: 情報接收

**E2E 測試檔案:** `tests/e2e/intelligence_reception_clean_test.go`

```go
// 測試案例
func TestAcceptIntelligence_Success()
func TestRejectIntelligence_Secret()
func TestRejectIntelligence_Text()
func TestRejectIntelligence_BackToSender()
func TestRejectIntelligence_Direct()
func TestAcceptIntelligence_NotTargetPlayer()
```

**實作步驟:**

1. **Use Case Layer**
   - 修改 `intelligence_usecase.go`
     - `AcceptIntelligence(ctx, gameID, playerID) error`
     - `RejectIntelligence(ctx, gameID, playerID) error`
     - `GetNextTargetPlayer(ctx, game, currentTarget) (*Player, error)`
     - `AutoAcceptIntelligence(ctx, transfer) error`

2. **Adapter Layer**
   - 新增 handler:
     - `AcceptIntelligence`
     - `RejectIntelligence`

---

### Task 3.1: 死亡條件

**E2E 測試檔案:** `tests/e2e/death_conditions_clean_test.go`

```go
// 測試案例
func TestDeath_ThreeBlackIntelligence()
func TestDeath_MoreThanThreeBlack()
func TestDeath_SkipDeadPlayerTurn()
func TestDeath_SkipDeadPlayerIntelligence()
func TestDeath_NotEnoughBlack()
func TestDeath_PriorityOverVictory()
func TestDeath_AllPlayersDead()
```

**實作步驟:**

1. **Use Case Layer**
   - 新增 `death_usecase.go`
     - `CheckDeath(ctx, playerID) (bool, error)`
     - `KillPlayer(ctx, playerID) error`
     - `CheckAllPlayersDead(ctx, gameID) (bool, error)`
   - 修改 `game_usecase.go`
     - `NextPlayer` 跳過死亡玩家
   - 修改 `intelligence_usecase.go`
     - 情報傳遞跳過死亡玩家

---

### Task 3.2: 勝利條件

**E2E 測試檔案:** `tests/e2e/victory_conditions_clean_test.go`

```go
// 測試案例
func TestVictory_UndercoverFront_ThreeRed()
func TestVictory_UndercoverFront_MoreThanThreeRed()
func TestVictory_MilitaryAgency_ThreeBlue()
func TestVictory_Bystander_CannotTrigger()
func TestVictory_NotReached()
func TestVictory_Draw()
```

**實作步驟:**

1. **Domain Layer**
   - 修改 `internal/domain/entity/game.go`
     - 新增 Winner 欄位

2. **Use Case Layer**
   - 修改 `player_usecase.go`
     - 重構 `CheckWin` 邏輯
     ```go
     // 正確邏輯:
     // 1. 計算每位玩家的紅/藍/黑情報數
     // 2. 潛伏戰線: 紅 >= 3 觸發勝利
     // 3. 軍情處: 藍 >= 3 觸發勝利
     // 4. 打醬油的: 不能單獨觸發
     // 5. 死亡條件優先於勝利條件
     ```

---

### Task 4.1: 遊戲查詢

**E2E 測試檔案:** `tests/e2e/game_query_clean_test.go`

```go
// 測試案例
func TestGetGameInfo_Success()
func TestGetGameInfo_WithPersonalInfo()
func TestGetGameInfo_Unauthorized()
func TestGetGameInfo_WithIntelligenceInTransit()
```

**實作步驟:**

1. **Adapter Layer**
   - 重構 `game_handler.go` 的 `GetGame`
   - 新增 Response DTO:
     ```go
     type GameInfoResponse struct {
         GameID            int                   `json:"gameId"`
         Status            string                `json:"status"`
         Phase             string                `json:"phase"`
         CurrentPlayerID   int                   `json:"currentPlayerId"`
         HasIntelInTransit bool                  `json:"hasIntelInTransit"`
         MyFaction         string                `json:"myFaction,omitempty"`
         MyRedCount        int                   `json:"myRedCount,omitempty"`
         MyBlueCount       int                   `json:"myBlueCount,omitempty"`
         MyBlackCount      int                   `json:"myBlackCount,omitempty"`
         MyHandCardCount   int                   `json:"myHandCardCount,omitempty"`
         Players           []PlayerPublicInfo    `json:"players"`
     }

     type PlayerPublicInfo struct {
         ID           int    `json:"id"`
         Name         string `json:"name"`
         Alive        bool   `json:"alive"`
         IntelCount   int    `json:"intelCount"`
     }
     ```

---

## 資料庫 Migration 清單

需新增的 Migration:

1. `xxx_add_game_room_fields.sql`
   - games 表新增: host_player_id, max_players, current_players, phase, winner

2. `xxx_create_game_players.sql`
   - 新增 game_players 表

3. `xxx_create_intelligence_transfers.sql`
   - 新增 intelligence_transfers 表

4. `xxx_create_action_passes.sql`
   - 新增 action_passes 表

---

## 測試執行順序

```bash
# Phase 1
go test -v ./tests/e2e/player_management_clean_test.go
go test -v ./tests/e2e/game_room_clean_test.go

# Phase 2
go test -v ./tests/e2e/game_initialization_clean_test.go
go test -v ./tests/e2e/action_phase_clean_test.go
go test -v ./tests/e2e/intelligence_phase_clean_test.go
go test -v ./tests/e2e/intelligence_reception_clean_test.go

# Phase 3
go test -v ./tests/e2e/death_conditions_clean_test.go
go test -v ./tests/e2e/victory_conditions_clean_test.go

# Phase 4
go test -v ./tests/e2e/game_query_clean_test.go

# 全部測試
go test -v ./tests/e2e/...
```

---

## Commit 規範

每完成一個 Task 後進行 commit:

```
feat: implement player management API
feat: implement game room management
feat: implement action phase with draw/pass
feat: implement intelligence phase
feat: implement intelligence reception
feat: implement death conditions
fix: correct victory conditions logic
feat: implement game query API
```
