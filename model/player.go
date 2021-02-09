package model

// Player 玩家信息
type Player struct {
	Index     int
	Name      string
	Avatar    int
	IsHuman   bool
	CardsLeft int
	Cards     CardGroup
}
