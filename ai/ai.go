package ai

import (
	"ChuDaDi/model"
)

// AI 智能出牌
type AI interface {
	// cards 手上持有的牌组
	// prev 上家出的牌组
	// 返回要打出的牌组
	Play(cards, prev model.CardGroup) model.CardGroup
}

// Types AI 类型
type Types int

const (
	// AITypesPlayable 只会打出可以打赢上一轮出牌的最小的牌
	AITypesPlayable = iota
)

// NewAI 获得一个 AI 示例
func NewAI(aiType Types) AI {
	switch aiType {
	case AITypesPlayable:
		return &Playable{}
	}

	return nil
}
