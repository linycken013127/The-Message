# Task 2.1: 遊戲初始化 - 身份配置

## 功能概述

實作遊戲開始時的初始化流程，包含：
1. 根據玩家人數配置身份卡
2. 為每位玩家建立 Player 實體
3. 初始化牌堆
4. 為所有玩家發牌（每人 3 張）
5. 設定第一位行動玩家

## 身份配置表

| 人數 | 潛伏戰線 | 軍情處 | 打醬油 |
|------|----------|--------|--------|
| 3人  | 1        | 1      | 1      |
| 4人  | ❌ 不支援（身份配置不平衡） |
| 5人  | 2        | 2      | 1      |
| 6人  | 2        | 2      | 2      |
| 7人  | 2        | 2      | 3      |
| 8人  | 3        | 3      | 2      |
| 9人  | 3        | 3      | 3      |

## 實作細節

### 1. 身份卡初始化 (`player_usecase.go`)

```go
func (uc *playerUseCase) InitIdentityCards(playersCount int) []string {
    identityCards := make([]string, 0, playersCount)
    var undercover, military, bystander int

    switch playersCount {
    case 3:
        undercover, military, bystander = 1, 1, 1
    case 5:
        undercover, military, bystander = 2, 2, 1
    case 6:
        undercover, military, bystander = 2, 2, 2
    case 7:
        undercover, military, bystander = 2, 2, 3
    case 8:
        undercover, military, bystander = 3, 3, 2
    case 9:
        undercover, military, bystander = 3, 3, 3
    default:
        return identityCards // 不支援的人數返回空陣列
    }

    // 加入身份卡並洗牌
    for i := 0; i < undercover; i++ {
        identityCards = append(identityCards, entity.IdentityUndercoverFront)
    }
    for i := 0; i < military; i++ {
        identityCards = append(identityCards, entity.IdentityMilitaryAgency)
    }
    for i := 0; i < bystander; i++ {
        identityCards = append(identityCards, entity.IdentityBystander)
    }

    identityCards = ShuffleIdentityCards(identityCards)
    return identityCards
}
```

### 2. 遊戲開始驗證 (`game.go`)

新增 4 人遊戲的驗證：

```go
func (g *Game) CanStart(accountID int) error {
    if g.HostAccountID != accountID {
        return ErrNotGameHost
    }
    if g.CurrentPlayers < GameMinPlayers {
        return ErrNotEnoughPlayers
    }
    // 4 人遊戲不支援（身份配置不平衡）
    if g.CurrentPlayers == 4 {
        return ErrPlayerCountNotSupported
    }
    if g.Status != GameRoomStatusWaiting {
        return ErrGameAlreadyStarted
    }
    return nil
}

var ErrPlayerCountNotSupported = errors.New("不支援此人數配置，請選擇 3、5-9 人")
```

### 3. 完整遊戲初始化流程 (`game_room_usecase.go`)

```go
func (uc *gameRoomUseCase) StartGame(ctx context.Context, gameID int, accountID int) (*entity.Game, error) {
    // 1. 取得遊戲房
    game, err := uc.gameRepo.GetGameById(ctx, gameID)
    if err != nil {
        return nil, entity.ErrGameNotFound
    }

    // 2. 檢查是否可以開始
    if err := game.CanStart(accountID); err != nil {
        return nil, err
    }

    // 3. 產生遊戲 Token
    token, err := generateGameToken(256)
    if err != nil {
        return nil, err
    }
    game.Token = token

    // 4. 取得所有加入遊戲的玩家（GamePlayer）
    gamePlayers, err := uc.gamePlayerRepo.GetGamePlayersByGameID(ctx, gameID)
    if err != nil {
        return nil, err
    }

    // 5. 初始化身份卡
    identityCards := uc.playerUseCase.InitIdentityCards(len(gamePlayers))
    if len(identityCards) == 0 {
        return nil, entity.ErrPlayerCountNotSupported
    }

    // 6. 為每個 GamePlayer 建立 Player 並分配身份
    var firstPlayerID int
    for i, gp := range gamePlayers {
        account, err := uc.accountRepo.GetAccountById(ctx, gp.AccountID)
        if err != nil {
            return nil, err
        }

        player, err := uc.playerUseCase.CreatePlayer(ctx, entity.NewPlayer(
            account.Name,
            gameID,
            identityCards[i],
            gp.JoinOrder,
        ))
        if err != nil {
            return nil, err
        }

        // 第一位加入的玩家為第一位行動玩家
        if gp.JoinOrder == 1 {
            firstPlayerID = player.ID
        }
    }

    // 7. 初始化牌堆
    if err := uc.gameUseCase.InitDeck(ctx, game); err != nil {
        return nil, err
    }

    // 8. 為所有玩家發牌（每人 3 張）
    if err := uc.gameUseCase.DrawCardsForAllPlayers(ctx, game); err != nil {
        return nil, err
    }

    // 9. 設定第一位行動玩家並更新遊戲狀態
    game.CurrentPlayerID = firstPlayerID
    game.Status = entity.GameRoomStatusPlaying
    err = uc.gameRepo.UpdateGame(ctx, game)
    if err != nil {
        return nil, err
    }

    return game, nil
}
```

## 測試案例

新增 10 個 E2E 測試案例於 `tests/e2e/game_initialization_clean_test.go`：

| 測試案例 | 說明 |
|----------|------|
| `TestGameInit_ThreePlayers_IdentityConfig` | 3 人遊戲身份配置 (1,1,1) |
| `TestGameInit_FivePlayers_IdentityConfig` | 5 人遊戲身份配置 (2,2,1) |
| `TestGameInit_SixPlayers_IdentityConfig` | 6 人遊戲身份配置 (2,2,2) |
| `TestGameInit_SevenPlayers_IdentityConfig` | 7 人遊戲身份配置 (2,2,3) |
| `TestGameInit_EightPlayers_IdentityConfig` | 8 人遊戲身份配置 (3,3,2) |
| `TestGameInit_NinePlayers_IdentityConfig` | 9 人遊戲身份配置 (3,3,3) |
| `TestGameInit_FourPlayers_NotSupported` | 4 人遊戲不支援 |
| `TestGameInit_PlayersReceiveThreeCards` | 每位玩家發 3 張牌 |
| `TestGameInit_FirstPlayerSet` | 設定第一位行動玩家 |
| `TestGameInit_DeckHas45Cards` | 牌堆剩餘 45 張（54 - 9 = 45） |

## 架構關係

```
GameRoomHandler
       ↓
GameRoomUseCase ──→ PlayerUseCase (InitIdentityCards, CreatePlayer)
       │
       └──────────→ GameUseCase (InitDeck, DrawCardsForAllPlayers)
       │
       └──────────→ GamePlayerRepo (GetGamePlayersByGameID)
       │
       └──────────→ AccountRepo (GetAccountById)
       │
       └──────────→ GameRepo (GetGameById, UpdateGame)
```

## 異動檔案

- `internal/usecase/player_usecase.go` - 擴充 InitIdentityCards 支援 3-9 人
- `internal/domain/entity/game.go` - 新增 4 人驗證與錯誤
- `internal/usecase/game_room_usecase.go` - 完整實作 StartGame 初始化流程
- `cmd/app/main.go` - 傳入 PlayerUseCase 和 GameUseCase 給 GameRoomUseCase
- `tests/e2e/suite_clean_test.go` - 更新測試套件依賴注入
- `tests/e2e/game_initialization_clean_test.go` - 新增 10 個測試案例

## 測試結果

```
=== RUN   TestCleanArchTestSuite
--- PASS: TestCleanArchTestSuite (27.00s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_ThreePlayers_IdentityConfig (0.79s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_FivePlayers_IdentityConfig (0.82s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_SixPlayers_IdentityConfig (1.08s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_SevenPlayers_IdentityConfig (0.86s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_EightPlayers_IdentityConfig (1.05s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_NinePlayers_IdentityConfig (1.07s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_FourPlayers_NotSupported (0.54s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_PlayersReceiveThreeCards (0.78s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_FirstPlayerSet (0.77s)
    --- PASS: TestCleanArchTestSuite/TestGameInit_DeckHas45Cards (0.79s)
PASS
ok      github.com/Game-as-a-Service/The-Message/tests/e2e
```
