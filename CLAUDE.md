# CLAUDE.md

此檔案為 Claude Code (claude.ai/code) 在本儲存庫中工作時提供指導。

## 專案概述

「The Message」是一個線上棋盤遊戲實現，採用全端架構：
- **後端**：Go 配合 Gin web framework、GORM ORM、MySQL 資料庫
- **前端**：React + TypeScript + Vite 配合 Recoil 狀態管理

### 架構模式

後端採用清潔架構，包含以下層級：
- **HTTP 處理器** (`service/delivery/http/v1/`)：路由定義和請求/回應處理
- **服務** (`service/service/`)：GameService、PlayerService、CardService、DeckService 的業務邏輯
- **儲存庫** (`service/repository/mysql/`)：使用 GORM 的資料存取層
- **模型與列舉** (`enums/`)：遊戲實體和常數
- **資料庫** (`database/`)：MySQL 的遷移和播種程式

前端使用 Recoil atoms 進行全域狀態管理，組織結構如下：
- `states/globalCards/`：共享卡片狀態 (functionalCardStates、lineUpCardStates)
- `states/personalCards/`：玩家特定狀態
- `components/`：UI 元件
- `pages/`：遊戲頁面 (entrance、遊戲棋盤、地圖、playerSection、桌面)
- `hooks/`：卡片相關邏輯的自訂 React hooks

## 後端開發

### 前置需求
- Go 1.22+
- Docker（用於 MySQL）
- golang-migrate（用於資料庫遷移）

### 設定資料庫

```bash
cd Backend
docker-compose up -d           # 啟動 MySQL 容器
go run ./cmd/migrate/migrate.go # 執行遷移
go run ./cmd/migrate/game_card_seeder.go # 播種遊戲卡片
```

### 常見後端命令

**資料庫操作：**
```bash
go run ./cmd/migrate/migrate.go   # 執行遷移升級
go run ./cmd/migrate/rollback.go  # 回滾一次遷移
go run ./cmd/migrate/refresh.go   # 清空資料庫並重新執行遷移
go run ./cmd/migrate/game_card_seeder.go # 播種卡片資料
```

**運行伺服器：**
```bash
cd Backend
go run ./cmd/app/main.go  # 在 port 8080 啟動伺服器
```

**生成 Swagger API 文檔：**
```bash
cd Backend
swag init -g ./cmd/app/main.go -output ./cmd/app/docs
# 在 http://127.0.0.1:8080/swagger/index.html 查看
```

**程式碼品質：**
```bash
cd Backend
go mod tidy                    # 更新 go.mod 和 go.sum
goimports -l -w .            # 格式化匯入
golangci-lint run ./...       # 執行語法檢查工具
```

**執行測試：**
```bash
cd Backend
go test ./tests/e2e/...       # 執行所有 e2e 測試
go test -v ./tests/e2e/game_api_test.go # 執行特定測試檔案
```

說明：測試使用 testify/suite 進行整合測試，針對 `config/test_database.go` 中定義的測試資料庫。測試通過 `suite_test.go` 設定，並在 `SetupTest()` 中自動播種資料庫。

## 前端開發

### 設定

```bash
cd Frontend/the-message
npm install
```

### 常見前端命令

**開發：**
```bash
npm run dev      # 啟動 Vite dev 伺服器並進行熱重載
```

**構建與預覽：**
```bash
npm run build    # 為生產環境構建（執行 TypeScript 檢查 + Vite）
npm run preview  # 在本機預覽生產構建
```

**程式碼品質：**
```bash
npm run lint     # 執行 ESLint 並自動修復，警告時失敗
npm run test     # 執行 Jest 測試
```

**執行單個測試：**
```bash
npm run test -- App.test.tsx  # 執行特定測試檔案
```

## API 整合

前端在 `localhost:8080/api/v1/` 與後端通訊。主要端點：

- `POST /api/v1/games` - 建立新遊戲（需在 JSON body 中提供玩家列表）
- `GET /api/v1/cards` - 獲取卡片資料
- `GET /api/v1/players/{id}` - 玩家資訊
- `POST /api/v1/players/{id}/cards` - 玩家卡片操作

Swagger 文檔自動位於：`http://127.0.0.1:8080/swagger/index.html`

## 資料庫遷移

- 遷移檔案位置：`Backend/database/migrations/`
- 使用 golang-migrate 以保持一致性
- 建立新遷移：`migrate create -ext sql -dir Backend/database/migrations -seq <migration_name>`

## 主要依賴項與工具

**後端：**
- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM
- `github.com/golang-migrate/migrate/v4` - 資料庫遷移
- `github.com/swaggo/gin-swagger` - Swagger 文檔
- `github.com/stretchr/testify/suite` - 測試框架

**前端：**
- `react` & `react-dom` - UI framework
- `recoil` - 狀態管理
- `chakra-ui` - 元件庫
- `vite` - 構建工具
- `jest` & `@testing-library/react` - 測試

## 開發工作流程說明

- 使用 `rg` 進行快速文字搜尋，而非 grep（尊重 gitignore）
- 透過 `TearDownTest()` 中的交易回滾來隔離每個測試的資料庫狀態
- Swagger 文檔從處理器檔案的註解自動生成
- 前端構建中強制執行 TypeScript 嚴格模式
- ESLint 警告會導致構建失敗（`--max-warnings 0`）
