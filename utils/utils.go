package utils

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// odit depertment clark -

type Pos struct {
	x, y int
}

var X, Y, Timer, MaxCoins = 1, 1, 100, 1000000000
var Width, Height int
var coins []Pos

func GetTerm() {
	Width, Height, _ = term.GetSize(int(os.Stdout.Fd()))
	Height--
}
func MakeRandomCoins() {
	if Width*Height < MaxCoins {
		MaxCoins = Width*Height - Width
	}
	for range MaxCoins {
		coins = append(coins, Pos{x: rand.Intn(Width), y: rand.Intn(Height)})
	}
}
func returnMatchedIdx(coins []Pos, x, y int) int {
	for idx, coin := range coins {
		if coin.x == x && coin.y == y {
			return idx
		}
	}
	return -1
}
func Draw() {
	var s strings.Builder
	// dummyCoins := append([]Pos{}, coins...)
	// for i := range Height {
	// 	if i == 0 {
	// 		continue
	// 	}
	// 	for j := range Width {
	// 		isAppend := false
	// 		for idx, coin := range dummyCoins {
	// 			if i == coin.y && j == coin.x {
	// 				if rIdx := returnMatchedIdx(coins, X, Y); rIdx != -1 {
	// 					// s.WriteString("X")
	// 					coins = append(coins[:rIdx], coins[rIdx+1:]...)
	// 				} else {
	// 					dummyCoins = append(dummyCoins[:idx], dummyCoins[idx+1:]...)
	// 				}
	// 				s.WriteString("✦")
	// 				isAppend = true
	// 				break
	// 			}
	// 		}
	// 		if isAppend {
	// 			continue
	// 		}
	// 		if i == Y && j == X {
	// 			s.WriteString("✪")
	// 		} else {
	// 			s.WriteString(" ")
	// 		}
	// 	}
	// }
	// s.WriteString(strings.Repeat("-", Width))
	s.WriteString(strings.Repeat("-", Width*Height))
	s.WriteString(fmt.Sprintf("timer: %d , X : %d , Y : %d , Score : %d , Width : %d , Height : %d", Timer, X, Y, MaxCoins-len(coins), Width, Height))
	fmt.Print("\033[2J\033[H" + s.String())
}

func SetPos(input byte) bool {
	changed := true
	switch input {
	case 'j':
		if Y < Height-1 {
			Y++
		}
	case 'k':
		if Y > 1 {
			Y--
		}
	case 'h':
		if X > 0 {
			X--
		}
	case 'l':
		if X < Width-1 {
			X++
		}
	default:
		changed = false
	}
	return changed
}

func HandleTermSizeChange() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGWINCH)
	for range sigc {
		GetTerm()
		if X >= Width {
			X = Width - 1
		}
		if Y >= Height {
			Y = Height - 1
		}

		Draw()
	}
}
