package model

import "sort"

// GroupTypes 牌组类型
type GroupTypes int

const (
	// GroupTypesInvalid 不符合规则
	GroupTypesInvalid GroupTypes = iota
	// GroupTypesSingle 单张
	GroupTypesSingle
	// GroupTypesPair 一啤/一对
	GroupTypesPair
	// GroupTypesTriple 三条
	GroupTypesTriple
	// GroupTypesQuadra 四条
	GroupTypesQuadra
	// GroupTypesStraight 顺子
	GroupTypesStraight
	// GroupTypesFlush 同花
	GroupTypesFlush
	// GroupTypesFull 三带二
	GroupTypesFull
	// GroupTypesFourOfAKind 四带一
	GroupTypesFourOfAKind
	// GroupTypesFlushStraight 同花顺
	GroupTypesFlushStraight
)

// CardGroup 一组牌
type CardGroup []*Card

// Len 牌的张数
func (cg CardGroup) Len() int {
	return len(cg)
}

// Less 比较大小，先比点数再比花色
func (cg CardGroup) Less(i, j int) bool {
	return cg[i].Number < cg[j].Number ||
		(cg[i].Number == cg[j].Number && cg[i].Flush < cg[j].Flush)
}

// Swap 交换两张牌
func (cg CardGroup) Swap(i, j int) {
	cg[i], cg[j] = cg[j], cg[i]
}

// String 列出所有牌的花色和点数
func (cg CardGroup) String() string {
	var s string
	for _, card := range cg {
		s += card.String() + " "
	}

	return s
}

// NotPlayed 获取未打出的牌
func (cg CardGroup) NotPlayed() CardGroup {
	var cgNP CardGroup
	for _, card := range cg {
		if !card.Played {
			cgNP = append(cgNP, card)
		}
	}

	sort.Sort(cgNP)
	return cgNP
}

// GetByNumber 根据牌的点数筛选对应的牌
func (cg CardGroup) GetByNumber(number int, offset int) CardGroup {
	var cards CardGroup
	for i := offset; i < cg.Len(); i++ {
		if number == cg[i].Number {
			cards = append(cards, cg[i])
		}
	}

	return cards
}
