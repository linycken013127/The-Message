package main

import (
	_ "github.com/Game-as-a-Service/The-Message/cmd/app/docs"
	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/handler"
	"github.com/Game-as-a-Service/The-Message/internal/adapter/sse"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/config"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			The Message API
// @description	This is an online version of the "The Message" board game backend API
// @host			127.0.0.1:8080
func main() {
	// 1. 初始化資料庫連線（Infrastructure Layer）
	db := config.NewDatabase()

	// 2. 初始化 Gin Engine 與 SSE（Adapter Layer）
	engine := gin.Default()
	sseServer := sse.NewSSEServer()

	// 3. 初始化 Repository（Infrastructure Layer）
	// Repository 實作依賴於 GORM，但介面定義於 Domain Layer
	gameRepo := mysql.NewGameRepository(db)
	playerRepo := mysql.NewPlayerRepository(db)
	cardRepo := mysql.NewCardRepository(db)
	deckRepo := mysql.NewDeckRepository(db)
	playerCardRepo := mysql.NewPlayerCardRepository(db)
	gameProgressRepo := mysql.NewGameProgressRepository(db)
	accountRepo := mysql.NewAccountRepository(db)
	gamePlayerRepo := mysql.NewGamePlayerRepository(db)
	actionPassRepo := mysql.NewActionPassRepository(db)
	intelligenceTransferRepo := mysql.NewIntelligenceTransferRepository(db)

	// 4. 初始化 Use Cases（Use Case Layer）
	// Use Case 只依賴於 Domain Layer 的 Repository 介面
	cardUseCase := usecase.NewCardUseCase(&usecase.CardUseCaseOptions{
		CardRepo:       cardRepo,
		PlayerRepo:     playerRepo,
		PlayerCardRepo: playerCardRepo,
		GameRepo:       gameRepo,
	})

	deckUseCase := usecase.NewDeckUseCase(&usecase.DeckUseCaseOptions{
		DeckRepo:    deckRepo,
		CardUseCase: cardUseCase,
	})

	playerUseCase := usecase.NewPlayerUseCase(&usecase.PlayerUseCaseOptions{
		PlayerRepo:       playerRepo,
		PlayerCardRepo:   playerCardRepo,
		GameRepo:         gameRepo,
		GameProgressRepo: gameProgressRepo,
	})

	gameUseCase := usecase.NewGameUseCase(&usecase.GameUseCaseOptions{
		GameRepo:      gameRepo,
		PlayerUseCase: playerUseCase,
		CardUseCase:   cardUseCase,
		DeckUseCase:   deckUseCase,
	})

	// 處理循環依賴：PlayerUseCase 需要 GameUseCase 來執行 NextPlayer
	playerUseCase.SetGameUseCase(gameUseCase)

	accountUseCase := usecase.NewAccountUseCase(&usecase.AccountUseCaseOptions{
		AccountRepo: accountRepo,
	})

	gameRoomUseCase := usecase.NewGameRoomUseCase(&usecase.GameRoomUseCaseOptions{
		GameRepo:       gameRepo,
		GamePlayerRepo: gamePlayerRepo,
		AccountRepo:    accountRepo,
		PlayerUseCase:  playerUseCase,
		GameUseCase:    gameUseCase,
	})

	actionPhaseUseCase := usecase.NewActionPhaseUseCase(&usecase.ActionPhaseUseCaseOptions{
		GameRepo:       gameRepo,
		PlayerRepo:     playerRepo,
		PlayerCardRepo: playerCardRepo,
		DeckUseCase:    deckUseCase,
		ActionPassRepo: actionPassRepo,
	})

	intelligencePhaseUseCase := usecase.NewIntelligencePhaseUseCase(&usecase.IntelligencePhaseUseCaseOptions{
		GameRepo:                 gameRepo,
		PlayerRepo:               playerRepo,
		PlayerCardRepo:           playerCardRepo,
		CardRepo:                 cardRepo,
		IntelligenceTransferRepo: intelligenceTransferRepo,
	})

	// 5. 註冊 HTTP Handlers（Adapter Layer）
	// Handler 只依賴於 Use Case 介面，不直接操作 Repository
	handler.RegisterGameHandler(&handler.GameHandlerOptions{
		Engine:        engine,
		GameUseCase:   gameUseCase,
		PlayerUseCase: playerUseCase,
		SSE:           sseServer,
	})

	handler.RegisterHeartbeatHandler(&handler.HeartbeatHandler{
		Engine: engine,
	})

	handler.RegisterCardHandler(&handler.CardHandlerOptions{
		Engine:      engine,
		CardUseCase: cardUseCase,
	})

	handler.RegisterPlayerHandler(&handler.PlayerHandlerOptions{
		Engine:        engine,
		PlayerUseCase: playerUseCase,
		GameUseCase:   gameUseCase,
		SSE:           sseServer,
	})

	handler.RegisterAccountHandler(&handler.AccountHandlerOptions{
		Engine:         engine,
		AccountUseCase: accountUseCase,
	})

	handler.RegisterGameRoomHandler(&handler.GameRoomHandlerOptions{
		Engine:          engine,
		GameRoomUseCase: gameRoomUseCase,
		AccountRepo:     accountRepo,
	})

	handler.RegisterActionPhaseHandler(&handler.ActionPhaseHandlerOptions{
		Engine:             engine,
		ActionPhaseUseCase: actionPhaseUseCase,
	})

	handler.RegisterIntelligencePhaseHandler(&handler.IntelligencePhaseHandlerOptions{
		Engine:                   engine,
		IntelligencePhaseUseCase: intelligencePhaseUseCase,
	})

	// 6. Swagger 文件
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 7. 啟動伺服器
	err := engine.Run(":8080")
	if err != nil {
		return
	}
}
