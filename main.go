package main

import (
	"ChuDaDi/ai"
	"ChuDaDi/model"
	"ChuDaDi/rules"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// PassCount 不出牌的玩家数量，不超过 3 人
var PassCount int

// AfterPlayed 出牌后的规则处理
func AfterPlayed(cards model.CardGroup) bool {
	for _, card := range cards {
		card.Played = true
	}

	if nil != cards {
		rules.PrevCards = cards
		rules.PrevPlayerIndex = rules.CurrentPlayerIndex
	}

	if 0 == cards.Len() { // 没人出牌
		PassCount++

		// 三个玩家都不出牌，又轮到自己出牌
		if PassCount >= 3 {
			rules.PrevCards = nil
			rules.PrevPlayerIndex = -1
		}
	} else {
		PassCount = 0
	}

	rules.Players[rules.CurrentPlayerIndex].CardsLeft -= cards.Len()
	if rules.Players[rules.CurrentPlayerIndex].CardsLeft <= 0 {
		fmt.Printf("%s玩家 %d 胜出！%s\n",
			strings.Repeat("*", 30),
			rules.CurrentPlayerIndex,
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
	if cards.Len() > rules.GetCurPlayer().CardsLeft {
		fmt.Println("出牌数不能多于玩家持有牌数！")
		return false
	}

	return rules.CheckCanPlay(rules.GetCurCards().NotPlayed(), rules.PrevCards, cards)
}

// GetCards 让玩家选择要出的牌
func GetCards() model.CardGroup {
	if -1 != rules.PrevPlayerIndex {
		fmt.Printf("上一轮玩家 %d 打出了：%s\n", rules.PrevPlayerIndex, rules.PrevCards)
	}

	fmt.Printf("现在轮到玩家 %d 开始出牌\n请输入牌前面的数字，用空格隔开，或直接按下回车不出牌：\n", rules.CurrentPlayerIndex)

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

		if cardIndex < 0 || cardIndex >= rules.GetCurCards().Len() {
			continue
		}

		cards = append(cards, rules.GetCurCards()[cardIndex])
	}

	return cards
}

func main() {
	AI := ai.NewAI(ai.AITypesPlayable).(*ai.Playable)

	for {
		// 初始化
		rules.InitAll()
		// rules.PrintAllCards()

		// 洗牌
		rules.Shuffle()
		// rules.PrintAllCards()

		// 发牌
		rules.DealOut()

		// 设置三人为 AI 玩家，一人为真人玩家
		// n := rand.Intn(len(rules.Players))
		for i := 0; i < len(rules.Players); i++ {
			// if 0 == (i+n)%4 {
			// rules.Players[i].IsHuman = true
			// } else {
			rules.Players[i].IsHuman = false
			// }
		}

		rules.PrintPlayersCards()

		// 开始出牌
		for {
			<-time.After(time.Millisecond * 100)

			var cards model.CardGroup
			if rules.GetCurPlayer().IsHuman { // 真人玩家出牌
				cards = GetCards()
			} else { // 电脑出牌
				now := time.Now().UnixNano()
				cards = AI.Play(rules.GetCurCards().NotPlayed(), rules.PrevCards)
				delta := time.Now().UnixNano() - now
				if delta > 0 {
					fmt.Printf("耗时 %d ns\n", delta)
				}
			}

			if nil == cards {
				fmt.Printf("玩家 %d 不出牌\n", rules.CurrentPlayerIndex)
			} else {
				fmt.Printf("玩家 %d 打出了：%s\n", rules.CurrentPlayerIndex, cards)
			}

			if !Play(cards, nil == rules.PrevCards) {
				continue
			}

			if AfterPlayed(cards) { // true: 有玩家胜出
				os.Exit(0)
				break
			}

			if rules.AutoPass(rules.PrevCards) {
				rules.PrevCards = nil
				rules.PrevPlayerIndex = -1
				fmt.Printf("玩家 %d 的出牌最大，其他人没有出牌机会，请继续出牌。。。\n", rules.CurrentPlayerIndex)
				continue
			} else {
				rules.CurrentPlayerIndex = (rules.CurrentPlayerIndex + 1) % 4

				// 跳过手上持有牌数少于出牌数的玩家
				for {
					if rules.PrevCards.Len() > rules.GetCurPlayer().CardsLeft {
						rules.CurrentPlayerIndex = (rules.CurrentPlayerIndex + 1) % 4
					} else {
						break
					}
				}
			}

			if 0 == rules.CurrentPlayerIndex {
				fmt.Println(strings.Repeat("-", 100))
				rules.PrintPlayersCards()
			}
		}
	}
}
