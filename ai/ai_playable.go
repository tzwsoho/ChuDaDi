package ai

import (
	"ChuDaDi/model"
	"ChuDaDi/rules"
)

// Playable 此 AI 只会打出当前持有的牌中，可打赢上一轮出牌的最小的牌组
type Playable struct {
}

// PickupGroup 从持有牌中选出特定类型的牌组
// 返回牌组和第一张牌的索引
func (ai Playable) PickupGroup(cards, prev model.CardGroup, gt model.GroupTypes, withCard *model.Card, offset int) (model.CardGroup, int) {
	switch gt {
	case model.GroupTypesSingle:
		if nil != withCard {
			return model.CardGroup{withCard}, 0
		}

		for _, card := range cards {
			ret := model.CardGroup{card}
			if rules.CompareCardGroup(prev, ret) {
				return ret, card.Number
			}
		}

	case model.GroupTypesPair:
		if cards.Len() < 2 {
			return nil, -1
		}

		var index int = -1
		for {
			var found bool
			for i := index + 1; i < cards.Len()-1; i++ {
				if nil != withCard && withCard.Number == cards[i].Number {
					if cards.GetByNumber(cards[i].Number, i).Len() >= 2 {
						index = i
						found = true
						break
					}
				} else if cards.GetByNumber(cards[i].Number, i).Len() >= 2 {
					index = i
					found = true
					break
				}
			}

			if !found {
				return nil, -1
			}

			ret := model.CardGroup{cards[index], cards[index+1]}
			if rules.CompareCardGroup(prev, ret) {
				return ret, index
			}

			index++
			if index >= cards.Len()-1 {
				break
			}
		}

	case model.GroupTypesTriple:
		if cards.Len() < 3 {
			return nil, -1
		}

		var index int = -1
		for {
			var found bool
			for i := index + 1; i < cards.Len()-2; i++ {
				if nil != withCard && withCard.Number == cards[i].Number {
					if cards.GetByNumber(cards[i].Number, i).Len() >= 3 {
						index = i
						found = true
						break
					}
				} else if cards.GetByNumber(cards[i].Number, i).Len() >= 3 {
					index = i
					found = true
					break
				}
			}

			if !found {
				return nil, -1
			}

			ret := model.CardGroup{cards[index], cards[index+1], cards[index+2]}
			if rules.CompareCardGroup(prev, ret) {
				return ret, index
			}

			index += 2
			if index >= cards.Len()-2 {
				break
			}
		}

	case model.GroupTypesQuadra:
		if cards.Len() < 4 {
			return nil, -1
		}

		var index int = -1
		for {
			var found bool
			for i := index + 1; i < cards.Len()-3; i++ {
				if nil != withCard && withCard.Number == cards[i].Number {
					if cards.GetByNumber(cards[i].Number, i).Len() >= 4 {
						index = i
						found = true
						break
					}
				} else if cards.GetByNumber(cards[i].Number, i).Len() >= 4 {
					index = i
					found = true
					break
				}
			}

			if !found {
				return nil, -1
			}

			ret := model.CardGroup{cards[index], cards[index+1], cards[index+2], cards[index+3]}
			if rules.CompareCardGroup(prev, ret) {
				return ret, index
			}

			index += 3
			if index >= cards.Len()-3 {
				break
			}
		}

	case model.GroupTypesStraight: // 没有 2AKQJ/32AKQ/432AK 这三种组合
		for i := offset; i < cards.Len()-4; i++ {
			if cards[i].Number >= 8 { // 8 - J
				break
			}

			var has bool = true
			ret := make([]model.CardGroup, 5)
			for j := 1; j < 5; j++ {
				ret[j] = cards.GetByNumber(cards[i].Number+j, i)
				if nil == ret[j] { // 没有对应点数的牌
					has = false
					break
				}
			}

			if has {
				if nil != withCard &&
					(withCard.Number != cards[i].Number || withCard.Flush != cards[i].Flush) &&
					(withCard.Number != ret[1][0].Number || withCard.Flush != ret[1][0].Flush) &&
					(withCard.Number != ret[2][0].Number || withCard.Flush != ret[2][0].Flush) &&
					(withCard.Number != ret[3][0].Number || withCard.Flush != ret[3][0].Flush) &&
					(withCard.Number != ret[4][0].Number || withCard.Flush != ret[4][0].Flush) {
					continue
				}

				return model.CardGroup{
					cards[i],
					ret[1][0],
					ret[2][0],
					ret[3][0],
					ret[4][0],
				}, i
			}
		}

		return nil, -1

	case model.GroupTypesFlush:
		var indices [][]int = make([][]int, 4)
		for i := offset; i < cards.Len(); i++ {
			indices[int(cards[i].Flush)] = append(indices[int(cards[i].Flush)], i)
		}

		for i := 0; i < len(indices); i++ {
			if len(indices[i]) >= 5 {
				if nil != withCard &&
					(withCard.Number != cards[indices[i][0]].Number || withCard.Flush != cards[indices[i][0]].Flush) &&
					(withCard.Number != cards[indices[i][1]].Number || withCard.Flush != cards[indices[i][1]].Flush) &&
					(withCard.Number != cards[indices[i][2]].Number || withCard.Flush != cards[indices[i][2]].Flush) &&
					(withCard.Number != cards[indices[i][3]].Number || withCard.Flush != cards[indices[i][3]].Flush) &&
					(withCard.Number != cards[indices[i][4]].Number || withCard.Flush != cards[indices[i][4]].Flush) {
					continue
				}

				return model.CardGroup{
					cards[indices[i][0]],
					cards[indices[i][1]],
					cards[indices[i][2]],
					cards[indices[i][3]],
					cards[indices[i][4]],
				}, indices[i][0]
			}
		}

	case model.GroupTypesFull:
		for i := offset; i < cards.Len()-2; i++ {
			cards1 := cards.GetByNumber(cards[i].Number, i)
			if cards1.Len() >= 3 { // 已找到三张牌组
				// 查找一对牌来搭配打出
				for j := 0; j < cards.Len(); j++ {
					// 不能使用和三张牌组一样点数的牌
					if cards[j].Number == cards[i].Number {
						continue
					}

					cards2 := cards.GetByNumber(cards[j].Number, j)
					if cards2.Len() >= 2 {
						return model.CardGroup{
							cards1[0],
							cards1[1],
							cards1[2],
							cards2[0],
							cards2[1],
						}, i
					}
				}
			}
		}

		return nil, -1

	case model.GroupTypesFourOfAKind:
		for i := offset; i < cards.Len()-3; i++ {
			cardsQuadra := cards.GetByNumber(cards[i].Number, i)
			if 4 == cardsQuadra.Len() { // 已找到四张牌组
				// 查找一张牌来搭配打出
				for j := 0; j < cards.Len(); j++ {
					// 不能使用和四张牌组一样点数的牌
					if cards[j].Number == cards[i].Number {
						continue
					}

					return model.CardGroup{
						cardsQuadra[0],
						cardsQuadra[1],
						cardsQuadra[2],
						cardsQuadra[3],
						cards[j],
					}, i
				}
			}
		}

		return nil, -1

	case model.GroupTypesFlushStraight:
		var indices [][]int = make([][]int, 4)
		for i := offset; i < cards.Len(); i++ {
			indices[int(cards[i].Flush)] = append(indices[int(cards[i].Flush)], i)
		}

		for i := 0; i < len(indices); i++ {
			if len(indices[i]) >= 5 {
				for j := 0; j < len(indices[i])-4; j++ { // 没有 2AKQJ/32AKQ/432AK 这三种组合
					if cards[indices[i][j]].Number >= 8 { // 8 - J
						break
					}

					if cards[indices[i][j]].Number == cards[indices[i][j+1]].Number-1 &&
						cards[indices[i][j+1]].Number == cards[indices[i][j+2]].Number-1 &&
						cards[indices[i][j+2]].Number == cards[indices[i][j+3]].Number-1 &&
						cards[indices[i][j+3]].Number == cards[indices[i][j+4]].Number-1 {
						return model.CardGroup{
							cards[indices[i][0]],
							cards[indices[i][1]],
							cards[indices[i][2]],
							cards[indices[i][3]],
							cards[indices[i][4]],
						}, indices[i][0]
					}
				}
			}
		}

		return nil, -1
	}

	return nil, -1
}

// PickupMost 打出最多的牌
func (ai Playable) PickupMost(cards model.CardGroup) model.CardGroup {
	gts := []model.GroupTypes{
		model.GroupTypesStraight,
		model.GroupTypesFlush,
		model.GroupTypesFull,
		model.GroupTypesTriple,
		model.GroupTypesPair,
		model.GroupTypesSingle,
	}

	var card *model.Card
	if rules.IsDiamond3(cards[0]) { // 首轮出牌必须带上 ♦3
		card = cards[0]
	}

	for _, gt := range gts {
		ret, _ := ai.PickupGroup(cards, nil, gt, card, 0)
		if nil != ret {
			return ret
		}
	}

	return nil
}

// Play 出牌
func (ai Playable) Play(cards, prev model.CardGroup) model.CardGroup {
	switch prev.Len() {
	case 0:
		return ai.PickupMost(cards)

	case 1:
		ret, _ := ai.PickupGroup(cards, prev, model.GroupTypesSingle, nil, 0)
		return ret

	case 2:
		ret, _ := ai.PickupGroup(cards, prev, model.GroupTypesPair, nil, 0)
		return ret

	case 3:
		ret, _ := ai.PickupGroup(cards, prev, model.GroupTypesTriple, nil, 0)
		return ret

	case 4:
		ret, _ := ai.PickupGroup(cards, prev, model.GroupTypesQuadra, nil, 0)
		return ret

	case 5:
		if cards.Len() < 5 {
			return nil
		}

		var offset int = 0
		gt, _ := rules.CheckGroupType(prev)
		for {
			// 根据出牌的牌型查找可以出的牌
			ret, off := ai.PickupGroup(cards, prev, gt, nil, offset)
			if nil == ret {
				// 找不到可以出的牌
				if model.GroupTypesFlushStraight == gt {
					return nil
				}

				// 查找下一类牌型
				gt++
				offset = 0
			}

			// 比较牌组是否比出牌要大
			if rules.CompareCardGroup(prev, ret) {
				return ret
			}

			// 查找下一个五张牌组
			offset = off + 1
		}
	}

	return nil
}
