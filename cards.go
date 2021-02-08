package main

import (
	"ChuDaDi/model"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Init 初始化所有牌
func Init() {
	PassCount = 0
	PrevCards = nil
	PrevPlayerIndex = -1
	CurrentPlayerIndex = -1

	for i := 0; i < len(Players); i++ {
		Players[i].CardsLeft = 0
		Players[i].Cards = make(model.CardGroup, 0)
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
	// now := int64(1612754657760430000)
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

// IsDiamond3 是否 ♦3
func IsDiamond3(card *model.Card) bool {
	return 0 == card.Number && model.FlushesDiamonds == card.Flush
}

// IsSpade2 是否 ♠2
func IsSpade2(card *model.Card) bool {
	return 12 == card.Number && model.FlushesSpades == card.Flush
}

// IsSpadeA 是否 ♠A
func IsSpadeA(card *model.Card) bool {
	return 11 == card.Number && model.FlushesSpades == card.Flush
}
