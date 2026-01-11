# Clean Architecture 重構 SOP

> **目標受眾**：無乾淨架構經驗的開發者
> **預期成果**：將現有後端程式完整重構為符合 Clean Architecture 原則的架構
> **文件版本**：v1.0

---

## 目錄

1. [乾淨架構核心概念](#1-乾淨架構核心概念)
2. [目前架構分析](#2-目前架構分析)
3. [目標架構設計](#3-目標架構設計)
4. [重構階段規劃](#4-重構階段規劃)
5. [詳細執行步驟](#5-詳細執行步驟)
6. [程式碼範例與模板](#6-程式碼範例與模板)
7. [檢查清單](#7-檢查清單)
8. [常見問題與解決方案](#8-常見問題與解決方案)

---

## 1. 乾淨架構核心概念

### 1.1 什麼是乾淨架構？

乾淨架構（Clean Architecture）是 Robert C. Martin（Uncle Bob）提出的軟體架構原則，核心理念是**依賴規則（Dependency Rule）**：

> **外層可以依賴內層，但內層絕不能依賴外層**

```
┌─────────────────────────────────────────────────────────┐
│                    外層 (Frameworks)                     │
│  ┌───────────────────────────────────────────────────┐  │
│  │              Interface Adapters                    │  │
│  │  ┌─────────────────────────────────────────────┐  │  │
│  │  │           Application Business Rules         │  │  │
│  │  │  ┌───────────────────────────────────────┐  │  │  │
│  │  │  │     Enterprise Business Rules         │  │  │  │
│  │  │  │           (Entities)                  │  │  │  │
│  │  │  └───────────────────────────────────────┘  │  │  │
│  │  │              (Use Cases)                    │  │  │
│  │  └─────────────────────────────────────────────┘  │  │
│  │        (Controllers, Presenters, Gateways)        │  │
│  └───────────────────────────────────────────────────┘  │
│              (Web, DB, Devices, External APIs)          │
└─────────────────────────────────────────────────────────┘
```

### 1.2 四層架構說明

| 層級 | 名稱 | 職責 | Go 對應目錄 |
|------|------|------|-------------|
| 第 1 層 | **Entities** | 企業業務規則、核心領域模型 | `domain/` |
| 第 2 層 | **Use Cases** | 應用業務規則、操作流程編排 | `usecase/` |
| 第 3 層 | **Interface Adapters** | 轉換資料格式、橋接內外層 | `adapter/` |
| 第 4 層 | **Frameworks & Drivers** | 框架、資料庫、Web 伺服器 | `infrastructure/` |

### 1.3 核心原則

#### 原則 1：依賴反轉（Dependency Inversion）
- 高層模組不應依賴低層模組，兩者都應依賴抽象
- 抽象不應依賴細節，細節應依賴抽象

```go
// ❌ 錯誤：Service 直接依賴具體實現
type GameService struct {
    repo *mysql.GameRepository  // 直接依賴 MySQL 實現
}

// ✅ 正確：Service 依賴介面
type GameService struct {
    repo repository.GameRepository  // 依賴介面
}
```

#### 原則 2：單一職責（Single Responsibility）
- 每個模組只有一個改變的理由
- 一個類別/結構體只做一件事

#### 原則 3：開放封閉（Open/Closed）
- 對擴展開放，對修改封閉
- 新增功能時不需修改現有程式碼

---

## 2. 目前架構分析

### 2.1 現有目錄結構

```
Backend/
├── cmd/                    # ✅ 正確：應用程式入口
│   ├── app/main.go
│   └── migrate/
├── config/                 # ⚠️ 需調整：應移至 infrastructure
├── database/               # ⚠️ 需調整：應移至 infrastructure
├── enums/                  # ⚠️ 需調整：應拆分至 domain
├── service/
│   ├── delivery/http/v1/   # ✅ 正確：HTTP 處理層
│   ├── service/            # ⚠️ 需調整：應重新命名為 usecase
│   ├── repository/         # ⚠️ 需調整：介面與實現應分離
│   └── request/            # ⚠️ 需調整：應移至 adapter/http
├── tests/                  # ✅ 正確：測試目錄
└── utils/                  # ⚠️ 需調整：應移至 pkg 或各層
```

### 2.2 目前問題

| 問題 | 嚴重程度 | 說明 |
|------|----------|------|
| 領域模型與 GORM 耦合 | 🔴 高 | Entity 直接使用 `gorm.Model`，違反依賴規則 |
| Service 職責過重 | 🟠 中 | `PlayerService` 306 行，混雜多種業務邏輯 |
| 缺乏 Use Case 層 | 🟠 中 | Service 同時處理業務邏輯與編排 |
| 錯誤處理不一致 | 🟠 中 | 部分使用 `panic()`，部分返回 error |
| 缺乏事務管理 | 🟠 中 | 複雜操作缺少原子性保證 |
| 列舉散落各處 | 🟡 低 | `enums/` 應整合進 domain 層 |

### 2.3 優點（保留）

- ✅ 已有基本的分層概念
- ✅ Repository 使用介面定義
- ✅ 依賴注入模式
- ✅ 完整的 E2E 測試框架
- ✅ API 版本化設計

---

## 3. 目標架構設計

### 3.1 目標目錄結構

```
Backend/
├── cmd/                           # 應用程式入口
│   ├── app/
│   │   └── main.go               # 主程式，負責 DI 組裝
│   └── migrate/
│       └── main.go               # 資料庫遷移工具
│
├── internal/                      # 內部程式碼（不對外暴露）
│   │
│   ├── domain/                    # 第 1 層：企業業務規則
│   │   ├── entity/               # 領域實體（純 Go 結構體）
│   │   │   ├── game.go
│   │   │   ├── player.go
│   │   │   ├── card.go
│   │   │   ├── deck.go
│   │   │   └── player_card.go
│   │   ├── valueobject/          # 值物件
│   │   │   ├── game_status.go
│   │   │   ├── player_status.go
│   │   │   ├── card_color.go
│   │   │   └── identity_card.go
│   │   ├── repository/           # Repository 介面定義
│   │   │   ├── game_repository.go
│   │   │   ├── player_repository.go
│   │   │   ├── card_repository.go
│   │   │   └── deck_repository.go
│   │   └── service/              # Domain Service 介面
│   │       └── identity_card_service.go
│   │
│   ├── usecase/                   # 第 2 層：應用業務規則
│   │   ├── game/
│   │   │   ├── create_game.go    # 建立遊戲 Use Case
│   │   │   ├── start_game.go     # 開始遊戲 Use Case
│   │   │   └── get_game.go       # 取得遊戲 Use Case
│   │   ├── player/
│   │   │   ├── play_card.go      # 出牌 Use Case
│   │   │   ├── transmit_intelligence.go
│   │   │   └── accept_card.go
│   │   └── dto/                  # 資料傳輸物件
│   │       ├── game_dto.go
│   │       └── player_dto.go
│   │
│   ├── adapter/                   # 第 3 層：介面轉接器
│   │   ├── http/                 # HTTP 轉接器
│   │   │   ├── handler/
│   │   │   │   ├── game_handler.go
│   │   │   │   ├── player_handler.go
│   │   │   │   └── health_handler.go
│   │   │   ├── middleware/
│   │   │   │   ├── cors.go
│   │   │   │   ├── recovery.go
│   │   │   │   └── logger.go
│   │   │   ├── request/
│   │   │   │   ├── game_request.go
│   │   │   │   └── player_request.go
│   │   │   ├── response/
│   │   │   │   ├── game_response.go
│   │   │   │   └── error_response.go
│   │   │   └── router.go
│   │   ├── presenter/            # 輸出格式轉換
│   │   │   ├── game_presenter.go
│   │   │   └── player_presenter.go
│   │   └── sse/                  # SSE 轉接器
│   │       └── event_handler.go
│   │
│   └── infrastructure/            # 第 4 層：框架與驅動
│       ├── persistence/          # 資料持久化
│       │   ├── mysql/
│       │   │   ├── connection.go
│       │   │   ├── game_repository.go
│       │   │   ├── player_repository.go
│       │   │   ├── card_repository.go
│       │   │   └── model/        # GORM 模型（與 Entity 分離）
│       │   │       ├── game_model.go
│       │   │       ├── player_model.go
│       │   │       └── mapper.go # Entity ↔ Model 轉換
│       │   └── migration/
│       │       └── migrations/
│       ├── config/               # 配置管理
│       │   ├── config.go
│       │   └── database.go
│       └── logger/               # 日誌系統
│           └── logger.go
│
├── pkg/                           # 可重用的公共套件
│   ├── errors/                   # 自定義錯誤
│   │   └── errors.go
│   └── utils/                    # 工具函式
│       └── random.go
│
├── api/                           # API 文件
│   └── swagger/
│
├── tests/                         # 測試
│   ├── e2e/
│   ├── integration/
│   └── unit/
│
├── docker-compose.yml
├── Makefile
└── go.mod
```

### 3.2 層級依賴圖

```
cmd/app/main.go (DI Container)
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│                    infrastructure/                           │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                     adapter/                         │    │
│  │  ┌─────────────────────────────────────────────┐    │    │
│  │  │                  usecase/                    │    │    │
│  │  │  ┌───────────────────────────────────────┐  │    │    │
│  │  │  │               domain/                  │  │    │    │
│  │  │  │   (Entity, ValueObject, Repository)   │  │    │    │
│  │  │  └───────────────────────────────────────┘  │    │    │
│  │  └─────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘

依賴方向：外 → 內（只能由外層依賴內層）
```

---

## 4. 重構階段規劃

### 階段概覽

```
階段 1: 建立 Domain 層
    │
    ▼
階段 2: 建立 Use Case 層
    │
    ▼
階段 3: 重構 Adapter 層
    │
    ▼
階段 4: 重構 Infrastructure 層
    │
    ▼
階段 5: 整合與測試
    │
    ▼
階段 6: 清理與文件
```

### 階段 1：建立 Domain 層

**目標**：建立純淨的領域模型，不依賴任何外部框架

| 任務 | 說明 | 優先級 |
|------|------|--------|
| 1.1 建立目錄結構 | 建立 `internal/domain/` | P0 |
| 1.2 建立 Entity | 將現有模型轉為純 Go 結構體 | P0 |
| 1.3 建立 Value Object | 提取列舉為值物件 | P1 |
| 1.4 定義 Repository 介面 | 將介面移至 domain 層 | P0 |
| 1.5 定義 Domain Service 介面 | 跨 Entity 的業務邏輯 | P1 |

### 階段 2：建立 Use Case 層

**目標**：將業務邏輯封裝為獨立的 Use Case

| 任務 | 說明 | 優先級 |
|------|------|--------|
| 2.1 識別 Use Case | 分析現有 Service 方法 | P0 |
| 2.2 建立 DTO | 定義輸入/輸出資料結構 | P0 |
| 2.3 實作 Game Use Cases | CreateGame, StartGame, GetGame | P0 |
| 2.4 實作 Player Use Cases | PlayCard, TransmitIntelligence | P0 |
| 2.5 建立 Use Case 介面 | 方便測試與擴展 | P1 |

### 階段 3：重構 Adapter 層

**目標**：將 HTTP 處理邏輯與業務邏輯解耦

| 任務 | 說明 | 優先級 |
|------|------|--------|
| 3.1 建立 Request/Response 結構 | 分離 HTTP 資料結構 | P0 |
| 3.2 重構 Handler | 只負責 HTTP 邏輯 | P0 |
| 3.3 建立 Presenter | 格式化輸出 | P1 |
| 3.4 統一錯誤處理 | 建立錯誤轉換機制 | P0 |
| 3.5 重構 Middleware | 提取共用中間件 | P1 |

### 階段 4：重構 Infrastructure 層

**目標**：將框架相關程式碼隔離

| 任務 | 說明 | 優先級 |
|------|------|--------|
| 4.1 分離 GORM Model | 建立獨立的 DB 模型 | P0 |
| 4.2 實作 Mapper | Entity ↔ Model 轉換 | P0 |
| 4.3 實作 Repository | 實現 domain 層介面 | P0 |
| 4.4 重構 Config | 集中配置管理 | P1 |
| 4.5 建立 Logger | 結構化日誌 | P2 |

### 階段 5：整合與測試

**目標**：確保重構後功能正確

| 任務 | 說明 | 優先級 |
|------|------|--------|
| 5.1 更新 main.go | 調整 DI 組裝 | P0 |
| 5.2 執行現有測試 | 確保不破壞現有功能 | P0 |
| 5.3 補充單元測試 | 針對 Use Case 層 | P1 |
| 5.4 補充整合測試 | 針對 Repository | P1 |

### 階段 6：清理與文件

**目標**：移除舊程式碼，更新文件

| 任務 | 說明 | 優先級 |
|------|------|--------|
| 6.1 移除舊目錄 | 刪除 service/, enums/ | P0 |
| 6.2 更新 CLAUDE.md | 更新架構說明 | P1 |
| 6.3 更新 Swagger | 重新生成 API 文件 | P1 |

---

## 5. 詳細執行步驟

### 步驟 1.1：建立目錄結構

```bash
cd Backend

# 建立 domain 層目錄
mkdir -p internal/domain/{entity,valueobject,repository,service}

# 建立 usecase 層目錄
mkdir -p internal/usecase/{game,player,card,dto}

# 建立 adapter 層目錄
mkdir -p internal/adapter/http/{handler,middleware,request,response}
mkdir -p internal/adapter/{presenter,sse}

# 建立 infrastructure 層目錄
mkdir -p internal/infrastructure/persistence/mysql/model
mkdir -p internal/infrastructure/{config,logger}
mkdir -p internal/infrastructure/persistence/migration

# 建立 pkg 目錄
mkdir -p pkg/{errors,utils}
```

### 步驟 1.2：建立 Entity（範例：Game）

**檔案**：`internal/domain/entity/game.go`

```go
package entity

import "time"

// Game 遊戲領域實體
// 注意：不依賴任何外部框架（無 gorm.Model）
type Game struct {
    ID              int
    Token           string
    Status          GameStatus
    CurrentPlayerID int
    Players         []Player
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// GameStatus 遊戲狀態值物件
type GameStatus string

const (
    GameStatusWaiting  GameStatus = "waiting"
    GameStatusPlaying  GameStatus = "playing"
    GameStatusFinished GameStatus = "finished"
)

// NewGame 建立新遊戲（工廠方法）
func NewGame(token string, players []string) (*Game, error) {
    if len(players) < 3 || len(players) > 9 {
        return nil, ErrInvalidPlayerCount
    }

    game := &Game{
        Token:     token,
        Status:    GameStatusWaiting,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    for i, name := range players {
        player := NewPlayer(name, i+1)
        game.Players = append(game.Players, *player)
    }

    return game, nil
}

// Start 開始遊戲
func (g *Game) Start() error {
    if g.Status != GameStatusWaiting {
        return ErrGameAlreadyStarted
    }
    g.Status = GameStatusPlaying
    g.UpdatedAt = time.Now()
    return nil
}

// IsPlaying 是否進行中
func (g *Game) IsPlaying() bool {
    return g.Status == GameStatusPlaying
}

// NextPlayer 切換到下一位玩家
func (g *Game) NextPlayer() error {
    if !g.IsPlaying() {
        return ErrGameNotPlaying
    }

    currentIndex := -1
    for i, p := range g.Players {
        if p.ID == g.CurrentPlayerID {
            currentIndex = i
            break
        }
    }

    if currentIndex == -1 {
        return ErrPlayerNotFound
    }

    // 找下一個存活的玩家
    for i := 1; i <= len(g.Players); i++ {
        nextIndex := (currentIndex + i) % len(g.Players)
        if g.Players[nextIndex].IsAlive() {
            g.CurrentPlayerID = g.Players[nextIndex].ID
            g.UpdatedAt = time.Now()
            return nil
        }
    }

    return ErrNoAlivePlayer
}
```

### 步驟 1.3：建立 Value Object

**檔案**：`internal/domain/valueobject/card_color.go`

```go
package valueobject

// CardColor 卡片顏色值物件
type CardColor string

const (
    CardColorRed    CardColor = "red"
    CardColorBlue   CardColor = "blue"
    CardColorBlack  CardColor = "black"
)

// IsValid 驗證顏色是否有效
func (c CardColor) IsValid() bool {
    switch c {
    case CardColorRed, CardColorBlue, CardColorBlack:
        return true
    default:
        return false
    }
}

// String 實作 Stringer 介面
func (c CardColor) String() string {
    return string(c)
}
```

### 步驟 1.4：定義 Repository 介面

**檔案**：`internal/domain/repository/game_repository.go`

```go
package repository

import (
    "context"
    "the-message/internal/domain/entity"
)

// GameRepository 遊戲倉儲介面
// 定義於 domain 層，由 infrastructure 層實作
type GameRepository interface {
    // Create 建立遊戲
    Create(ctx context.Context, game *entity.Game) error

    // FindByID 根據 ID 查詢遊戲
    FindByID(ctx context.Context, id int) (*entity.Game, error)

    // FindByToken 根據 Token 查詢遊戲
    FindByToken(ctx context.Context, token string) (*entity.Game, error)

    // FindWithPlayers 查詢遊戲及其玩家
    FindWithPlayers(ctx context.Context, id int) (*entity.Game, error)

    // Update 更新遊戲
    Update(ctx context.Context, game *entity.Game) error

    // Delete 刪除遊戲
    Delete(ctx context.Context, id int) error
}
```

### 步驟 2.1：識別 Use Case

分析現有 Service 方法，整理為 Use Case 清單：

| 現有方法 | Use Case 名稱 | 輸入 | 輸出 |
|---------|--------------|------|------|
| `GameService.InitGame` | `CreateGameUseCase` | 玩家名單 | Game |
| `GameService.GetGame` | `GetGameUseCase` | GameID | Game |
| `PlayerService.PlayCard` | `PlayCardUseCase` | PlayerID, CardID | Result |
| `PlayerService.TransmitIntelligence` | `TransmitIntelligenceUseCase` | PlayerID, CardID, TargetID | Result |
| `PlayerService.AcceptCard` | `AcceptCardUseCase` | PlayerID, CardID | Result |

### 步驟 2.2：建立 Use Case（範例：CreateGame）

**檔案**：`internal/usecase/game/create_game.go`

```go
package game

import (
    "context"
    "crypto/rand"
    "encoding/hex"

    "the-message/internal/domain/entity"
    "the-message/internal/domain/repository"
    "the-message/internal/usecase/dto"
    "the-message/pkg/errors"
)

// CreateGameInput 建立遊戲輸入
type CreateGameInput struct {
    Players []string
}

// CreateGameOutput 建立遊戲輸出
type CreateGameOutput struct {
    GameID  int
    Token   string
    Players []dto.PlayerDTO
}

// CreateGameUseCase 建立遊戲用例
type CreateGameUseCase interface {
    Execute(ctx context.Context, input CreateGameInput) (*CreateGameOutput, error)
}

// createGameUseCase 實作
type createGameUseCase struct {
    gameRepo   repository.GameRepository
    playerRepo repository.PlayerRepository
    cardRepo   repository.CardRepository
    deckRepo   repository.DeckRepository
}

// NewCreateGameUseCase 建立用例實例
func NewCreateGameUseCase(
    gameRepo repository.GameRepository,
    playerRepo repository.PlayerRepository,
    cardRepo repository.CardRepository,
    deckRepo repository.DeckRepository,
) CreateGameUseCase {
    return &createGameUseCase{
        gameRepo:   gameRepo,
        playerRepo: playerRepo,
        cardRepo:   cardRepo,
        deckRepo:   deckRepo,
    }
}

// Execute 執行建立遊戲
func (uc *createGameUseCase) Execute(ctx context.Context, input CreateGameInput) (*CreateGameOutput, error) {
    // 1. 驗證輸入
    if len(input.Players) < 3 || len(input.Players) > 9 {
        return nil, errors.NewValidationError("players count must be between 3 and 9")
    }

    // 2. 生成遊戲 Token
    token, err := generateToken()
    if err != nil {
        return nil, errors.Wrap(err, "failed to generate token")
    }

    // 3. 建立遊戲實體
    game, err := entity.NewGame(token, input.Players)
    if err != nil {
        return nil, errors.Wrap(err, "failed to create game entity")
    }

    // 4. 儲存遊戲
    if err := uc.gameRepo.Create(ctx, game); err != nil {
        return nil, errors.Wrap(err, "failed to save game")
    }

    // 5. 初始化牌組
    if err := uc.initializeDeck(ctx, game.ID); err != nil {
        return nil, errors.Wrap(err, "failed to initialize deck")
    }

    // 6. 發放初始手牌
    if err := uc.dealInitialCards(ctx, game); err != nil {
        return nil, errors.Wrap(err, "failed to deal initial cards")
    }

    // 7. 分配身份卡
    if err := uc.assignIdentityCards(ctx, game); err != nil {
        return nil, errors.Wrap(err, "failed to assign identity cards")
    }

    // 8. 組裝輸出
    output := &CreateGameOutput{
        GameID:  game.ID,
        Token:   game.Token,
        Players: make([]dto.PlayerDTO, len(game.Players)),
    }

    for i, p := range game.Players {
        output.Players[i] = dto.ToPlayerDTO(&p)
    }

    return output, nil
}

// initializeDeck 初始化牌組
func (uc *createGameUseCase) initializeDeck(ctx context.Context, gameID int) error {
    cards, err := uc.cardRepo.FindAll(ctx)
    if err != nil {
        return err
    }

    // 洗牌邏輯
    shuffledCards := shuffleCards(cards)

    for _, card := range shuffledCards {
        deck := &entity.Deck{
            GameID: gameID,
            CardID: card.ID,
        }
        if err := uc.deckRepo.Create(ctx, deck); err != nil {
            return err
        }
    }

    return nil
}

// dealInitialCards 發放初始手牌
func (uc *createGameUseCase) dealInitialCards(ctx context.Context, game *entity.Game) error {
    const initialHandSize = 3

    for i := range game.Players {
        for j := 0; j < initialHandSize; j++ {
            card, err := uc.deckRepo.DrawCard(ctx, game.ID)
            if err != nil {
                return err
            }
            game.Players[i].HandCards = append(game.Players[i].HandCards, *card)
        }
    }

    return nil
}

// assignIdentityCards 分配身份卡
func (uc *createGameUseCase) assignIdentityCards(ctx context.Context, game *entity.Game) error {
    identities := entity.GenerateIdentityCards(len(game.Players))

    for i := range game.Players {
        game.Players[i].IdentityCard = identities[i]
        if err := uc.playerRepo.Update(ctx, &game.Players[i]); err != nil {
            return err
        }
    }

    return nil
}

// generateToken 生成遊戲 Token
func generateToken() (string, error) {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

// shuffleCards 洗牌（Fisher-Yates 算法）
func shuffleCards(cards []entity.Card) []entity.Card {
    // ... 洗牌實作
    return cards
}
```

### 步驟 2.3：建立 DTO

**檔案**：`internal/usecase/dto/game_dto.go`

```go
package dto

import (
    "time"
    "the-message/internal/domain/entity"
)

// GameDTO 遊戲資料傳輸物件
type GameDTO struct {
    ID              int          `json:"id"`
    Token           string       `json:"token"`
    Status          string       `json:"status"`
    CurrentPlayerID int          `json:"currentPlayerId"`
    Players         []PlayerDTO  `json:"players"`
    CreatedAt       time.Time    `json:"createdAt"`
}

// ToGameDTO 將 Entity 轉換為 DTO
func ToGameDTO(game *entity.Game) GameDTO {
    dto := GameDTO{
        ID:              game.ID,
        Token:           game.Token,
        Status:          string(game.Status),
        CurrentPlayerID: game.CurrentPlayerID,
        CreatedAt:       game.CreatedAt,
    }

    for _, p := range game.Players {
        dto.Players = append(dto.Players, ToPlayerDTO(&p))
    }

    return dto
}

// PlayerDTO 玩家資料傳輸物件
type PlayerDTO struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    OrderNumber  int    `json:"orderNumber"`
    Status       string `json:"status"`
    IdentityCard string `json:"identityCard,omitempty"` // 只對玩家本人顯示
    HandCards    []CardDTO `json:"handCards,omitempty"`
}

// ToPlayerDTO 將 Entity 轉換為 DTO
func ToPlayerDTO(player *entity.Player) PlayerDTO {
    dto := PlayerDTO{
        ID:          player.ID,
        Name:        player.Name,
        OrderNumber: player.OrderNumber,
        Status:      string(player.Status),
    }

    for _, c := range player.HandCards {
        dto.HandCards = append(dto.HandCards, ToCardDTO(&c))
    }

    return dto
}

// CardDTO 卡片資料傳輸物件
type CardDTO struct {
    ID               int    `json:"id"`
    Name             string `json:"name"`
    Color            string `json:"color"`
    IntelligenceType int    `json:"intelligenceType"`
}

// ToCardDTO 將 Entity 轉換為 DTO
func ToCardDTO(card *entity.Card) CardDTO {
    return CardDTO{
        ID:               card.ID,
        Name:             card.Name,
        Color:            string(card.Color),
        IntelligenceType: card.IntelligenceType,
    }
}
```

### 步驟 3.1：重構 Handler

**檔案**：`internal/adapter/http/handler/game_handler.go`

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "the-message/internal/adapter/http/request"
    "the-message/internal/adapter/http/response"
    "the-message/internal/usecase/game"
    "the-message/pkg/errors"
)

// GameHandler 遊戲 HTTP 處理器
type GameHandler struct {
    createGameUC game.CreateGameUseCase
    getGameUC    game.GetGameUseCase
}

// NewGameHandler 建立處理器
func NewGameHandler(
    createGameUC game.CreateGameUseCase,
    getGameUC game.GetGameUseCase,
) *GameHandler {
    return &GameHandler{
        createGameUC: createGameUC,
        getGameUC:    getGameUC,
    }
}

// RegisterRoutes 註冊路由
func (h *GameHandler) RegisterRoutes(r *gin.RouterGroup) {
    games := r.Group("/games")
    {
        games.POST("", h.CreateGame)
        games.GET("/:id", h.GetGame)
    }
}

// CreateGame 建立遊戲
// @Summary 建立新遊戲
// @Tags games
// @Accept json
// @Produce json
// @Param request body request.CreateGameRequest true "遊戲資訊"
// @Success 201 {object} response.GameResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/games [post]
func (h *GameHandler) CreateGame(c *gin.Context) {
    // 1. 綁定並驗證請求
    var req request.CreateGameRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request body")
        return
    }

    if err := req.Validate(); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    // 2. 轉換為 Use Case 輸入
    input := game.CreateGameInput{
        Players: req.Players,
    }

    // 3. 執行 Use Case
    output, err := h.createGameUC.Execute(c.Request.Context(), input)
    if err != nil {
        h.handleError(c, err)
        return
    }

    // 4. 回傳結果
    response.Success(c, http.StatusCreated, response.ToGameResponse(output))
}

// GetGame 取得遊戲
// @Summary 取得遊戲詳情
// @Tags games
// @Produce json
// @Param id path int true "遊戲 ID"
// @Success 200 {object} response.GameResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/games/{id} [get]
func (h *GameHandler) GetGame(c *gin.Context) {
    id, err := parseIntParam(c, "id")
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid game ID")
        return
    }

    input := game.GetGameInput{GameID: id}
    output, err := h.getGameUC.Execute(c.Request.Context(), input)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.Success(c, http.StatusOK, response.ToGameResponse(output))
}

// handleError 統一錯誤處理
func (h *GameHandler) handleError(c *gin.Context, err error) {
    switch {
    case errors.IsValidationError(err):
        response.Error(c, http.StatusBadRequest, err.Error())
    case errors.IsNotFoundError(err):
        response.Error(c, http.StatusNotFound, err.Error())
    case errors.IsConflictError(err):
        response.Error(c, http.StatusConflict, err.Error())
    default:
        // 記錄內部錯誤
        response.Error(c, http.StatusInternalServerError, "Internal server error")
    }
}

// parseIntParam 解析整數路徑參數
func parseIntParam(c *gin.Context, name string) (int, error) {
    // ... 實作
    return 0, nil
}
```

### 步驟 3.2：建立 Request/Response 結構

**檔案**：`internal/adapter/http/request/game_request.go`

```go
package request

import "the-message/pkg/errors"

// CreateGameRequest 建立遊戲請求
type CreateGameRequest struct {
    Players []string `json:"players" binding:"required"`
}

// Validate 驗證請求
func (r *CreateGameRequest) Validate() error {
    if len(r.Players) < 3 {
        return errors.NewValidationError("at least 3 players required")
    }
    if len(r.Players) > 9 {
        return errors.NewValidationError("at most 9 players allowed")
    }

    // 檢查玩家名稱不重複
    names := make(map[string]bool)
    for _, name := range r.Players {
        if name == "" {
            return errors.NewValidationError("player name cannot be empty")
        }
        if names[name] {
            return errors.NewValidationError("duplicate player name: " + name)
        }
        names[name] = true
    }

    return nil
}
```

**檔案**：`internal/adapter/http/response/response.go`

```go
package response

import "github.com/gin-gonic/gin"

// Response 通用回應結構
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo 錯誤資訊
type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// Success 成功回應
func Success(c *gin.Context, statusCode int, data interface{}) {
    c.JSON(statusCode, Response{
        Success: true,
        Data:    data,
    })
}

// Error 錯誤回應
func Error(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, Response{
        Success: false,
        Error: &ErrorInfo{
            Message: message,
        },
    })
}

// GameResponse 遊戲回應
type GameResponse struct {
    ID              int              `json:"id"`
    Token           string           `json:"token"`
    Status          string           `json:"status"`
    CurrentPlayerID int              `json:"currentPlayerId"`
    Players         []PlayerResponse `json:"players"`
}

// PlayerResponse 玩家回應
type PlayerResponse struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    OrderNumber int    `json:"orderNumber"`
    Status      string `json:"status"`
}
```

### 步驟 4.1：分離 GORM Model

**檔案**：`internal/infrastructure/persistence/mysql/model/game_model.go`

```go
package model

import (
    "time"

    "gorm.io/gorm"

    "the-message/internal/domain/entity"
)

// GameModel GORM 資料庫模型
type GameModel struct {
    gorm.Model
    ID              int           `gorm:"primaryKey;autoIncrement"`
    Token           string        `gorm:"type:varchar(64);uniqueIndex"`
    Status          string        `gorm:"type:varchar(20)"`
    CurrentPlayerID int           `gorm:"column:current_player_id"`
    Players         []PlayerModel `gorm:"foreignKey:GameID"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// TableName 指定表名
func (GameModel) TableName() string {
    return "games"
}

// ToEntity 轉換為領域實體
func (m *GameModel) ToEntity() *entity.Game {
    game := &entity.Game{
        ID:              m.ID,
        Token:           m.Token,
        Status:          entity.GameStatus(m.Status),
        CurrentPlayerID: m.CurrentPlayerID,
        CreatedAt:       m.CreatedAt,
        UpdatedAt:       m.UpdatedAt,
    }

    for _, pm := range m.Players {
        game.Players = append(game.Players, *pm.ToEntity())
    }

    return game
}

// FromEntity 從領域實體轉換
func GameModelFromEntity(e *entity.Game) *GameModel {
    model := &GameModel{
        ID:              e.ID,
        Token:           e.Token,
        Status:          string(e.Status),
        CurrentPlayerID: e.CurrentPlayerID,
        CreatedAt:       e.CreatedAt,
        UpdatedAt:       e.UpdatedAt,
    }

    for _, p := range e.Players {
        model.Players = append(model.Players, *PlayerModelFromEntity(&p))
    }

    return model
}
```

### 步驟 4.2：實作 Repository

**檔案**：`internal/infrastructure/persistence/mysql/game_repository.go`

```go
package mysql

import (
    "context"

    "gorm.io/gorm"

    "the-message/internal/domain/entity"
    "the-message/internal/domain/repository"
    "the-message/internal/infrastructure/persistence/mysql/model"
    "the-message/pkg/errors"
)

// gameRepository MySQL 遊戲倉儲實作
type gameRepository struct {
    db *gorm.DB
}

// NewGameRepository 建立倉儲
func NewGameRepository(db *gorm.DB) repository.GameRepository {
    return &gameRepository{db: db}
}

// Create 建立遊戲
func (r *gameRepository) Create(ctx context.Context, game *entity.Game) error {
    gameModel := model.GameModelFromEntity(game)

    result := r.db.WithContext(ctx).Create(gameModel)
    if result.Error != nil {
        return errors.Wrap(result.Error, "failed to create game")
    }

    // 回填生成的 ID
    game.ID = gameModel.ID
    return nil
}

// FindByID 根據 ID 查詢
func (r *gameRepository) FindByID(ctx context.Context, id int) (*entity.Game, error) {
    var gameModel model.GameModel

    result := r.db.WithContext(ctx).First(&gameModel, id)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, errors.NewNotFoundError("game not found")
        }
        return nil, errors.Wrap(result.Error, "failed to find game")
    }

    return gameModel.ToEntity(), nil
}

// FindByToken 根據 Token 查詢
func (r *gameRepository) FindByToken(ctx context.Context, token string) (*entity.Game, error) {
    var gameModel model.GameModel

    result := r.db.WithContext(ctx).Where("token = ?", token).First(&gameModel)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, errors.NewNotFoundError("game not found")
        }
        return nil, errors.Wrap(result.Error, "failed to find game")
    }

    return gameModel.ToEntity(), nil
}

// FindWithPlayers 查詢遊戲及其玩家
func (r *gameRepository) FindWithPlayers(ctx context.Context, id int) (*entity.Game, error) {
    var gameModel model.GameModel

    result := r.db.WithContext(ctx).
        Preload("Players").
        First(&gameModel, id)

    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, errors.NewNotFoundError("game not found")
        }
        return nil, errors.Wrap(result.Error, "failed to find game with players")
    }

    return gameModel.ToEntity(), nil
}

// Update 更新遊戲
func (r *gameRepository) Update(ctx context.Context, game *entity.Game) error {
    gameModel := model.GameModelFromEntity(game)

    result := r.db.WithContext(ctx).Save(gameModel)
    if result.Error != nil {
        return errors.Wrap(result.Error, "failed to update game")
    }

    return nil
}

// Delete 刪除遊戲
func (r *gameRepository) Delete(ctx context.Context, id int) error {
    result := r.db.WithContext(ctx).Delete(&model.GameModel{}, id)
    if result.Error != nil {
        return errors.Wrap(result.Error, "failed to delete game")
    }

    if result.RowsAffected == 0 {
        return errors.NewNotFoundError("game not found")
    }

    return nil
}
```

### 步驟 5.1：更新 main.go（DI 組裝）

**檔案**：`cmd/app/main.go`

```go
package main

import (
    "log"

    "github.com/gin-gonic/gin"

    // Infrastructure
    "the-message/internal/infrastructure/config"
    "the-message/internal/infrastructure/persistence/mysql"

    // Use Cases
    gameUC "the-message/internal/usecase/game"
    playerUC "the-message/internal/usecase/player"

    // Adapters
    "the-message/internal/adapter/http/handler"
    "the-message/internal/adapter/http/middleware"
)

func main() {
    // 1. 載入配置
    cfg := config.Load()

    // 2. 初始化資料庫連線
    db, err := mysql.NewConnection(cfg.Database)
    if err != nil {
        log.Fatalf("Failed to connect database: %v", err)
    }

    // 3. 初始化 Repository
    gameRepo := mysql.NewGameRepository(db)
    playerRepo := mysql.NewPlayerRepository(db)
    cardRepo := mysql.NewCardRepository(db)
    deckRepo := mysql.NewDeckRepository(db)

    // 4. 初始化 Use Cases
    createGameUC := gameUC.NewCreateGameUseCase(gameRepo, playerRepo, cardRepo, deckRepo)
    getGameUC := gameUC.NewGetGameUseCase(gameRepo)
    playCardUC := playerUC.NewPlayCardUseCase(playerRepo, cardRepo)

    // 5. 初始化 HTTP Handlers
    gameHandler := handler.NewGameHandler(createGameUC, getGameUC)
    playerHandler := handler.NewPlayerHandler(playCardUC)
    healthHandler := handler.NewHealthHandler()

    // 6. 設定 Gin Router
    r := gin.Default()

    // 7. 註冊中間件
    r.Use(middleware.CORS())
    r.Use(middleware.Recovery())
    r.Use(middleware.Logger())

    // 8. 註冊路由
    api := r.Group("/api/v1")
    {
        gameHandler.RegisterRoutes(api)
        playerHandler.RegisterRoutes(api)
        healthHandler.RegisterRoutes(api)
    }

    // 9. 啟動伺服器
    log.Printf("Server starting on port %s", cfg.Server.Port)
    if err := r.Run(":" + cfg.Server.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

---

## 6. 程式碼範例與模板

### 6.1 Entity 模板

```go
package entity

import "time"

// EntityName 領域實體
type EntityName struct {
    ID        int
    // 其他欄位...
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewEntityName 工廠方法
func NewEntityName(/* params */) (*EntityName, error) {
    // 1. 驗證參數
    // 2. 建立實體
    // 3. 回傳
    return &EntityName{}, nil
}

// BusinessMethod 業務方法
func (e *EntityName) BusinessMethod() error {
    // 業務邏輯
    return nil
}
```

### 6.2 Use Case 模板

```go
package usecase

import (
    "context"
    "the-message/internal/domain/repository"
)

// UseCaseNameInput 輸入
type UseCaseNameInput struct {
    // 輸入欄位
}

// UseCaseNameOutput 輸出
type UseCaseNameOutput struct {
    // 輸出欄位
}

// UseCaseNameUseCase 介面
type UseCaseNameUseCase interface {
    Execute(ctx context.Context, input UseCaseNameInput) (*UseCaseNameOutput, error)
}

// useCaseNameUseCase 實作
type useCaseNameUseCase struct {
    repo repository.SomeRepository
}

// NewUseCaseNameUseCase 建構函式
func NewUseCaseNameUseCase(repo repository.SomeRepository) UseCaseNameUseCase {
    return &useCaseNameUseCase{repo: repo}
}

// Execute 執行
func (uc *useCaseNameUseCase) Execute(ctx context.Context, input UseCaseNameInput) (*UseCaseNameOutput, error) {
    // 1. 驗證輸入
    // 2. 執行業務邏輯
    // 3. 組裝輸出
    return &UseCaseNameOutput{}, nil
}
```

### 6.3 Handler 模板

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "the-message/internal/adapter/http/response"
    "the-message/internal/usecase"
)

// HandlerName HTTP 處理器
type HandlerName struct {
    someUC usecase.SomeUseCase
}

// NewHandlerName 建構函式
func NewHandlerName(someUC usecase.SomeUseCase) *HandlerName {
    return &HandlerName{someUC: someUC}
}

// RegisterRoutes 註冊路由
func (h *HandlerName) RegisterRoutes(r *gin.RouterGroup) {
    group := r.Group("/resource")
    {
        group.POST("", h.Create)
        group.GET("/:id", h.GetByID)
    }
}

// Create 建立資源
func (h *HandlerName) Create(c *gin.Context) {
    // 1. 綁定請求
    // 2. 驗證
    // 3. 呼叫 Use Case
    // 4. 回傳結果
    response.Success(c, http.StatusCreated, nil)
}

// GetByID 取得資源
func (h *HandlerName) GetByID(c *gin.Context) {
    // 1. 解析參數
    // 2. 呼叫 Use Case
    // 3. 回傳結果
    response.Success(c, http.StatusOK, nil)
}
```

### 6.4 Repository 實作模板

```go
package mysql

import (
    "context"

    "gorm.io/gorm"

    "the-message/internal/domain/entity"
    "the-message/internal/domain/repository"
    "the-message/internal/infrastructure/persistence/mysql/model"
    "the-message/pkg/errors"
)

// entityRepository 倉儲實作
type entityRepository struct {
    db *gorm.DB
}

// NewEntityRepository 建構函式
func NewEntityRepository(db *gorm.DB) repository.EntityRepository {
    return &entityRepository{db: db}
}

// Create 建立
func (r *entityRepository) Create(ctx context.Context, entity *entity.Entity) error {
    m := model.FromEntity(entity)
    result := r.db.WithContext(ctx).Create(m)
    if result.Error != nil {
        return errors.Wrap(result.Error, "failed to create entity")
    }
    entity.ID = m.ID
    return nil
}

// FindByID 查詢
func (r *entityRepository) FindByID(ctx context.Context, id int) (*entity.Entity, error) {
    var m model.EntityModel
    result := r.db.WithContext(ctx).First(&m, id)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, errors.NewNotFoundError("entity not found")
        }
        return nil, errors.Wrap(result.Error, "failed to find entity")
    }
    return m.ToEntity(), nil
}
```

---

## 7. 檢查清單

### 階段 1 完成檢查

- [ ] `internal/domain/entity/` 目錄已建立
- [ ] 所有 Entity 不依賴 GORM
- [ ] Value Object 已定義（GameStatus, PlayerStatus, CardColor 等）
- [ ] Repository 介面定義於 `domain/repository/`
- [ ] Domain Service 介面已定義（如需要）
- [ ] 單元測試通過

### 階段 2 完成檢查

- [ ] `internal/usecase/` 目錄已建立
- [ ] 每個 Use Case 有明確的 Input/Output
- [ ] Use Case 只依賴 domain 層介面
- [ ] DTO 定義於 `usecase/dto/`
- [ ] 單元測試通過（可 mock repository）

### 階段 3 完成檢查

- [ ] `internal/adapter/http/` 目錄已建立
- [ ] Handler 只處理 HTTP 邏輯
- [ ] Request 驗證邏輯獨立
- [ ] Response 格式統一
- [ ] 錯誤處理統一
- [ ] 中間件已提取

### 階段 4 完成檢查

- [ ] GORM Model 與 Entity 完全分離
- [ ] Mapper 轉換邏輯正確
- [ ] Repository 實作符合介面
- [ ] 資料庫連線配置獨立
- [ ] 整合測試通過

### 階段 5 完成檢查

- [ ] `main.go` DI 組裝正確
- [ ] 所有 E2E 測試通過
- [ ] 新增單元測試覆蓋率 > 80%
- [ ] API 功能與重構前一致

### 階段 6 完成檢查

- [ ] 舊目錄已清理
- [ ] 無未使用的程式碼
- [ ] CLAUDE.md 已更新
- [ ] Swagger 文件已更新
- [ ] README 已更新

---

## 8. 常見問題與解決方案

### Q1: 為什麼要分離 GORM Model 與 Entity？

**問題**：目前 Entity 直接使用 `gorm.Model`，這會導致領域層依賴資料庫框架。

**解決方案**：
```go
// ❌ 舊做法：Entity 依賴 GORM
type Game struct {
    gorm.Model  // 違反依賴規則
    Name string
}

// ✅ 新做法：純 Entity + 獨立 Model
// domain/entity/game.go
type Game struct {
    ID   int
    Name string
}

// infrastructure/persistence/mysql/model/game_model.go
type GameModel struct {
    gorm.Model
    Name string
}

// Mapper 負責轉換
func (m *GameModel) ToEntity() *entity.Game { ... }
func FromEntity(e *entity.Game) *GameModel { ... }
```

### Q2: Use Case 之間可以互相呼叫嗎？

**答案**：一般不建議。如果有共用邏輯：

1. **提取為 Domain Service**：如果是領域邏輯
2. **提取為共用的私有函式**：如果是編排邏輯
3. **在 Handler 層組合**：如果是獨立的操作

```go
// ✅ 使用 Domain Service
type IdentityCardService interface {
    AssignCards(players []entity.Player) error
}

// CreateGameUseCase 使用 Domain Service
type createGameUseCase struct {
    identityCardSvc domain.IdentityCardService
}
```

### Q3: 如何處理跨多個 Repository 的事務？

**解決方案**：使用 Unit of Work 模式或在 Use Case 層管理事務。

```go
// 方法 1：傳入事務管理器
type TransactionManager interface {
    Execute(ctx context.Context, fn func(ctx context.Context) error) error
}

func (uc *createGameUseCase) Execute(ctx context.Context, input CreateGameInput) error {
    return uc.txManager.Execute(ctx, func(txCtx context.Context) error {
        // 在事務中執行多個 repository 操作
        if err := uc.gameRepo.Create(txCtx, game); err != nil {
            return err
        }
        if err := uc.playerRepo.CreateBatch(txCtx, players); err != nil {
            return err
        }
        return nil
    })
}
```

### Q4: 如何處理認證/授權？

**建議**：在 Adapter 層（Middleware）處理認證，在 Use Case 層處理授權。

```go
// middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 驗證 token
        userID, err := validateToken(c.GetHeader("Authorization"))
        if err != nil {
            c.AbortWithStatus(401)
            return
        }
        c.Set("userID", userID)
        c.Next()
    }
}

// usecase
func (uc *playCardUseCase) Execute(ctx context.Context, input PlayCardInput) error {
    // 授權檢查
    player, err := uc.playerRepo.FindByID(ctx, input.PlayerID)
    if err != nil {
        return err
    }
    if player.UserID != input.CurrentUserID {
        return errors.NewForbiddenError("not your turn")
    }
    // ...
}
```

### Q5: 日誌應該放在哪一層？

**建議**：
- **基礎設施日誌**（連線錯誤、效能監控）：Infrastructure 層
- **業務日誌**（重要操作記錄）：Use Case 層
- **請求日誌**（HTTP 請求/回應）：Adapter 層 Middleware

```go
// 使用依賴注入傳入 Logger
type createGameUseCase struct {
    logger Logger  // 介面定義於 pkg 或 domain
    // ...
}
```

### Q6: 如何逐步重構而不破壞現有功能？

**策略**：使用 Strangler Fig Pattern（絞殺者模式）

1. 建立新架構目錄，與舊架構並存
2. 每次重構一個 Use Case
3. 新舊路由並行運行，通過 feature flag 切換
4. 測試確認後，刪除舊程式碼

```go
// main.go
if cfg.UseNewArchitecture {
    // 新架構路由
    newHandler.RegisterRoutes(api)
} else {
    // 舊架構路由
    oldHandler.RegisterRoutes(api)
}
```

---

## 附錄 A：推薦閱讀

1. **Clean Architecture** - Robert C. Martin
2. **Implementing Domain-Driven Design** - Vaughn Vernon
3. **Go with Domain** - Three Dots Labs (免費線上)
4. **Standard Go Project Layout** - https://github.com/golang-standards/project-layout

## 附錄 B：工具推薦

| 工具 | 用途 |
|------|------|
| `wire` | 依賴注入程式碼生成 |
| `mockery` | Mock 生成工具 |
| `golangci-lint` | 程式碼品質檢查 |
| `goimports` | 自動整理 import |

## 附錄 C：專案特定注意事項

### 遊戲邏輯特殊處理

1. **狀態機**：遊戲狀態轉換應封裝於 Entity 方法中
2. **隨機性**：洗牌、抽卡邏輯應可被測試（使用種子）
3. **並發**：多玩家操作需考慮 race condition
4. **實時通訊**：SSE 事件應在 Adapter 層處理

---

> **文件維護者**：Backend Team
> **最後更新**：2026-01-11
> **版本**：1.0
