package main

import (
	"fmt"
	"strings"
)

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
				fmt.Printf("%-2d - %-4s\t\t", j, card)

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
