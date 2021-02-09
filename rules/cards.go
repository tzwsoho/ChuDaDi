package rules

import (
	"ChuDaDi/model"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// AllCards 所有 52 张牌（没有大小鬼/小丑/大王）
var AllCards model.CardGroup

// PrevCards 上一轮玩家出的牌
var PrevCards model.CardGroup

// PrevPlayerIndex 最后出牌玩家的索引
var PrevPlayerIndex int = -1

// Players 所有 4 名玩家
var Players [4]*model.Player

// CurrentPlayerIndex 当前出牌玩家的索引
var CurrentPlayerIndex int = -1

// InitAll 初始化
func InitAll() {
	PrevCards = nil
	PrevPlayerIndex = -1
	CurrentPlayerIndex = -1

	for i := 0; i < len(Players); i++ {
		Players[i] = &model.Player{
			Index:     i,
			Name:      fmt.Sprintf("Player %d", i),
			Avatar:    i % 10,
			IsHuman:   false,
			CardsLeft: 0,
			Cards:     make(model.CardGroup, 0),
		}
	}

	AllCards = make(model.CardGroup, 0)
	for i := model.FlushesDiamonds; i <= model.FlushesSpades; i++ {
		for j := 0; j < 13; j++ {
			AllCards = append(AllCards, &model.Card{
				Played: false,
				Number: j,
				Flush:  model.Flushes(i),
			})
		}
	}
}

// Shuffle 洗牌
func Shuffle() {
	// 1612779063003819500 玩家 3 同花顺
	// now := int64(1612853908550452700)
	now := time.Now().UnixNano()
	fmt.Printf("本轮随机种子：%d\n", now)

	rand.Seed(now)
	for i := 0; i < AllCards.Len(); i++ {
		AllCards.Swap(0, rand.Intn(AllCards.Len()))
	}
}

// DealOut 发牌
func DealOut() {
	// 轮流发牌
	for i := 0; i < AllCards.Len(); i++ {
		if IsDiamond3(AllCards[i]) {
			CurrentPlayerIndex = i % 4
		}

		Players[i%4].CardsLeft++
		Players[i%4].Cards = append(Players[i%4].Cards, AllCards[i])
	}

	// 每个玩家手上的牌按从小到大排序
	for i := 0; i < len(Players); i++ {
		sort.Sort(Players[i].Cards)
	}
}

// GetPlayer 获取玩家
func GetPlayer(index int) *model.Player {
	if index < 0 || index >= len(Players) {
		return nil
	}

	return Players[index]
}

// GetCurPlayer 获取当前游戏的玩家
func GetCurPlayer() *model.Player {
	if CurrentPlayerIndex < 0 || CurrentPlayerIndex >= len(Players) {
		return nil
	}

	return Players[CurrentPlayerIndex]
}

// GetCurCards 获取当前游戏玩家手上的牌
func GetCurCards() model.CardGroup {
	if CurrentPlayerIndex < 0 || CurrentPlayerIndex >= len(Players) {
		return nil
	}

	return Players[CurrentPlayerIndex].Cards
}
