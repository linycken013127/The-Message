package seeders

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql"
	"gorm.io/gorm"
)

var actionColors = map[string]map[string]int{
	entity.CardNameLockOn:      {"紅": 5, "藍": 5, "黑": 4},
	entity.CardNameLureAway:    {"紅": 2, "藍": 2, "黑": 4},
	entity.CardNameIntercept:   {"紅": 1, "藍": 1, "黑": 6},
	entity.CardNameDiversion:   {"紅": 2, "藍": 2, "黑": 2},
	entity.CardNameDecipher:    {"紅": 3, "藍": 3, "黑": 1},
	entity.CardNameBurn:        {"紅": 1, "藍": 1, "黑": 4},
	entity.CardNameSeeThrough:  {"紅": 3, "藍": 3, "黑": 6},
	entity.CardNameProbe:       {"紅": 6, "藍": 6, "黑": 6},
	entity.CardNameBlurOfTruth: {"紅": 1, "藍": 1},
}

// Run 執行所有 Seeder
func Run(db *gorm.DB) {
	SeederCards(db)
}

// SeederCards 種子卡片資料
func SeederCards(db *gorm.DB) {
	cardRepo := mysql.NewCardRepository(db)

	for actionType, colors := range actionColors {
		for color, count := range colors {
			for i := 0; i < count; i++ {
				cardRepo.CreateCard(context.TODO(), entity.NewCard(
					actionType,
					color,
					entity.ToIntelligenceType(actionType),
				))
			}
		}
	}
}
