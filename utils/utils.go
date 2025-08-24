package utils

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// odit depertment clark -

type Pos struct {
	x, y int
}

var X, Y, Timer, MaxCoins = 1, 1, 10, 30
var Width, Height int
var CoinsMap = make(map[Pos]bool)
var B []byte
var InitTimer = Timer
var Instraction = "use hjkl to move, q or ESC to quit"

const (
	NC          = "\033[0m"
	BrightGreen = "\033[38;5;46m"
)

func GetTerm() error {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	parts := strings.Fields(string(out))
	Height, _ = strconv.Atoi(parts[0])
	Width, _ = strconv.Atoi(parts[1])
	Height--
	return nil
}

func MakeRandomCoins() {
	if Width*Height < MaxCoins {
		MaxCoins = Width*Height - Width
	}
	counter := 0
	for counter < MaxCoins {
		pos := Pos{x: rand.Intn(Width), y: rand.Intn(Height)}
		if _, ok := CoinsMap[pos]; ok {
			continue
		}
		CoinsMap[pos] = true
		counter++
	}
}

func InitCoins() {
	var s strings.Builder
	for coin := range CoinsMap {
		if coin.x == X && coin.y == Y {
			continue
		}
		RandomColor := fmt.Sprintf("\033[38;5;%dm", rand.Intn(255))
		fmt.Fprintf(&s, "\033[%d;%dH%s0%s", coin.y, coin.x, RandomColor, NC)
	}
	fmt.Print(s.String())
}

func Move(oldX, oldY int) {
	var s strings.Builder
	fmt.Fprintf(&s, "\033[%d;%dH\033[K%s  |  Timer: %d , X : %d , Y : %d , Score : %d , Width : %d , Height : %d", Height+1, 0, Instraction, Timer, X, Y, MaxCoins-len(CoinsMap), Width, Height)
	if !(oldX == X && oldY == Y) {
		delete(CoinsMap, Pos{x: X, y: Y})
		fmt.Fprintf(&s, "\033[%d;%dH ", oldY, oldX)
		fmt.Fprintf(&s, "\033[%d;%dH%s@%s", Y, X, BrightGreen, NC)
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

		InitCoins()
		Move(X, Y)
	}
}
func ShowResult() {
	text := fmt.Sprintf("SCORE : %d IN %ds", MaxCoins-len(CoinsMap), InitTimer-Timer)
	textWidth := len(text)
	boxWidth := textWidth + 8
	lineCount := 5

	s := strings.Builder{}
	cursorX := (Width - boxWidth) / 2
	cursorY := (Height - lineCount) / 2

	for i := range lineCount {
		s.WriteString(fmt.Sprintf("\033[%d;%dH", cursorY+i, cursorX))
		switch i {
		case 0, lineCount - 1:
			s.WriteString(strings.Repeat("-", boxWidth))
		case lineCount / 2:
			padding := (boxWidth - textWidth) / 2
			s.WriteString(strings.Repeat(" ", padding))
			s.WriteString(text)
			s.WriteString(strings.Repeat(" ", padding))
		default:
			s.WriteString(strings.Repeat(" ", boxWidth))
		}
	}
	fmt.Print(s.String())
}
