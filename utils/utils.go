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

var X, Y, Timer, MaxCoins = 1, 1, 10, 4588
var Width, Height int
var coins []Pos
var B []byte

func GetTerm() {
	Width, Height, _ = term.GetSize(int(os.Stdout.Fd()))
	Height--
}
func MakeRandomCoins() {
	if Width*Height < MaxCoins {
		MaxCoins = Width*Height - Width
	}
	var coinMap = make(map[Pos]bool)
	counter := 0
	for counter < MaxCoins {
		pos := Pos{x: rand.Intn(Width), y: rand.Intn(Height)}
		if _, ok := coinMap[pos]; ok {
			continue
		}
		coinMap[pos] = true
		coins = append(coins, pos)
		counter++
	}
}
func Move(oldX, oldY int) {
	var s strings.Builder
	fmt.Fprintf(&s, "\033[%d;%dHtimer: %d , X : %d , Y : %d , Score : %d , Width : %d , Height : %d , Coins : %d", Height+1, 0, Timer, X, Y, MaxCoins-len(coins), Width, Height, len(coins))
	if !(oldX == X && oldY == Y) {
		fmt.Fprintf(&s, "\033[%d;%dH ", oldY, oldX)
		j := 0
		for _, coin := range coins {
			if coin.x == X && coin.y == Y {
				continue
			}
			coins[j] = coin
			fmt.Fprintf(&s, "\033[%d;%dH0", coin.y, coin.x)
			j++
		}
		coins = coins[:j]
		fmt.Fprintf(&s, "\033[%d;%dH@", Y, X)
	}
	fmt.Print(s.String())
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

		Move(X, Y)
	}
}
func ShowResult() {
	text := fmt.Sprintf("SCORE : %d", MaxCoins-len(coins))
	textWidth := len(text)
	s := strings.Builder{}
	lineCount := 5
	cursorX := (Width - textWidth) / 2
	cursorY := (Height - lineCount) / 2
	s.WriteString(fmt.Sprintf("\033[%d;%dH", cursorY, cursorX))
	for i := range lineCount {
		s.WriteString(fmt.Sprintf("\033[%d;%dH", cursorY+i, cursorX))
		switch i {
		case 0, lineCount - 1:
			s.WriteString(strings.Repeat("-", textWidth*2))
		case lineCount / 2:
			s.WriteString(strings.Repeat(" ", textWidth/2))
			s.WriteString(text)
			s.WriteString(strings.Repeat(" ", textWidth/2))
		default:
			s.WriteString(strings.Repeat(" ", textWidth*2))
		}
	}
	fmt.Print(s.String())
}
