package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Flushes ♠Spades/♥Hearts/♣Clubs/♦Diamonds
type Flushes int

const (
	// FlushesDiamonds ♦ Diamonds
	FlushesDiamonds Flushes = iota
	// FlushesClubs ♣ Clubs
	FlushesClubs
	// FlushesHearts ♥ Hearts
	FlushesHearts
	// FlushesSpades ♠ Spades
	FlushesSpades
)

func (f Flushes) String() string {
	switch f {
	case FlushesDiamonds:
		return "♦"

	case FlushesClubs:
		return "♣"

	case FlushesHearts:
		return "♥"

	case FlushesSpades:
		return "♠"
	}

	return ""
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

	return c.Flush.String() + " " + cardNumber
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Player 玩家信息
type Player struct {
	Index     int
	Name      string
	Avatar    int
	CardsLeft int
	Cards     CardGroup
}

// AllCards 所有 52 张牌（没有大小鬼/小丑/大王）
var AllCards CardGroup

// PrevCards 上一轮玩家出的牌
var PrevCards CardGroup

// PrevPlayerIndex 最后出牌玩家的索引
var PrevPlayerIndex int = -1

// Players 所有 4 名玩家
var Players [4]Player

// CurrentPlayerIndex 当前出牌玩家的索引
var CurrentPlayerIndex int = -1

// IsStarted 是否刚开局
var IsStarted bool

// PassCount 不出牌的玩家数量，不超过 3 人
var PassCount int

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// PrintAllCards 打印所有牌的花色和点数
func PrintAllCards() {
	for n, card := range AllCards {
		if 0 == n%4 {
			fmt.Println("")
		}

		fmt.Printf("%-2d - %-4s(%t)\t\t", n, card, card.Played)
	}

	fmt.Println("")
	fmt.Println(strings.Repeat("*", 100))
}

// PrintPlayersCards 打印所有玩家手上牌的花色和点数
func PrintPlayersCards() {
	for i := 0; i < len(Players); i++ {
		if i == CurrentPlayerIndex {
			fmt.Print("----> ")
		}

		fmt.Printf("Player %d:\n", i)

		var n int
		for j, card := range Players[i].Cards {
			if !card.Played {
				fmt.Printf("%-2d - %-4s(%t)\t\t", j, card, card.Played)

				n++
				if 0 == n%4 {
					fmt.Println("")
				}
			}
		}

		fmt.Println("")
		fmt.Println(strings.Repeat("*", 100))
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Init 初始化所有牌
func Init() {
	PassCount = 0
	IsStarted = true
	PrevCards = nil
	PrevPlayerIndex = -1
	CurrentPlayerIndex = -1

	for i := 0; i < len(Players); i++ {
		Players[i].CardsLeft = 0
		Players[i].Cards = make(CardGroup, 0)
	}

	AllCards = make(CardGroup, 0)
	for i := FlushesDiamonds; i <= FlushesSpades; i++ {
		for j := 0; j < 13; j++ {
			AllCards = append(AllCards, &Card{
				Played: false,
				Number: j,
				Flush:  Flushes(i),
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
func IsDiamond3(card *Card) bool {
	return 0 == card.Number && FlushesDiamonds == card.Flush
}

// IsSpade2 是否 ♠2
func IsSpade2(card *Card) bool {
	return 12 == card.Number && FlushesSpades == card.Flush
}

// IsSpadeA 是否 ♠A
func IsSpadeA(card *Card) bool {
	return 11 == card.Number && FlushesSpades == card.Flush
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

// CheckGroupType 出牌的类型和最大牌的点数（顺子/同花顺/三带二/四带一时有效）
func CheckGroupType(cards CardGroup) (GroupTypes, int) {
	// 单张（Single）
	if 1 == cards.Len() {
		return GroupTypesSingle, -1
	}

	// 一啤/一对（Pairs）：点数必须相同
	if 2 == cards.Len() {
		if cards[0].Number == cards[1].Number {
			return GroupTypesPair, -1
		}

		return GroupTypesInvalid, -1
	}

	// 三条（Triples）：点数必须相同
	if 3 == cards.Len() {
		if cards[0].Number == cards[1].Number && cards[1].Number == cards[2].Number {
			return GroupTypesTriple, -1
		}

		return GroupTypesInvalid, -1
	}

	// 四条（Four-of-a-kind）：点数必须相同
	if 4 == cards.Len() {
		if cards[0].Number == cards[1].Number && cards[1].Number == cards[2].Number && cards[2].Number == cards[3].Number {
			return GroupTypesQuadra, -1
		}

		return GroupTypesInvalid, -1
	}

	// 五张牌有多种出法
	if 5 == cards.Len() {
		// 保证外部调用前已排序
		// sort.Sort(cards)

		var (
			cardNo int
			gt1    GroupTypes = GroupTypesInvalid
			gt2    GroupTypes = GroupTypesInvalid
		)

		// 顺子/蛇（Straight）：点数相邻，但没有 2AKQJ/32AKQ/432AK 这三种组合
		if cards[0].Number == cards[1].Number-1 &&
			cards[1].Number == cards[2].Number-1 &&
			cards[2].Number == cards[3].Number-1 &&
			cards[3].Number == cards[4].Number-1 {

			// 8 - J，即 2AKQJ 组合，其余两种非法组合不符合索引相邻故无需判断
			if 8 == cards[0].Number {
				return GroupTypesInvalid, -1
			}

			cardNo = cards[4].Number
			gt1 = GroupTypesStraight
		} else if 11 == cards[3].Number && 12 == cards[4].Number && // 11 - A, 12 - 2
			0 == cards[0].Number && 1 == cards[1].Number && 2 == cards[2].Number { // A2345
			cardNo = cards[2].Number // 最大点数为 5
			gt1 = GroupTypesStraight
		} else if 12 == cards[4].Number && 0 == cards[0].Number &&
			1 == cards[1].Number && 2 == cards[2].Number && 3 == cards[3].Number { // 23456
			cardNo = cards[3].Number // 最大点数为 6
			gt1 = GroupTypesStraight
		}

		// 同花（Flush）：花色一致
		if cards[0].Flush == cards[1].Flush &&
			cards[1].Flush == cards[2].Flush &&
			cards[2].Flush == cards[3].Flush &&
			cards[3].Flush == cards[4].Flush {
			gt2 = GroupTypesFlush
		}

		if GroupTypesStraight == gt1 && GroupTypesFlush == gt2 { // 同花顺
			return GroupTypesFlushStraight, cardNo
		} else if GroupTypesStraight == gt1 { // 顺子
			return GroupTypesStraight, cardNo
		} else if GroupTypesFlush == gt2 { // 同花
			return GroupTypesFlush, -1
		} else { // 不是顺子也不是同花
			var (
				cardNo1 int = cards[0].Number
				cards1  int = 1
				cardNo2 int = -1
				cards2  int = 0
			)

			for i := 1; i < cards.Len(); i++ {
				if cards[i].Number == cardNo1 { // 累加相同点数牌的数量
					cards1++
				} else if -1 == cardNo2 { // 出现一种不同点数的牌
					cards2++
					cardNo2 = cards[i].Number
				} else if cardNo2 == cards[i].Number { // 累加相同点数牌的数量
					cards2++
				} else { // 出现第三种不同点数的牌，不符合出牌规则
					return GroupTypesInvalid, -1
				}
			}

			if 3 == cards1 { // 葫芦/夫佬/俘虏（Full House）：三个点数相同带二个点数相同
				return GroupTypesFull, cardNo1
			} else if 3 == cards2 {
				return GroupTypesFull, cardNo2
			} else if 4 == cards1 { // 金刚/铁扇（Four-of-a-kind）：四个点数相同带另一张牌
				return GroupTypesFourOfAKind, cardNo1
			} else if 4 == cards2 {
				return GroupTypesFourOfAKind, cardNo2
			}
		}
	}

	return GroupTypesInvalid, -1
}

// CompareCardGroup 比较牌组大小
func CompareCardGroup(prev, cur CardGroup) bool {
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
	case GroupTypesSingle:
		return cur[0].Number > prev[0].Number || (cur[0].Number == prev[0].Number && cur[0].Flush > prev[0].Flush)

	case GroupTypesPair: // ♠3♦3 > ♥3♣3
		return cur[1].Number > prev[1].Number || (cur[1].Number == prev[1].Number && cur[1].Flush > prev[1].Flush)

	case GroupTypesTriple:
		return cur[0].Number > prev[0].Number

	case GroupTypesQuadra:
		return cur[0].Number > prev[0].Number

	// 五张牌大小顺序：同花顺 > 四带一 > 三带二 > 同花 > 顺子
	case GroupTypesStraight:
		return gtCur > GroupTypesStraight ||
			(gtCur == GroupTypesStraight &&
				(cnCur > cnPrev || // 同为顺子比较最大牌的大小
					(cnCur == cnPrev && cur[4].Flush > prev[4].Flush)))

	case GroupTypesFlush:
		return gtCur > GroupTypesFlush ||
			(gtCur == GroupTypesFlush &&
				cur[0].Flush > prev[0].Flush) // 同为同花比较花色大小

	case GroupTypesFull:
		return gtCur > GroupTypesFull ||
			(gtCur == GroupTypesStraight &&
				cnCur > cnPrev) // 同为三带二比较最大牌的大小

	case GroupTypesFourOfAKind:
		return gtCur > GroupTypesFourOfAKind ||
			(gtCur == GroupTypesFourOfAKind &&
				cnCur > cnPrev) // 同为四带一比较最大牌的大小

	case GroupTypesFlushStraight:
		return gtCur == GroupTypesFlushStraight &&
			(cnCur > cnPrev || // 同为同花顺比较最大牌的大小
				(cnCur == cnPrev && cur[4].Flush > prev[4].Flush))
	}

	return false
}

// AutoPass 如果上一轮出的牌比其他任何牌都要大，则自动通过
func AutoPass(prev CardGroup) bool {
	sort.Sort(prev)
	gtPrev, _ := CheckGroupType(prev)
	switch gtPrev {
	case GroupTypesSingle: // ♠2
		if IsSpade2(prev[0]) {
			return true
		}

	case GroupTypesPair: // ♠2♥2/♠2♣2/♠2♦2
		if IsSpade2(prev[1]) {
			return true
		}

	case GroupTypesTriple: // ♠2♥2♣2/♠2♥2♦2/♠2♣2♦2
		if IsSpade2(prev[2]) {
			return true
		}

	case GroupTypesQuadra: // ♠2♥2♣2♦2
		if IsSpade2(prev[3]) {
			return true
		}

	case GroupTypesFlushStraight: // ♠A♠K♠Q♠J♠10
		if IsSpadeA(prev[4]) {
			return true
		}
	}

	// TODO 比未出的所有牌大也算最大
	return false
}

// CheckCanPlay 是否可出牌
func CheckCanPlay(prev, cur CardGroup) bool {
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

	if gt, _ := CheckGroupType(cur); GroupTypesInvalid == gt {
		fmt.Println("出牌不符合规则！")
		return false
	}

	if nil == prev { // 刚开局/上一轮无人出牌/上一轮出的牌比所有人的牌都大
		if IsStarted && !hasDiamond3 { // 刚开局，所出牌中必须带有 ♦3
			fmt.Println("所出牌中必须带有 ♦3 ！")
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

// AfterPlayed 出牌后的规则处理
func AfterPlayed(cards CardGroup) bool {
	IsStarted = false

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
		fmt.Printf("玩家 %d 胜出！\n", CurrentPlayerIndex)
		return true
	}

	return false
}

// Play 出牌
func Play(cards CardGroup, mustPlay bool) bool {
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
func GetCards() CardGroup {
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

	var cards CardGroup
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
		IsStarted = true
		for {
			cards := GetCards()
			fmt.Printf("你打出了：%s\n", cards)

			if !Play(cards, nil == PrevCards) {
				continue
			}

			if AfterPlayed(cards) { // true: 有玩家胜出
				break
			}

			if AutoPass(PrevCards) {
				PrevCards = nil
				PrevPlayerIndex = -1
				fmt.Println("出牌最大，其他人没有机会出牌，请继续出牌。。。")
				continue
			} else {
				CurrentPlayerIndex = (CurrentPlayerIndex + 1) % 4
			}

			fmt.Println(strings.Repeat("-", 100))
			PrintPlayersCards()
		}
	}
}
