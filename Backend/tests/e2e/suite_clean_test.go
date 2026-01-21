package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/handler"
	"github.com/Game-as-a-Service/The-Message/internal/adapter/sse"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/config"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/seeders"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// CleanArchTestSuite Clean Architecture 整合測試套件
type CleanArchTestSuite struct {
	suite.Suite
	db                       *gorm.DB
	tx                       *gorm.DB
	server                   *httptest.Server
	gameRepo                 repository.GameRepository
	playerRepo               repository.PlayerRepository
	playerCardRepo           repository.PlayerCardRepository
	accountRepo              repository.AccountRepository
	gamePlayerRepo           repository.GamePlayerRepository
	actionPassRepo           repository.ActionPassRepository
	intelligenceTransferRepo repository.IntelligenceTransferRepository
	gameUseCase              usecase.GameUseCase
	playerUseCase            usecase.PlayerUseCase
	accountUseCase           usecase.AccountUseCase
	gameRoomUseCase          usecase.GameRoomUseCase
	actionPhaseUseCase       usecase.ActionPhaseUseCase
	intelligencePhaseUseCase usecase.IntelligencePhaseUseCase
	deckUseCase              usecase.DeckUseCase
}

func (suite *CleanArchTestSuite) SetupSuite() {
	sourceURL := config.GetSourceURL()

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := config.BaseTestDSN()
	val := url.Values{}
	val.Add("multiStatements", "true")
	dsn = fmt.Sprintf("%s?%s", dsn, val.Encode())

	m, err := config.NewMigration(dsn, sourceURL)
	if err != nil {
		panic(err)
	}

	err = m.Down()
	if err != nil {
		if err.Error() == "no change" {
			fmt.Println("no change")
		} else {
			panic(err)
		}
	}

	err = m.Up()
	if err != nil {
		if err.Error() == "no change" {
			fmt.Println("no change")
		} else {
			panic(err)
		}
	}
	db := config.NewDatabase()

	seeders.SeederCards(db)

	engine := gin.Default()
	sseServer := sse.NewSSEServer()

	// 初始化 Repository（Infrastructure Layer）
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

	// 初始化 Use Cases（Use Case Layer）
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

	// 處理循環依賴
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

	// 註冊 HTTP Handlers
	handler.RegisterGameHandler(&handler.GameHandlerOptions{
		Engine:        engine,
		GameUseCase:   gameUseCase,
		PlayerUseCase: playerUseCase,
		SSE:           sseServer,
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

	server := httptest.NewServer(engine)

	suite.db = db
	suite.server = server
	suite.gameRepo = gameRepo
	suite.playerRepo = playerRepo
	suite.gameUseCase = gameUseCase
	suite.playerUseCase = playerUseCase
	suite.playerCardRepo = playerCardRepo
	suite.accountRepo = accountRepo
	suite.accountUseCase = accountUseCase
	suite.gamePlayerRepo = gamePlayerRepo
	suite.gameRoomUseCase = gameRoomUseCase
	suite.actionPassRepo = actionPassRepo
	suite.actionPhaseUseCase = actionPhaseUseCase
	suite.intelligenceTransferRepo = intelligenceTransferRepo
	suite.intelligencePhaseUseCase = intelligencePhaseUseCase
	suite.deckUseCase = deckUseCase
}

func (suite *CleanArchTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	err := sqlDB.Close()
	if err != nil {
		return
	}

	suite.server.Close()
}

func (suite *CleanArchTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()

	config.RunRefresh()
	db := config.NewDatabase()
	seeders.Run(db)
}

func (suite *CleanArchTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func TestCleanArchTestSuite(t *testing.T) {
	suite.Run(t, new(CleanArchTestSuite))
}

func (suite *CleanArchTestSuite) responseJson(resp *http.Response) map[string]interface{} {
	var responseMap map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		suite.T().Fatalf("Failed to decode JSON: %v", err)
	}
	return responseMap
}

func (suite *CleanArchTestSuite) requestJson(api string, jsonBody []byte, method string) *http.Response {
	req, err := http.NewRequest(method, suite.server.URL+api, bytes.NewBuffer(jsonBody))
	if err != nil {
		suite.T().Fatalf("Failed to send request: %v", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func (suite *CleanArchTestSuite) responseTest(resp *http.Response) interface{} {
	var responseMap interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		suite.T().Fatalf("Failed to decode JSON: %v", err)
	}
	return responseMap
}
