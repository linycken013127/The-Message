# The Message - 後端 Clean Architecture 技術文件

## 目錄

1. [專案概述](#專案概述)
2. [Clean Architecture 架構說明](#clean-architecture-架構說明)
3. [目錄結構](#目錄結構)
4. [各層詳細說明](#各層詳細說明)
5. [API 文件](#api-文件)
6. [開發指南](#開發指南)
7. [測試指南](#測試指南)
8. [部署說明](#部署說明)

---

## 專案概述

「風聲」（The Message）是一個桌遊線上版的後端 API 服務，使用 Go 語言開發，採用 Clean Architecture（乾淨架構）設計模式。

### 技術棧

| 技術 | 版本 | 用途 |
|------|------|------|
| Go | 1.22+ | 主要開發語言 |
| Gin | v1.9.1 | HTTP Web 框架 |
| GORM | v1.25.7 | ORM 資料庫操作 |
| MySQL | 8.0+ | 資料庫 |
| golang-migrate | v4 | 資料庫遷移 |
| Swagger | swaggo | API 文件生成 |
| testify | v1.8.4 | 測試框架 |

---

## Clean Architecture 架構說明

本專案採用 Clean Architecture，遵循依賴反轉原則（Dependency Inversion Principle），確保業務邏輯與外部框架解耦。

### 架構圖

```
┌─────────────────────────────────────────────────────────────┐
│                    External Systems                          │
│  (HTTP Client, Database, Third-party Services)              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Adapter Layer (適配器層)                    │
│  ┌─────────────────┐  ┌─────────────────────────────────┐  │
│  │   HTTP Handler  │  │    SSE Server                   │  │
│  │   (Gin Router)  │  │    (Server-Sent Events)         │  │
│  └─────────────────┘  └─────────────────────────────────┘  │
│  internal/adapter/http/handler + internal/adapter/sse      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Use Case Layer (用例層)                    │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────┐│
│  │ GameUseCase  │ │PlayerUseCase │ │CardUseCase/DeckUseCase│
│  └──────────────┘ └──────────────┘ └──────────────────────┘│
│                  internal/usecase                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer (領域層)                     │
│  ┌──────────────────┐  ┌─────────────────────────────────┐ │
│  │     Entities     │  │     Repository Interfaces       │ │
│  │  (Game, Player,  │  │  (GameRepository, PlayerRepo,   │ │
│  │   Card, Deck,    │  │   CardRepository, DeckRepository│ │
│  │   PlayerCard,    │  │   PlayerCardRepository,         │ │
│  │   GameProgress)  │  │   GameProgressRepository)       │ │
│  └──────────────────┘  └─────────────────────────────────┘ │
│         internal/domain/entity + internal/domain/repository │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│               Infrastructure Layer (基礎設施層)              │
│  ┌─────────────────┐  ┌─────────────────────────────────┐  │
│  │  MySQL Repos    │  │  Config & Seeders               │  │
│  │  (GORM Models)  │  │  (Database, Migration)          │  │
│  └─────────────────┘  └─────────────────────────────────┘  │
│            internal/infrastructure                          │
└─────────────────────────────────────────────────────────────┘
```

### 依賴規則

- **Domain Layer** 是核心，不依賴任何外部層
- **Use Case Layer** 只依賴 Domain Layer
- **Adapter Layer** 依賴 Use Case Layer 和 Domain Layer
- **Infrastructure Layer** 實作 Domain Layer 定義的介面

---

## 目錄結構

```
Backend/
├── cmd/
│   └── app/
│       └── main.go              # 應用程式入口點
├── internal/
│   ├── domain/                  # 領域層（核心業務邏輯）
│   │   ├── entity/              # 領域實體
│   │   │   ├── game.go          # 遊戲實體
│   │   │   ├── player.go        # 玩家實體
│   │   │   ├── card.go          # 卡片實體
│   │   │   ├── deck.go          # 牌組實體
│   │   │   ├── player_card.go   # 玩家手牌實體
│   │   │   └── game_progress.go # 遊戲進度實體
│   │   └── repository/          # 倉儲介面
│   │       ├── game_repository.go
│   │       ├── player_repository.go
│   │       ├── card_repository.go
│   │       ├── deck_repository.go
│   │       ├── player_card_repository.go
│   │       └── game_progress_repository.go
│   ├── usecase/                 # 用例層（應用業務邏輯）
│   │   ├── game_usecase.go      # 遊戲用例
│   │   ├── player_usecase.go    # 玩家用例
│   │   ├── card_usecase.go      # 卡片用例
│   │   ├── deck_usecase.go      # 牌組用例
│   │   └── dto/                 # 資料傳輸物件
│   │       └── dto.go
│   ├── adapter/                 # 適配器層（外部介面）
│   │   ├── http/
│   │   │   ├── handler/         # HTTP 處理器
│   │   │   │   ├── game_handler.go
│   │   │   │   ├── player_handler.go
│   │   │   │   ├── card_handler.go
│   │   │   │   └── heartbeat_handler.go
│   │   │   ├── request/         # HTTP 請求結構
│   │   │   │   └── request.go
│   │   │   └── response/        # HTTP 回應結構
│   │   │       └── response.go
│   │   └── sse/                 # Server-Sent Events
│   │       └── sse_handler.go
│   └── infrastructure/          # 基礎設施層
│       ├── config/              # 設定
│       │   ├── database.go
│       │   ├── migration.go
│       │   └── test_database.go
│       ├── persistence/
│       │   └── mysql/           # MySQL 實作
│       │       ├── model/       # GORM 模型
│       │       │   ├── game_model.go
│       │       │   ├── player_model.go
│       │       │   ├── card_model.go
│       │       │   ├── deck_model.go
│       │       │   ├── player_card_model.go
│       │       │   └── game_progress_model.go
│       │       ├── game_repository.go
│       │       ├── player_repository.go
│       │       ├── card_repository.go
│       │       ├── deck_repository.go
│       │       ├── player_card_repository.go
│       │       └── game_progress_repository.go
│       └── seeders/             # 資料種子
│           └── database_seeder.go
├── pkg/                         # 共用套件
│   └── errors/
│       └── errors.go
├── database/
│   └── migrations/              # 資料庫遷移檔
├── tests/
│   └── e2e/                     # 端對端測試
└── docs/                        # 文件
```

---

## 各層詳細說明

### Domain Layer（領域層）

領域層是系統的核心，包含純業務邏輯，不依賴任何外部框架。

#### 實體（Entities）

| 實體 | 說明 | 主要屬性 |
|------|------|----------|
| Game | 遊戲 | ID, Token, Status, CurrentPlayerID, Players |
| Player | 玩家 | ID, Name, GameID, IdentityCard, Status, OrderNumber |
| Card | 卡片 | ID, Name, Color, IntelligenceType |
| Deck | 牌組 | ID, GameID, CardID |
| PlayerCard | 玩家手牌 | ID, PlayerID, GameID, CardID, Type |
| GameProgress | 遊戲進度 | ID, PlayerID, GameID, CardID, Action, TargetPlayerID |

#### 常數定義

**遊戲狀態：**
```go
GameStatusStart                = "開始遊戲"
GameStatusActionCardStage      = "功能牌階段"
GameStatusTransmitIntelligence = "情報牌階段"
GameStatusEnd                  = "結束遊戲"
```

**玩家狀態：**
```go
PlayerStatusAlive = "生存"
PlayerStatusDead  = "死亡"
```

**身份卡：**
```go
IdentityUndercoverFront = "潛伏戰線"
IdentityMilitaryAgency  = "軍情處"
IdentityBystander       = "打醬油"
```

**卡片顏色：**
```go
CardColorRed   = "紅"
CardColorBlue  = "藍"
CardColorBlack = "黑"
```

#### 倉儲介面

```go
// GameRepository 遊戲倉儲介面
type GameRepository interface {
    CreateGame(ctx context.Context, game *entity.Game) (*entity.Game, error)
    GetGameById(ctx context.Context, id int) (*entity.Game, error)
    GetGameWithPlayers(ctx context.Context, id int) (*entity.Game, error)
    UpdateGame(ctx context.Context, game *entity.Game) error
    DeleteGame(ctx context.Context, id int) error
}
```

### Use Case Layer（用例層）

用例層封裝應用程式的業務邏輯，協調領域物件完成特定任務。

#### GameUseCase

| 方法 | 說明 |
|------|------|
| InitGame | 初始化遊戲（產生 Token、建立遊戲記錄） |
| InitDeck | 初始化牌組 |
| DrawCard | 抽牌 |
| DrawCardsForAllPlayers | 為所有玩家抽牌 |
| GetGameById | 根據 ID 取得遊戲 |
| CreateGame | 建立遊戲 |
| DeleteGame | 刪除遊戲 |
| UpdateCurrentPlayer | 更新當前玩家 |
| NextPlayer | 切換到下一位玩家 |
| UpdateStatus | 更新遊戲狀態 |

#### PlayerUseCase

| 方法 | 說明 |
|------|------|
| InitPlayers | 初始化玩家（分配身份卡） |
| CanPlayCard | 檢查玩家是否可以出牌 |
| PlayCard | 出牌 |
| TransmitIntelligenceCard | 傳情報卡 |
| AcceptCard | 接收卡片 |
| CheckWin | 檢查勝利條件 |

### Adapter Layer（適配器層）

適配器層負責將外部請求轉換為用例層可理解的格式。

#### HTTP Handler

| Handler | 路由 | 方法 | 說明 |
|---------|------|------|------|
| GameHandler | /api/v1/games | POST | 開始新遊戲 |
| GameHandler | /api/v1/games/:gameId/events | GET | SSE 遊戲事件 |
| PlayerHandler | /api/v1/players/:playerId/player-cards | POST | 出牌 |
| PlayerHandler | /api/v1/player/:playerId/transmit-intelligence | POST | 傳情報 |
| PlayerHandler | /api/v1/players/:playerId/accept | POST | 接收/拒絕卡片 |
| CardHandler | /api/v1/player/:playerId/player-cards/ | GET | 取得玩家手牌 |
| HeartbeatHandler | /api/v1/heartbeat | GET | 健康檢查 |

### Infrastructure Layer（基礎設施層）

基礎設施層實作領域層定義的介面，處理與外部系統的互動。

#### MySQL Repository

每個 Repository 實作對應的領域介面，使用 GORM 進行資料庫操作。

#### Model 轉換

每個 GORM Model 都提供：
- `ToEntity()` - 將 Model 轉換為領域實體
- `FromEntity()` - 將領域實體轉換為 Model

---

## API 文件

### 1. 開始新遊戲

**POST** `/api/v1/games`

**請求：**
```json
{
  "players": [
    {"id": "player1_id", "name": "玩家A"},
    {"id": "player2_id", "name": "玩家B"},
    {"id": "player3_id", "name": "玩家C"}
  ]
}
```

**回應：**
```json
{
  "Id": 1,
  "Token": "abc123..."
}
```

---

### 2. 取得遊戲事件（SSE）

**GET** `/api/v1/games/:gameId/events`

**回應：** Server-Sent Events 串流

```
event: message
data: {"game_id": 1, "status": "功能牌階段", "current_player": 1}
```

---

### 3. 出牌

**POST** `/api/v1/players/:playerId/player-cards`

**請求：**
```json
{
  "card_id": 5
}
```

**回應：**
```json
{
  "result": true
}
```

---

### 4. 傳情報

**POST** `/api/v1/player/:playerId/transmit-intelligence`

**請求：**
```json
{
  "card_id": 5
}
```

**回應：**
```json
{
  "result": true
}
```

---

### 5. 接收/拒絕卡片

**POST** `/api/v1/players/:playerId/accept`

**請求：**
```json
{
  "accept": true
}
```

**回應：**
```json
{
  "result": true
}
```

---

### 6. 取得玩家手牌

**GET** `/api/v1/player/:playerId/player-cards/`

**回應：**
```json
{
  "player_cards": [
    {"id": 1, "name": "鎖定", "color": "紅"},
    {"id": 2, "name": "試探", "color": "藍"},
    {"id": 3, "name": "截獲", "color": "黑"}
  ]
}
```

---

### 7. 健康檢查

**GET** `/api/v1/heartbeat`

**回應：** HTTP 204 No Content

---

## 開發指南

### 環境設定

1. **安裝 Go 1.22+**

2. **設定環境變數**（`.env` 檔案）：
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_DATABASE=the_message
```

3. **安裝相依套件**：
```bash
go mod tidy
```

4. **執行資料庫遷移**：
```bash
go run cmd/migrate/migrate.go
```

5. **種子資料**：
```bash
go run cmd/migrate/game_card_seeder.go
```

6. **啟動伺服器**：
```bash
go run cmd/app/main.go
```

### 新增功能的步驟

1. **Domain Layer**：
   - 在 `entity/` 新增或修改實體
   - 在 `repository/` 定義倉儲介面

2. **Use Case Layer**：
   - 在 `usecase/` 新增用例
   - 定義介面和實作

3. **Infrastructure Layer**：
   - 在 `mysql/model/` 新增 GORM Model
   - 在 `mysql/` 實作倉儲介面

4. **Adapter Layer**：
   - 在 `handler/` 新增 HTTP Handler
   - 在 `request/` 和 `response/` 定義 DTO

5. **更新 main.go**：
   - 註冊新的依賴

---

## 測試指南

### 執行測試

```bash
# 執行所有測試
go test ./...

# 執行 e2e 測試
go test ./tests/e2e/...

# 執行特定測試
go test -run TestStartGameE2E ./tests/e2e/...

# 顯示詳細輸出
go test -v ./tests/e2e/...
```

### 測試結構

- `tests/e2e/suite_clean_test.go` - Clean Architecture 測試套件
- `tests/e2e/game_clean_test.go` - 遊戲相關測試
- `tests/e2e/player_clean_test.go` - 玩家相關測試
- `tests/e2e/player_card_clean_test.go` - 玩家手牌測試

---

## 部署說明

### Docker 部署

```bash
# 建置映像
docker build -t the-message .

# 執行容器
docker run -p 8080:8080 --env-file .env the-message
```

### Docker Compose

```bash
docker-compose up -d
```

---

## 附錄

### 遊戲規則簡述

1. 遊戲開始時，每位玩家會被分配一個秘密身份
2. 玩家輪流出牌，傳遞情報
3. 當玩家收集到足夠的特定顏色情報時獲勝：
   - 軍情處：收集 3 張紅色情報
   - 潛伏戰線：收集 3 張藍色情報
   - 任一方收集 5 張任意情報

### 聯繫資訊

如有問題，請聯繫專案維護人員。
