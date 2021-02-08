package model

// CardNumber 牌面点数
const CardNumber string = "34567890JQKA2"

// Card 一张牌
type Card struct {
	Played bool    // 是否已打出
	Number int     // 点数索引：0 - 3, 1 - 4, 2 - 5, 3 - 6, 4 - 7, 5 - 8, 6 - 9, 7 - 10, 8 - J, 9 - Q, 10 - K, 11 - A, 12 - 2
	Flush  Flushes // 花色：FlushesSpades - ♠, FlushesHearts - ♥, FlushesClubs - ♣, FlushesDiamonds - ♦
}

func (c Card) String() string {
	var cardNumber string = string(CardNumber[c.Number])
	if "0" == cardNumber {
		cardNumber = "10"
	}

	return c.Flush.String() + cardNumber
}
