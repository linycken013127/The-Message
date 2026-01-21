# Task 3.2 - 勝利條件實作

## 功能概述

勝利條件是「風聲」遊戲中決定遊戲結束的核心機制。當特定陣營的玩家收集足夠的情報時，該陣營獲勝。勝利條件在死亡條件之後判定。

## 勝利規則

1. **潛伏戰線勝利**：任一潛伏戰線玩家情報區中紅色情報 >= 3 時，潛伏戰線全體獲勝
2. **軍情處勝利**：任一軍情處玩家情報區中藍色情報 >= 3 時，軍情處全體獲勝
3. **打醬油無法觸發**：打醬油的玩家即使收集足夠情報也無法單獨觸發勝利
4. **死亡優先**：死亡條件優先於勝利條件判定（先判斷死亡，再判斷勝利）

## 架構設計

### Use Case Layer

#### IntelligencePhaseUseCase 更新
```
internal/usecase/intelligence_phase_usecase.go
```

新增方法：
- `checkAndHandleVictory`: 檢查並處理勝利條件
- `countPlayerColoredIntelligence`: 計算玩家紅色和藍色情報數量
- `endGameWithWinner`: 以指定勝利者結束遊戲

更新結構：
- `AcceptIntelligenceResult`: 新增 `Winner` 欄位
- `RejectIntelligenceResult`: 新增 `Winner` 欄位

## 測試案例

### E2E 測試
```
tests/e2e/victory_conditions_clean_test.go
```

1. **TestVictory_UndercoverFront_ThreeRed**: 潛伏戰線收集 3 張紅色情報獲勝
2. **TestVictory_UndercoverFront_MoreThanThreeRed**: 潛伏戰線收集超過 3 張紅色情報獲勝
3. **TestVictory_MilitaryAgency_ThreeBlue**: 軍情處收集 3 張藍色情報獲勝
4. **TestVictory_Bystander_CannotTrigger**: 打醬油的無法單獨觸發勝利
5. **TestVictory_NotReached**: 情報數量不足時遊戲繼續
6. **TestVictory_DeathPriorityOverVictory**: 死亡條件優先於勝利條件

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
│       → 否: 遊戲繼續，不檢查勝利     │
├─────────────────────────────────────┤
│ 否:                                 │
│   → 檢查勝利條件                    │
└─────────────────────────────────────┘
    ↓
檢查勝利條件
    ↓
┌─────────────────────────────────────┐
│ 玩家身份是打醬油?                    │
├─────────────────────────────────────┤
│ 是: → 不觸發勝利，遊戲繼續           │
├─────────────────────────────────────┤
│ 否: → 檢查對應陣營勝利條件           │
│     潛伏戰線: 紅色情報 >= 3?         │
│     軍情處:   藍色情報 >= 3?         │
│       → 是: 該陣營獲勝，遊戲結束     │
│       → 否: 遊戲繼續                │
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
    PlayerDied bool   // 玩家是否死亡
    GameEnded  bool   // 遊戲是否結束（平局或勝利）
    Winner     string // 勝利者身份（潛伏戰線/軍情處），空字串表示平局或遊戲繼續
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
    GameEnded      bool   // 遊戲是否結束（平局或勝利）
    Winner         string // 勝利者身份（潛伏戰線/軍情處），空字串表示平局或遊戲繼續
}
```

## 勝利條件判定表

### 潛伏戰線
| 紅色情報數 | 結果 |
|----------|------|
| 0 | 未達成 |
| 1 | 未達成 |
| 2 | 未達成 |
| 3 | 潛伏戰線獲勝 |
| 4+ | 潛伏戰線獲勝 |

### 軍情處
| 藍色情報數 | 結果 |
|----------|------|
| 0 | 未達成 |
| 1 | 未達成 |
| 2 | 未達成 |
| 3 | 軍情處獲勝 |
| 4+ | 軍情處獲勝 |

### 打醬油
| 情報數 | 結果 |
|--------|------|
| 任意 | 無法觸發勝利 |

## 身份常數

```go
// 身份卡常數
const (
    IdentityUndercoverFront = "潛伏戰線"
    IdentityMilitaryAgency  = "軍情處"
    IdentityBystander       = "打醬油"
)
```

## 判定優先順序

1. **死亡條件**：黑色情報 >= 3 → 玩家死亡
2. **平局條件**：所有玩家都死亡 → 遊戲結束（平局）
3. **勝利條件**：陣營情報達標 → 該陣營獲勝

**重要**：死亡條件優先於勝利條件。如果玩家同時達到死亡條件和勝利條件，玩家會先死亡，不會觸發勝利。

## 相關文件

- Task 2.4: 情報接收 - docs/task7.md
- Task 3.1: 死亡條件 - docs/task8.md
- 勝利條件 Feature: features/07-victory-conditions.feature

## 新增方法說明

### checkAndHandleVictory
```go
func (uc *intelligencePhaseUseCase) checkAndHandleVictory(ctx context.Context, playerID int, gameID int) (gameEnded bool, winner string, err error)
```
- 檢查玩家是否達成勝利條件
- 打醬油身份不會觸發勝利
- 回傳遊戲是否結束及勝利者身份

### countPlayerColoredIntelligence
```go
func (uc *intelligencePhaseUseCase) countPlayerColoredIntelligence(ctx context.Context, playerID int) (redCount int, blueCount int, err error)
```
- 計算玩家情報區中的紅色和藍色情報數量

### endGameWithWinner
```go
func (uc *intelligencePhaseUseCase) endGameWithWinner(ctx context.Context, gameID int, winner string) error
```
- 以指定勝利者結束遊戲
- 設定 game.Status 為 ENDED
- 設定 game.Winner 為勝利者身份
