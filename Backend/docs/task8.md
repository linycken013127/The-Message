# Task 3.1 - 死亡條件實作

## 功能概述

死亡條件是「風聲」遊戲中的核心機制之一。當玩家收集到 3 張或以上的黑色情報時，該玩家會死亡。死亡的玩家會被跳過，且死亡條件優先於勝利條件。

## 死亡規則

1. **死亡條件**：玩家情報區中黑色情報數量 >= 3 時死亡
2. **死亡優先**：死亡條件優先於勝利條件判定
3. **跳過死亡玩家**：情報傳遞和回合切換都會跳過死亡玩家
4. **平局條件**：當所有玩家都死亡時，遊戲以平局結束

## 架構設計

### Domain Layer

#### Entity: Game 新增 Winner 欄位
```
internal/domain/entity/game.go
```
- 新增 `Winner` 欄位：勝利者身份（潛伏戰線/軍情處），空字串表示平局

#### Repository: PlayerRepository 新增 UpdatePlayer
```
internal/domain/repository/player_repository.go
```
- 新增 `UpdatePlayer` 方法：更新玩家狀態（包含死亡狀態）

### Infrastructure Layer

#### Database Migration
```
database/migrations/000012_add_winner_to_games.up.sql
database/migrations/000012_add_winner_to_games.down.sql
```
- 為 games 表新增 winner 欄位

#### Model: GameModel 更新
```
internal/infrastructure/persistence/mysql/model/game_model.go
```
- 新增 Winner 欄位
- 更新 ToEntity 和 FromEntity 方法

#### Repository Implementation
```
internal/infrastructure/persistence/mysql/player_repository.go
```
- 實作 UpdatePlayer 方法

### Use Case Layer

#### DeathUseCase（獨立用例）
```
internal/usecase/death_usecase.go
```
- `CheckAndHandleDeath`: 檢查並處理玩家死亡
- `KillPlayer`: 殺死玩家
- `CheckAllPlayersDead`: 檢查是否所有玩家都死亡
- `CountPlayerBlackIntelligence`: 計算玩家黑色情報數量
- `EndGameAsDraw`: 以平局結束遊戲

#### IntelligencePhaseUseCase 更新
```
internal/usecase/intelligence_phase_usecase.go
```
- 更新 `completeIntelligenceTransfer`：在接收情報後檢查死亡
- 新增 `checkAndHandleDeath`: 檢查並處理死亡
- 新增 `countPlayerBlackIntelligence`: 計算黑色情報數
- 新增 `checkAllPlayersDead`: 檢查所有玩家是否死亡
- 新增 `endGameAsDraw`: 平局結束遊戲
- `getNextAlivePlayerID` 已存在：自動跳過死亡玩家

## 測試案例

### E2E 測試
```
tests/e2e/death_conditions_clean_test.go
```

1. **TestDeath_ThreeBlackIntelligence**: 玩家收集 3 張黑色情報死亡
2. **TestDeath_MoreThanThreeBlack**: 玩家收集超過 3 張黑色情報死亡
3. **TestDeath_NotEnoughBlack**: 黑色情報少於 3 張時玩家保持存活
4. **TestDeath_SkipDeadPlayerTurn**: 死亡玩家被跳過回合
5. **TestDeath_SkipDeadPlayerIntelligence**: 情報傳遞時跳過死亡玩家
6. **TestDeath_PriorityOverVictory**: 死亡條件優先於勝利條件
7. **TestDeath_AllPlayersDead**: 所有玩家都死亡時遊戲結束為平局

## 遊戲流程

```
玩家接收情報
    ↓
情報加入情報區
    ↓
計算黑色情報數量
    ↓
┌─────────────────────────────────────┐
│ 黑色情報 >= 3?                       │
├─────────────────────────────────────┤
│ 是:                                 │
│   → 玩家死亡（狀態改為「死亡」）      │
│   → 檢查是否所有玩家都死亡           │
│       → 是: 遊戲結束（平局）         │
│       → 否: 遊戲繼續                │
├─────────────────────────────────────┤
│ 否:                                 │
│   → 遊戲繼續                        │
└─────────────────────────────────────┘
    ↓
進入下一回合（行動階段）
```

## 結果結構

### AcceptIntelligenceResult 更新
```go
type AcceptIntelligenceResult struct {
    PlayerID   int
    CardID     int
    PlayerDied bool // 玩家是否死亡
    GameEnded  bool // 遊戲是否結束（平局）
}
```

### RejectIntelligenceResult 更新
```go
type RejectIntelligenceResult struct {
    Transfer       *entity.IntelligenceTransfer
    AutoAccepted   bool   // 是否自動接收
    AutoAcceptedBy int    // 自動接收的玩家 ID
    Message        string // 回傳訊息
    PlayerDied     bool   // 玩家是否死亡
    GameEnded      bool   // 遊戲是否結束（平局）
}
```

## 死亡條件判定表

| 黑色情報數 | 結果 |
|----------|------|
| 0 | 存活 |
| 1 | 存活 |
| 2 | 存活 |
| 3 | 死亡 |
| 4+ | 死亡 |

## 平局條件

當遊戲中所有玩家都死亡時：
- `game.Status` 設為 `ENDED`
- `game.Winner` 設為空字串（表示平局）

## 死亡優先於勝利

即使玩家同時達成勝利條件（如潛伏戰線收集 3 紅）和死亡條件（3 黑），死亡條件優先：
1. 先判定死亡
2. 玩家死亡後不觸發勝利
3. 遊戲繼續進行

## 相關文件

- Task 2.3: 情報階段 - docs/task6.md
- Task 2.4: 情報接收 - docs/task7.md
- Task 3.2: 勝利條件（待實作）

## 資料庫變更

### 新增遷移檔案
- `000012_add_winner_to_games.up.sql`
- `000012_add_winner_to_games.down.sql`

### games 表新增欄位
| 欄位 | 類型 | 說明 |
|------|------|------|
| winner | VARCHAR(50) | 勝利者身份，空字串表示平局 |
