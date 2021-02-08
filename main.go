package main

import (
	"ChuDaDi/model"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// AllCards 所有 52 张牌（没有大小鬼/小丑/大王）
var AllCards model.CardGroup

// PrevCards 上一轮玩家出的牌
var PrevCards model.CardGroup

// PrevPlayerIndex 最后出牌玩家的索引
var PrevPlayerIndex int = -1

// Players 所有 4 名玩家
var Players [4]model.Player

// CurrentPlayerIndex 当前出牌玩家的索引
var CurrentPlayerIndex int = -1

// PassCount 不出牌的玩家数量，不超过 3 人
var PassCount int

////////////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// AfterPlayed 出牌后的规则处理
func AfterPlayed(cards model.CardGroup) bool {
	for _, card := range cards {
		card.Played = true
	}

	if nil != cards {
		PrevCards = cards
		PrevPlayerIndex = CurrentPlayerIndex
	}

	if 0 == cards.Len() { // 没人出牌
		PassCount++

		// 三个玩家都不出牌，又轮到自己出牌
		if PassCount >= 3 {
			PrevCards = nil
			PrevPlayerIndex = -1
		}
	} else {
		PassCount = 0
	}

	Players[CurrentPlayerIndex].CardsLeft -= cards.Len()
	if Players[CurrentPlayerIndex].CardsLeft <= 0 {
		fmt.Printf("%s玩家 %d 胜出！%s\n",
			strings.Repeat("*", 30),
			CurrentPlayerIndex,
			strings.Repeat("*", 30))

		fmt.Println(strings.Repeat("-\\|/", 25))
		return true
	}

	return false
}

// Play 出牌
func Play(cards model.CardGroup, mustPlay bool) bool {
	// 本轮必须出牌：其他人没有出牌
	if mustPlay && nil == cards {
		fmt.Println("本轮必须出牌！")
		return false
	}

	// 每次不能出多于 5 张牌
	if cards.Len() > 5 {
		fmt.Println("每次不能出多于 5 张牌！")
		return false
	}

	// 出牌数不能多于玩家持有牌数
	if cards.Len() > Players[CurrentPlayerIndex].CardsLeft {
		fmt.Println("出牌数不能多于玩家持有牌数！")
		return false
	}

	return CheckCanPlay(PrevCards, cards)
}

// GetCards 让玩家选择要出的牌
func GetCards() model.CardGroup {
	if -1 != PrevPlayerIndex {
		fmt.Printf("上一轮玩家 %d 打出了：%s\n", PrevPlayerIndex, PrevCards)
	}

	fmt.Printf("现在轮到玩家 %d 开始出牌\n请输入牌前面的数字，用空格隔开，或直接按下回车不出牌：\n", CurrentPlayerIndex)

	var (
		strIndices string
		err        error
	)
	for {
		reader := bufio.NewReader(os.Stdin)
		strIndices, err = reader.ReadString('\n')
		if nil != err {
			continue
		}

		strIndices = strings.TrimRight(strIndices, "\r\n")
		break
	}

	var cards model.CardGroup
	indices := strings.Split(strIndices, " ")
	for _, index := range indices {
		cardIndex, e := strconv.Atoi(index)
		if nil != e {
			continue
		}

		if cardIndex < 0 || cardIndex >= Players[CurrentPlayerIndex].Cards.Len() {
			continue
		}

		cards = append(cards, Players[CurrentPlayerIndex].Cards[cardIndex])
	}

	return cards
}

func main() {
	for {
		// 初始化
		Init()
		// PrintAllCards()

		// 洗牌
		Shuffle()
		// PrintAllCards()

		// 发牌
		DealOut()
		PrintPlayersCards()

		// 开始出牌
		for {
			cards := GetCards()
			fmt.Printf("玩家 %d 打出了：%s\n", CurrentPlayerIndex, cards)

			if !Play(cards, nil == PrevCards) {
				continue
			}

			if AfterPlayed(cards) { // true: 有玩家胜出
				break
			}

			if AutoPass(PrevCards) {
				PrevCards = nil
				PrevPlayerIndex = -1
				fmt.Printf("玩家 %d 的出牌最大，其他人没有出牌机会，请继续出牌。。。\n", CurrentPlayerIndex)
				continue
			} else {
				CurrentPlayerIndex = (CurrentPlayerIndex + 1) % 4

				// 跳过手上持有牌数少于出牌数的玩家
				for {
					if PrevCards.Len() > Players[CurrentPlayerIndex].CardsLeft {
						CurrentPlayerIndex = (CurrentPlayerIndex + 1) % 4
					} else {
						break
					}
				}
			}

			fmt.Println(strings.Repeat("-", 100))
			PrintPlayersCards()
		}
	}
}
