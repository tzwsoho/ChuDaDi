package rules

import (
	"ChuDaDi/model"
	"fmt"
	"sort"
)

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

// CheckGroupType 出牌的类型和最大牌的点数（顺子/同花顺/三带二/四带一时有效）
func CheckGroupType(cards model.CardGroup) (model.GroupTypes, int) {
	// 单张（Single）
	if 1 == cards.Len() {
		return model.GroupTypesSingle, -1
	}

	// 一啤/一对（Pairs）：点数必须相同
	if 2 == cards.Len() {
		if cards[0].Number == cards[1].Number {
			return model.GroupTypesPair, -1
		}

		return model.GroupTypesInvalid, -1
	}

	// 三条（Triples）：点数必须相同
	if 3 == cards.Len() {
		if cards[0].Number == cards[1].Number && cards[1].Number == cards[2].Number {
			return model.GroupTypesTriple, -1
		}

		return model.GroupTypesInvalid, -1
	}

	// 四条（Four-of-a-kind）：点数必须相同
	if 4 == cards.Len() {
		if cards[0].Number == cards[1].Number && cards[1].Number == cards[2].Number && cards[2].Number == cards[3].Number {
			return model.GroupTypesQuadra, -1
		}

		return model.GroupTypesInvalid, -1
	}

	// 五张牌有多种出法
	if 5 == cards.Len() {
		// 保证外部调用前已排序
		// sort.Sort(cards)

		var (
			cardNo int
			gt1    model.GroupTypes = model.GroupTypesInvalid
			gt2    model.GroupTypes = model.GroupTypesInvalid
		)

		// 顺子/蛇（Straight）：点数相邻，但没有 2AKQJ/32AKQ/432AK 这三种组合
		if cards[0].Number == cards[1].Number-1 &&
			cards[1].Number == cards[2].Number-1 &&
			cards[2].Number == cards[3].Number-1 &&
			cards[3].Number == cards[4].Number-1 {

			// 8 - J，即 2AKQJ 组合，其余两种非法组合不符合索引相邻故无需判断
			if 8 == cards[0].Number {
				return model.GroupTypesInvalid, -1
			}

			cardNo = cards[4].Number
			gt1 = model.GroupTypesStraight
		} else if 11 == cards[3].Number && 12 == cards[4].Number && // 11 - A, 12 - 2
			0 == cards[0].Number && 1 == cards[1].Number && 2 == cards[2].Number { // A2345
			cardNo = cards[2].Number // 最大点数为 5
			gt1 = model.GroupTypesStraight
		} else if 12 == cards[4].Number && 0 == cards[0].Number &&
			1 == cards[1].Number && 2 == cards[2].Number && 3 == cards[3].Number { // 23456
			cardNo = cards[3].Number // 最大点数为 6
			gt1 = model.GroupTypesStraight
		}

		// 同花（Flush）：花色一致
		if cards[0].Flush == cards[1].Flush &&
			cards[1].Flush == cards[2].Flush &&
			cards[2].Flush == cards[3].Flush &&
			cards[3].Flush == cards[4].Flush {
			gt2 = model.GroupTypesFlush
		}

		if model.GroupTypesStraight == gt1 && model.GroupTypesFlush == gt2 { // 同花顺
			return model.GroupTypesFlushStraight, cardNo
		} else if model.GroupTypesStraight == gt1 { // 顺子
			return model.GroupTypesStraight, cardNo
		} else if model.GroupTypesFlush == gt2 { // 同花
			return model.GroupTypesFlush, -1
		} else { // 不是顺子也不是同花
			for i := 0; i < cards.Len()-2; i++ {
				if 3 == cards.GetByNumber(cards[i].Number, 0).Len() { // 葫芦/夫佬/俘虏（Full House）：三张点数相同带二张点数相同
					return model.GroupTypesFull, cards[i].Number
				} else if 4 == cards.GetByNumber(cards[i].Number, 0).Len() { // 金刚/铁扇（Four-of-a-kind）：四张点数相同带另一张牌
					return model.GroupTypesFourOfAKind, cards[i].Number
				}
			}
		}
	}

	return model.GroupTypesInvalid, -1
}

// CompareCardGroup 比较牌组大小
func CompareCardGroup(prev, cur model.CardGroup) bool {
	// 上一轮没有出牌
	if nil == prev {
		return true
	}

	// 出牌数量要与上一轮一致
	if prev.Len() != cur.Len() {
		return false
	}

	// sort.Sort(prev)
	// sort.Sort(cur)

	gtPrev, cnPrev := CheckGroupType(prev)
	gtCur, cnCur := CheckGroupType(cur)

	switch gtPrev {
	case model.GroupTypesSingle:
		return cur[0].Number > prev[0].Number || (cur[0].Number == prev[0].Number && cur[0].Flush > prev[0].Flush)

	case model.GroupTypesPair: // ♠3♦3 > ♥3♣3
		return cur[1].Number > prev[1].Number || (cur[1].Number == prev[1].Number && cur[1].Flush > prev[1].Flush)

	case model.GroupTypesTriple:
		return cur[0].Number > prev[0].Number

	case model.GroupTypesQuadra:
		return cur[0].Number > prev[0].Number

	// 五张牌大小顺序：同花顺 > 四带一 > 三带二 > 同花 > 顺子
	case model.GroupTypesStraight:
		return gtCur > model.GroupTypesStraight ||
			(gtCur == model.GroupTypesStraight &&
				(cnCur > cnPrev || // 同为顺子比较最大牌的大小
					(cnCur == cnPrev && cur[4].Flush > prev[4].Flush)))

	case model.GroupTypesFlush:
		return gtCur > model.GroupTypesFlush ||
			(gtCur == model.GroupTypesFlush &&
				cur[0].Flush > prev[0].Flush) // 同为同花比较花色大小

	case model.GroupTypesFull:
		return gtCur > model.GroupTypesFull ||
			(gtCur == model.GroupTypesStraight &&
				cnCur > cnPrev) // 同为三带二比较最大牌的大小

	case model.GroupTypesFourOfAKind:
		return gtCur > model.GroupTypesFourOfAKind ||
			(gtCur == model.GroupTypesFourOfAKind &&
				cnCur > cnPrev) // 同为四带一比较最大牌的大小

	case model.GroupTypesFlushStraight:
		return gtCur == model.GroupTypesFlushStraight &&
			(cnCur > cnPrev || // 同为同花顺比较最大牌的大小
				(cnCur == cnPrev && cur[4].Flush > prev[4].Flush))
	}

	return false
}

// AutoPass 如果上一轮出的牌比其他任何牌都要大，或比其他人持有牌的数量要多，则自动通过
func AutoPass(prev model.CardGroup) bool {
	sort.Sort(prev)
	gtPrev, _ := CheckGroupType(prev)
	switch gtPrev {
	case model.GroupTypesSingle: // ♠2
		if IsSpade2(prev[0]) {
			return true
		}

	case model.GroupTypesPair: // ♠2♥2/♠2♣2/♠2♦2
		if IsSpade2(prev[1]) {
			return true
		}

	case model.GroupTypesTriple: // ♠2♥2♣2/♠2♥2♦2/♠2♣2♦2
		if IsSpade2(prev[2]) {
			return true
		}

	case model.GroupTypesQuadra: // ♠2♥2♣2♦2
		if IsSpade2(prev[3]) {
			return true
		}

	case model.GroupTypesFlushStraight: // ♠A♠K♠Q♠J♠10
		if IsSpadeA(prev[4]) {
			return true
		}
	}

	// 上一轮打出的牌是否比其他玩家手上持有的牌都要多
	var found bool
	for i := 0; i < len(Players); i++ {
		if prev.Len() < Players[i].CardsLeft {
			found = true
			break
		}
	}

	if !found {
		return true
	}

	// TODO 比未出的所有牌大也算最大
	return false
}

// CheckCanPlay 是否可出牌
func CheckCanPlay(cards, prev, cur model.CardGroup) bool {
	// 本轮不出牌
	if nil == cur {
		return true
	}

	// 检查出牌是否符合规则
	sort.Sort(cur)

	var hasDiamond3 bool
	for _, card := range cur {
		if card.Played {
			fmt.Println("不能重复打出相同的牌！")
			return false
		}

		if IsDiamond3(card) {
			hasDiamond3 = true
		}
	}

	if gt, _ := CheckGroupType(cur); model.GroupTypesInvalid == gt {
		fmt.Println("出牌不符合规则！")
		return false
	}

	if nil == prev { // 刚开局/上一轮无人出牌/上一轮出的牌比所有人的牌都大
		if IsDiamond3(cards[0]) && !hasDiamond3 { // 开局出牌中必须带有 ♦3
			fmt.Println("开局出牌中必须带有 ♦3 ！")
			return false
		}

		// 符合规则就可以出牌
		return true
	}

	// 上一轮出的牌比其他任何牌组都要大
	if AutoPass(prev) {
		return false
	}

	// cur 必须大于 prev 才能出
	if !CompareCardGroup(prev, cur) {
		fmt.Printf("出的牌数量要与前一轮一致，而且必须比 %s 要大！\n", prev)
		return false
	}

	return true
}
