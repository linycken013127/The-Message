package sse

import (
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/gin-gonic/gin"
)

// SSEEventPublisher SSE 事件發佈器
type SSEEventPublisher struct {
	event *Event
}

// NewSSEEventPublisher 建立 SSE 事件發佈器
func NewSSEEventPublisher(event *Event) repository.EventPublisher {
	return &SSEEventPublisher{event: event}
}

// PublishGameStarted 發佈遊戲開始事件
func (p *SSEEventPublisher) PublishGameStarted(event repository.GameEvent) {
	p.event.Message <- gin.H{
		"message":     event.Message,
		"status":      "started",
		"game_id":     event.GameID,
		"next_player": event.NextPlayer,
	}
}

// PublishCardPlayed 發佈出牌事件
func (p *SSEEventPublisher) PublishCardPlayed(event repository.GameEvent) {
	p.event.Message <- gin.H{
		"game_id":     event.GameID,
		"status":      event.Status,
		"message":     event.Message,
		"card":        event.Card,
		"next_player": event.NextPlayer,
	}
}

// PublishGameWon 發佈遊戲勝利事件
func (p *SSEEventPublisher) PublishGameWon(event repository.GameEvent) {
	p.event.Message <- gin.H{
		"game_id": event.GameID,
		"status":  event.Status,
		"message": event.Message,
		"winner":  event.Winner,
	}
}

// PublishGameStatus 發佈遊戲狀態事件
func (p *SSEEventPublisher) PublishGameStatus(event repository.GameEvent) {
	p.event.Message <- gin.H{
		"message":        event.Status,
		"status":         event.Status,
		"game_id":        event.GameID,
		"current_player": event.CurrentPlayer,
	}
}
