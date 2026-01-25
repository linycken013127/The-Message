package repository

// GameEvent 遊戲事件類型
type GameEvent struct {
	GameID        int
	Status        string
	Message       string
	CurrentPlayer int
	NextPlayer    int
	Card          string
	Winner        string
}

// EventPublisher 事件發佈介面
type EventPublisher interface {
	// PublishGameStarted 發佈遊戲開始事件
	PublishGameStarted(event GameEvent)
	// PublishCardPlayed 發佈出牌事件
	PublishCardPlayed(event GameEvent)
	// PublishGameWon 發佈遊戲勝利事件
	PublishGameWon(event GameEvent)
	// PublishGameStatus 發佈遊戲狀態事件
	PublishGameStatus(event GameEvent)
}
