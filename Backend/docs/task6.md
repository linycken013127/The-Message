# Task 2.3 - 情報階段實作

## 功能概述

情報階段是「風聲」遊戲中的重要環節。在行動階段結束後，當前玩家必須傳遞一張情報牌給其他玩家。情報牌依據類型有不同的傳遞方式：

1. **密電（Secret Telegram）**：蓋牌向右傳遞給下一位存活玩家
2. **文件（Document）**：明牌向右傳遞給下一位存活玩家
3. **直達（Direct）**：蓋牌直接傳遞給指定的存活玩家

## API 端點

### 1. 傳遞情報牌
- **URL**: `POST /api/v1/games/:gameId/intelligence/pass-card`
- **請求參數**:
  - `player_id`: 傳遞情報的玩家 ID
  - `card_id`: 要傳遞的卡片 ID
  - `target_player_id`: 目標玩家 ID（僅直達情報需要）
- **成功回應**:
```json
{
  "success": true,
  "message": "情報牌已送出",
  "transfer": {
    "id": 1,
    "game_id": 1,
    "card_id": 10,
    "sender_player_id": 1,
    "current_target_player_id": 2,
    "face_up": false,
    "status": "IN_TRANSIT"
  }
}
```

### 2. 取得正在傳遞的情報
- **URL**: `GET /api/v1/games/:gameId/intelligence/active`
- **成功回應**:
```json
{
  "has_active_transfer": true,
  "transfer": {
    "id": 1,
    "game_id": 1,
    "card_id": 10,
    "sender_player_id": 1,
    "current_target_player_id": 2,
    "face_up": false,
    "status": "IN_TRANSIT"
  }
}
```

## 架構設計

### Domain Layer

#### Entity: IntelligenceTransfer
```
internal/domain/entity/intelligence_transfer.go
```
- 追蹤情報傳遞狀態
- 包含狀態常數：IN_TRANSIT, COMPLETED
- 定義錯誤類型：
  - `ErrCannotTargetSelf`: 直達情報不能指定自己
  - `ErrCannotTargetDeadPlayer`: 不能指定已死亡的玩家
  - `ErrDirectCardNeedsTarget`: 直達情報需要指定目標
  - `ErrNoIntelligenceInTransit`: 目前沒有情報傳遞中
  - `ErrNotIntelligenceTarget`: 不是當前情報的接收者
  - `ErrNotInIntelligencePhase`: 不在情報階段

#### Repository Interface: IntelligenceTransferRepository
```
internal/domain/repository/intelligence_transfer_repository.go
```
- `CreateIntelligenceTransfer`: 建立情報傳遞記錄
- `GetIntelligenceTransferByID`: 根據 ID 取得傳遞記錄
- `GetActiveTransferByGameID`: 取得遊戲中正在傳遞的情報
- `UpdateIntelligenceTransfer`: 更新傳遞狀態
- `DeleteIntelligenceTransfer`: 刪除傳遞記錄

### Infrastructure Layer

#### Model: IntelligenceTransferModel
```
internal/infrastructure/persistence/mysql/model/intelligence_transfer_model.go
```
- GORM 模型定義
- Entity <-> Model 轉換方法

#### Repository Implementation
```
internal/infrastructure/persistence/mysql/intelligence_transfer_repository.go
```
- 實作 IntelligenceTransferRepository 介面
- 使用 GORM 進行資料存取

#### Database Migration
```
database/migrations/000011_create_intelligence_transfers.up.sql
database/migrations/000011_create_intelligence_transfers.down.sql
```
- 建立 intelligence_transfers 資料表
- 包含外鍵約束和索引

### Use Case Layer

#### IntelligencePhaseUseCase
```
internal/usecase/intelligence_phase_usecase.go
```
- `PassIntelligenceCard`: 傳遞情報牌邏輯
  - 驗證遊戲狀態（進行中、情報階段）
  - 驗證當前玩家回合
  - 驗證手牌存在
  - 根據卡片類型決定傳遞方式
  - 建立傳遞記錄
  - 更新當前玩家為目標玩家
- `GetActiveTransfer`: 取得正在傳遞的情報

### Adapter Layer (HTTP Handler)

#### IntelligencePhaseHandler
```
internal/adapter/http/handler/intelligence_phase_handler.go
```
- 處理 HTTP 請求
- 參數驗證與錯誤處理
- 回傳 JSON 格式回應

## 測試案例

### E2E 測試
```
tests/e2e/intelligence_phase_clean_test.go
```

1. **TestPassIntelligenceCard_SecretTelegram**: 測試傳遞密電（蓋牌向右傳）
2. **TestPassIntelligenceCard_Document**: 測試傳遞文件（明牌向右傳）
3. **TestPassIntelligenceCard_Direct**: 測試傳遞直達（蓋牌指定玩家）
4. **TestPassIntelligenceCard_DirectWithoutTarget**: 測試直達卡未指定目標
5. **TestPassIntelligenceCard_DirectCannotTargetSelf**: 測試直達卡不能指定自己
6. **TestPassIntelligenceCard_NotCurrentPlayer**: 測試非當前玩家傳遞情報
7. **TestPassIntelligenceCard_NotInIntelligencePhase**: 測試非情報階段傳遞情報
8. **TestPassIntelligenceCard_CardNotInHand**: 測試傳遞不在手牌中的卡
9. **TestGetActiveTransfer_HasActiveTransfer**: 測試有正在傳遞的情報
10. **TestGetActiveTransfer_NoActiveTransfer**: 測試沒有正在傳遞的情報

## 情報類型對照表

| 情報類型 | IntelligenceType | 傳遞方式 | 明牌/蓋牌 |
|---------|------------------|---------|----------|
| 密電    | 1                | 向右傳遞 | 蓋牌     |
| 直達    | 2                | 指定玩家 | 蓋牌     |
| 文件    | 3                | 向右傳遞 | 明牌     |

## 卡片與情報類型對應

根據 `entity.ToIntelligenceType`:
- **密電**: 鎖定、調虎離山、試探、破譯
- **直達**: 截獲、燒毀、識破
- **文件**: 退回、真偽莫辯

## 遊戲流程

1. 行動階段結束（所有玩家都 Pass）
2. 遊戲進入情報階段
3. 當前玩家選擇一張手牌作為情報傳遞
4. 根據卡片類型決定傳遞方式：
   - 密電/文件：自動傳給右邊下一位存活玩家
   - 直達：需要指定目標玩家
5. 建立 IntelligenceTransfer 記錄
6. 當前玩家更新為情報接收者

## 相關文件

- Task 2.1: 遊戲初始化（身份配置）- docs/task4.md
- Task 2.2: 行動階段 - docs/task5.md
