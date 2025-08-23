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
	fmt.Print("\033[?1049h")
	fmt.Print("\x1b[?25l")
	utils.MakeRandomCoins()
	if err != nil {
		log.Fatal(err)
	}
	quit, redraw := make(chan bool), make(chan bool, 1)
	go func() {
		for range utils.Timer {
			time.Sleep(time.Second)
			utils.Timer--
			redraw <- true
		}
		quit <- true
	}()

	go func() {
		input := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(input)
			if err != nil || n == 0 {
				continue
			}
			if input[0] == 'q' {
				quit <- true
			}

			if utils.SetPos(input[0]) {
				// redraw <- true
				utils.Draw()
			}
		}
	}()
	utils.Draw()
	for {
		select {
		case <-quit:
			return
		case <-redraw:
			utils.Draw()
		}
	}

}
