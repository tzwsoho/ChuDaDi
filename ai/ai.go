package ai

import "ChuDaDi/model"

// AI 智能出牌
type AI interface {
	AIPlay(prev model.CardGroup)
}
