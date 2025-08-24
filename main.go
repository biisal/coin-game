package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/biisal/fun-game/utils"
	"golang.org/x/term"
)

func main() {
	utils.GetTerm()
	utils.MakeRandomCoins()
	utils.X = rand.Intn(utils.Width)
	utils.Y = rand.Intn(utils.Height)
	oldState, err := term.MakeRaw(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fmt.Print("\033[?1049l")
		fmt.Print("\x1b[?25h")
		term.Restore(int(os.Stdout.Fd()), oldState)
	}()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\033[?1049h")
	fmt.Print("\x1b[?25l")

	showResult, redraw, quit := make(chan bool), make(chan bool, 1), make(chan bool)
	redraw <- true
	go func() {
		for range utils.Timer {
			time.Sleep(time.Second)
			utils.Timer--
			redraw <- true
		}
		showResult <- true
	}()

	go func() {
		input := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(input)
			if err != nil || n == 0 {
				continue
			}
			if input[0] == 'q' {
				if utils.Timer == 0 {
					showResult <- true
					continue
				}
				showResult <- true
			}
			if input[0] == 0x1b {
				quit <- true
			}
			var oldX, oldY = utils.X, utils.Y
			if utils.SetPos(input[0]) {
				utils.Move(oldX, oldY)
			}
		}
	}()
	utils.Move(0, 0)
	for {
		select {
		case <-showResult:
			utils.ShowResult()
		case <-redraw:
			utils.Move(utils.X, utils.Y)
		case <-quit:
			return
		}

	}

}
