# Task 1.2: 遊戲房管理

## 功能概述

實作遊戲房管理功能，包含建立遊戲房、加入遊戲房、開始遊戲三個主要 API。

### Feature 對應
- `02-game-room.feature`

### API 端點
| Method | Path | 說明 | 認證 |
|--------|------|------|------|
| POST | `/api/v1/games` | 建立遊戲房 | 需要 X-Account-ID |
| POST | `/api/v1/games/:gameId/join` | 加入遊戲房 | 需要 X-Account-ID |
| POST | `/api/v1/games/:gameId/start` | 開始遊戲 | 需要 X-Account-ID (房主) |

## 架構設計

### Clean Architecture 分層

```
┌─────────────────────────────────────────────────────────────┐
│                    Adapter Layer                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ GameRoomHandler                                      │   │
│  │ - AuthMiddleware (X-Account-ID 認證)                │   │
│  │ - CreateGameRoom / JoinGameRoom / StartGame         │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Use Case Layer                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ GameRoomUseCase                                      │   │
│  │ - CreateGameRoom: 建立遊戲房並自動加入房主          │   │
│  │ - JoinGameRoom: 驗證並加入遊戲房                    │   │
│  │ - StartGame: 驗證房主並開始遊戲                     │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer                              │
│  ┌──────────────────┐  ┌──────────────────────────────┐    │
│  │ Game Entity      │  │ GamePlayer Entity            │    │
│  │ - HostAccountID  │  │ - GameID                     │    │
│  │ - MaxPlayers     │  │ - AccountID                  │    │
│  │ - CurrentPlayers │  │ - JoinOrder                  │    │
│  │ - CanJoin()      │  └──────────────────────────────┘    │
│  │ - CanStart()     │                                      │
│  └──────────────────┘  ┌──────────────────────────────┐    │
│                        │ GamePlayerRepository         │    │
│                        │ (Interface)                  │    │
│                        └──────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                Infrastructure Layer                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ GamePlayerRepository (MySQL Implementation)         │   │
│  │ - CreateGamePlayer / GetGamePlayersByGameID        │   │
│  │ - ExistsByGameIDAndAccountID / CountByGameID       │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Migrations                                          │   │
│  │ - 000009_add_game_room_fields (games 表新增欄位)   │   │
│  │ - game_players 表 (多對多關聯)                     │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 實作細節

### 1. Domain Layer

#### Game Entity 擴展 (`internal/domain/entity/game.go`)
```go
// 新增欄位
type Game struct {
    // ... 原有欄位
    HostAccountID   int          // 房主帳號 ID
    MaxPlayers      int          // 最大玩家數
    CurrentPlayers  int          // 目前玩家數
    GamePlayers     []GamePlayer // 遊戲玩家列表
}

// 遊戲房狀態常數
const (
    GameRoomStatusWaiting = "WAITING"  // 等待中
    GameRoomStatusPlaying = "PLAYING"  // 遊戲中
    GameRoomStatusEnded   = "ENDED"    // 已結束
)

// 玩家數限制
const (
    GameMinPlayers = 3  // 最少 3 人
    GameMaxPlayers = 9  // 最多 9 人
)

// 業務邏輯方法
func NewGameRoom(hostAccountID int, maxPlayers int) *Game
func (g *Game) CanJoin() error
func (g *Game) CanStart(accountID int) error
```

#### GamePlayer Entity (`internal/domain/entity/game_player.go`)
```go
type GamePlayer struct {
    ID        int
    GameID    int
    AccountID int
    JoinOrder int       // 加入順序
    CreatedAt time.Time
}
```

#### GamePlayerRepository Interface (`internal/domain/repository/game_player_repository.go`)
```go
type GamePlayerRepository interface {
    CreateGamePlayer(ctx context.Context, gamePlayer *entity.GamePlayer) (*entity.GamePlayer, error)
    GetGamePlayersByGameID(ctx context.Context, gameID int) ([]*entity.GamePlayer, error)
    ExistsByGameIDAndAccountID(ctx context.Context, gameID int, accountID int) (bool, error)
    CountByGameID(ctx context.Context, gameID int) (int, error)
}
```

### 2. Infrastructure Layer

#### Database Migration (`database/migrations/000009_add_game_room_fields.up.sql`)
```sql
-- 為 games 表新增遊戲房相關欄位
ALTER TABLE games ADD COLUMN host_account_id INT NOT NULL DEFAULT 0;
ALTER TABLE games ADD COLUMN max_players INT NOT NULL DEFAULT 9;
ALTER TABLE games ADD COLUMN current_players INT NOT NULL DEFAULT 0;

-- 建立 game_players 表（記錄遊戲與帳號的多對多關係）
CREATE TABLE IF NOT EXISTS game_players (
    id INT AUTO_INCREMENT PRIMARY KEY,
    game_id INT NOT NULL,
    account_id INT NOT NULL,
    join_order INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    UNIQUE KEY unique_game_account (game_id, account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### GamePlayerRepository Implementation (`internal/infrastructure/persistence/mysql/game_player_repository.go`)
- 使用 GORM 實作 Repository 介面
- GamePlayerModel 與 GamePlayer Entity 互相轉換

### 3. Use Case Layer

#### GameRoomUseCase (`internal/usecase/game_room_usecase.go`)
```go
type GameRoomUseCase interface {
    CreateGameRoom(ctx context.Context, accountID int, maxPlayers int) (*entity.Game, error)
    JoinGameRoom(ctx context.Context, gameID int, accountID int) error
    StartGame(ctx context.Context, gameID int, accountID int) (*entity.Game, error)
}
```

**業務邏輯：**
- `CreateGameRoom`: 驗證帳號存在 → 建立遊戲 → 房主自動加入
- `JoinGameRoom`: 驗證遊戲可加入 → 檢查是否重複加入 → 新增 GamePlayer
- `StartGame`: 驗證房主權限 → 檢查人數 ≥ 3 → 更新遊戲狀態

### 4. Adapter Layer

#### GameRoomHandler (`internal/adapter/http/handler/game_room_handler.go`)
```go
// AuthMiddleware - 從 Header 取得 X-Account-ID
func AuthMiddleware() gin.HandlerFunc

// Handler 方法
func (h *GameRoomHandler) CreateGameRoom(c *gin.Context)
func (h *GameRoomHandler) JoinGameRoom(c *gin.Context)
func (h *GameRoomHandler) StartGame(c *gin.Context)
```

#### Request/Response (`internal/adapter/http/request/request.go`, `response/response.go`)
```go
// Request
type CreateGameRoomRequest struct {
    MaxPlayers int `json:"maxPlayers"`
}

// Response
type CreateGameRoomResponse struct {
    GameID       int    `json:"gameId"`
    Status       string `json:"status"`
    HostPlayerID int    `json:"hostPlayerId"`
    MaxPlayers   int    `json:"maxPlayers"`
}
```

## E2E 測試案例

### 測試檔案: `tests/e2e/game_room_clean_test.go`

| 測試名稱 | 說明 | 預期結果 |
|----------|------|----------|
| TestCreateGameRoom_Success | 建立遊戲房成功 | 200 OK |
| TestCreateGameRoom_Unauthorized | 未認證建立遊戲房 | 401 Unauthorized |
| TestJoinGameRoom_Success | 加入遊戲房成功 | 200 OK |
| TestJoinGameRoom_AlreadyJoined | 重複加入遊戲房 | 400 Bad Request |
| TestJoinGameRoom_GameFull | 遊戲房已滿 | 400 Bad Request |
| TestJoinGameRoom_GameAlreadyStarted | 遊戲已開始 | 400 Bad Request |
| TestStartGame_AsHost_Success | 房主開始遊戲 | 200 OK |
| TestStartGame_NotHost_Fail | 非房主開始遊戲 | 400 Bad Request |
| TestStartGame_NotEnoughPlayers | 人數不足開始遊戲 | 400 Bad Request |
| TestStartGame_ExactlyThreePlayers | 剛好 3 人開始遊戲 | 200 OK |

## 錯誤代碼

| 錯誤 | 狀態碼 | 說明 |
|------|--------|------|
| ErrUnauthorized | 401 | 未提供 X-Account-ID |
| ErrAccountNotFound | 400 | 帳號不存在 |
| ErrGameNotFound | 400 | 遊戲不存在 |
| ErrGameNotWaiting | 400 | 遊戲非等待狀態 |
| ErrGameFull | 400 | 遊戲房已滿 |
| ErrAlreadyJoined | 400 | 已加入該遊戲房 |
| ErrNotHost | 400 | 非房主無法開始遊戲 |
| ErrNotEnoughPlayers | 400 | 玩家人數不足 (< 3) |

## 相關檔案清單

### 新增檔案
- `internal/domain/entity/game_player.go`
- `internal/domain/repository/game_player_repository.go`
- `internal/infrastructure/persistence/mysql/game_player_repository.go`
- `internal/infrastructure/persistence/mysql/model/game_player_model.go`
- `internal/usecase/game_room_usecase.go`
- `internal/adapter/http/handler/game_room_handler.go`
- `database/migrations/000009_add_game_room_fields.up.sql`
- `database/migrations/000009_add_game_room_fields.down.sql`
- `tests/e2e/game_room_clean_test.go`

### 修改檔案
- `internal/domain/entity/game.go` - 新增遊戲房欄位與方法
- `internal/infrastructure/persistence/mysql/model/game_model.go` - 新增欄位
- `internal/adapter/http/request/request.go` - 新增 Request 結構
- `internal/adapter/http/response/response.go` - 新增 Response 結構
- `internal/adapter/http/handler/game_handler.go` - 路由改為 `/api/v1/games/start-legacy`
- `cmd/app/main.go` - 註冊新的 Handler 與 Repository
- `tests/e2e/suite_clean_test.go` - 新增測試依賴

## 與 Task 1.1 的關聯

Task 1.2 建立在 Task 1.1 的基礎上：
- 使用 Task 1.1 建立的 `Account` 實體進行玩家認證
- 使用 `AccountRepository` 驗證帳號存在
- 遊戲房的 `HostAccountID` 關聯到 `accounts` 表

## 測試執行結果

```bash
go test -v ./tests/e2e/... -run "TestCleanArchTestSuite"
```

所有 10 個遊戲房相關測試案例通過：
- TestCreateGameRoom_Success ✓
- TestCreateGameRoom_Unauthorized ✓
- TestJoinGameRoom_Success ✓
- TestJoinGameRoom_AlreadyJoined ✓
- TestJoinGameRoom_GameFull ✓
- TestJoinGameRoom_GameAlreadyStarted ✓
- TestStartGame_AsHost_Success ✓
- TestStartGame_NotHost_Fail ✓
- TestStartGame_NotEnoughPlayers ✓
- TestStartGame_ExactlyThreePlayers ✓
