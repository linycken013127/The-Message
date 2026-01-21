# Task 2.4 - 情報接收實作

## 功能概述

情報接收是「風聲」遊戲中情報階段的核心環節。當情報傳遞到某位玩家時，該玩家可以選擇：

1. **接收情報**：將情報牌加入自己的情報區
2. **拒絕情報**：將情報傳給下一位玩家（密電/文件）或自動歸屬發送者（直達）

## API 端點

### 1. 接收情報
- **URL**: `POST /api/v1/games/:gameId/intelligence/accept`
- **請求參數**:
  - `player_id`: 接收情報的玩家 ID
- **成功回應**:
```json
{
  "success": true,
  "message": "情報已接收",
  "player_id": 2,
  "card_id": 10
}
```

### 2. 拒絕情報
- **URL**: `POST /api/v1/games/:gameId/intelligence/reject`
- **請求參數**:
  - `player_id`: 拒絕情報的玩家 ID
- **成功回應（傳給下一位玩家）**:
```json
{
  "success": true,
  "message": "情報已拒絕",
  "transfer": {
    "id": 1,
    "game_id": 1,
    "card_id": 10,
    "sender_player_id": 1,
    "current_target_player_id": 3,
    "face_up": false,
    "status": "IN_TRANSIT"
  }
}
```
- **成功回應（自動接收）**:
```json
{
  "success": true,
  "message": "情報已被發送者自動接收",
  "auto_accepted": true,
  "auto_accepted_by": 1
}
```

## 架構設計

### Use Case Layer

#### IntelligencePhaseUseCase 新增方法
```
internal/usecase/intelligence_phase_usecase.go
```
- `AcceptIntelligence`: 接收情報
  - 驗證是否有情報傳遞中
  - 驗證玩家是否為目標玩家
  - 將情報加入玩家情報區
  - 完成情報傳遞並進入行動階段
- `RejectIntelligence`: 拒絕情報
  - 驗證是否有情報傳遞中
  - 驗證玩家是否為目標玩家
  - 根據情報類型處理：
    - 直達：自動歸屬發送者
    - 密電/文件：傳給下一位玩家，若回到發送者則自動接收
- `completeIntelligenceTransfer`: 完成情報傳遞（內部方法）
  - 將卡片加入接收者情報區
  - 標記並刪除情報傳遞記錄
  - 更新遊戲狀態為行動階段

#### 新增結構體
```go
// AcceptIntelligenceResult 接收情報結果
type AcceptIntelligenceResult struct {
    PlayerID int
    CardID   int
}

// RejectIntelligenceResult 拒絕情報結果
type RejectIntelligenceResult struct {
    Transfer          *entity.IntelligenceTransfer
    AutoAccepted      bool   // 是否自動接收
    AutoAcceptedBy    int    // 自動接收的玩家 ID
    Message           string // 回傳訊息
}
```

### Adapter Layer (HTTP Handler)

#### IntelligencePhaseHandler 新增方法
```
internal/adapter/http/handler/intelligence_phase_handler.go
```
- `AcceptIntelligence`: 處理接收情報請求
- `RejectIntelligence`: 處理拒絕情報請求

## 測試案例

### E2E 測試
```
tests/e2e/intelligence_reception_clean_test.go
```

1. **TestAcceptIntelligence_Success**: 測試玩家成功接收情報
2. **TestAcceptIntelligence_NotTargetPlayer**: 測試非目標玩家嘗試接收情報
3. **TestAcceptIntelligence_NoTransferInProgress**: 測試沒有情報傳遞時嘗試接收
4. **TestRejectIntelligence_Secret**: 測試玩家拒絕接收密電情報
5. **TestRejectIntelligence_Document**: 測試玩家拒絕接收文件情報
6. **TestRejectIntelligence_BackToSender**: 測試密電情報回到發送者自動接收
7. **TestRejectIntelligence_Direct**: 測試直達情報被拒絕後自動歸屬發送者
8. **TestRejectIntelligence_NotTargetPlayer**: 測試非目標玩家嘗試拒絕情報

## 情報接收規則

### 接收情報流程
1. 情報目標玩家選擇「接收」
2. 情報牌加入該玩家的情報區（類型為 `intelligence`）
3. 情報傳遞結束，刪除傳遞記錄
4. 遊戲進入行動階段，當前玩家為接收者

### 拒絕情報流程

#### 密電/文件情報
1. 情報目標玩家選擇「拒絕」
2. 情報傳給下一位存活玩家（向右傳）
3. 若下一位是發送者，則發送者自動接收
4. 自動接收後，遊戲進入行動階段

#### 直達情報
1. 情報目標玩家選擇「拒絕」
2. 情報自動歸屬發送者
3. 遊戲進入行動階段

## 情報類型對照表

| 情報類型 | 拒絕後行為 | 自動接收條件 |
|---------|----------|------------|
| 密電    | 向右傳遞 | 回到發送者 |
| 文件    | 向右傳遞 | 回到發送者 |
| 直達    | 立即歸屬發送者 | 被拒絕時 |

## 遊戲流程整合

```
情報階段開始
    ↓
當前玩家傳遞情報牌
    ↓
目標玩家收到情報
    ↓
┌─────────────────────────────────────┐
│  選擇：接收 or 拒絕                  │
├─────────────────────────────────────┤
│ 接收:                               │
│   → 情報加入情報區                   │
│   → 進入行動階段                     │
│   → 當前玩家 = 接收者                │
├─────────────────────────────────────┤
│ 拒絕:                               │
│   密電/文件:                        │
│     → 傳給下一位玩家                 │
│     → 若回到發送者則自動接收          │
│   直達:                             │
│     → 發送者自動接收                 │
│     → 進入行動階段                   │
└─────────────────────────────────────┘
```

## 相關文件

- Task 2.1: 遊戲初始化（身份配置）- docs/task4.md
- Task 2.2: 行動階段 - docs/task5.md
- Task 2.3: 情報階段 - docs/task6.md

## 錯誤處理

| 錯誤情況 | 錯誤訊息 | HTTP 狀態碼 |
|---------|---------|------------|
| 沒有情報傳遞中 | 目前沒有情報傳遞中 | 400 |
| 非目標玩家 | 你不是當前情報的接收者 | 400 |
| 無效的遊戲 ID | 無效的遊戲 ID | 400 |
| 請求格式錯誤 | 請求格式錯誤 | 400 |
