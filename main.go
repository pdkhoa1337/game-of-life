package main

import (
	"bufio"
	"log"
	"os"
	"strconv"

	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

//	var matrix [][]int = [][]int{
//		{-1, -1}, {-1, 0}, {-1, 1},
//		{0, -1}, {0, 1},
//		{1, -1}, {1, 0}, {1, -1},
//	}
var matrix [][]int = [][]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

var xmax int
var ymax int
var drawFlag = false

// var neighborCounts [][]int8

// func initNC() {
// 	neighborCounts = make([][]int8, ymax)
// 	for i := range neighborCounts {
// 		neighborCounts[i] = make([]int8, xmax)
//     for j := range neighborCounts[i]{
//       neighborCounts[i][j] = 0
//     }
// 	}
// }

func writeFile(filename string, x [][]int8) error {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString("\n")
	if err != nil {
		return err
	}
	for i := range x {
		for j := range x[i] {
			_, err = w.WriteString(strconv.FormatInt(int64(x[i][j]), 10))
			if err != nil {
				return err
			}
		}
		_, err = w.WriteRune('\n')
		if err != nil {
			return err
		}
	}
	err = w.Flush()
	return err
}

func alive(living int8, x int, y int, game [][]int8) bool {
	var neighborCount int8 = 0
	for _, v := range matrix {
		var nx int = v[1] + x
		var ny int = v[0] + y
		if nx < (xmax-1) && ny < (ymax-1) && nx >= 0 && ny >= 0 {
			neighborCount += game[ny][nx]
		}
	}
	// neighborCounts[y][x] = neighborCount
	if living > 0 {
		return neighborCount == 2 || neighborCount == 3
	} else {
		return neighborCount == 3
	}
}

func main() {
	var interval int64 = 2

	var err error
	if len(os.Args) > 1 {
		interval, err = strconv.ParseInt(os.Args[1], 10, 64)
	}
	if err != nil {
		panic(err)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	cellStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	// Draw initial boxes

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	xmax, ymax = s.Size()
	game := make([][]int8, ymax)
	for i := range game {
		game[i] = make([]int8, xmax)
		for j := range game[i] {
			if rand.Int()%2 == 0 {
				game[i][j] = 1
				s.SetContent(j, i, tcell.RuneBlock, nil, cellStyle)
			}
		}
	}
	go func() {
		for {
			s.Show()
			time.Sleep(time.Duration(interval) * 100 * time.Millisecond)
			for drawFlag {
			}
			s.Clear()
			newGame := make([][]int8, ymax)
			for i := range game {
				newGame[i] = make([]int8, xmax)
			}
			for i, v := range game {
				for j, cell := range v {
					if alive(cell, j, i, game) {
						newGame[i][j] = 1
						s.SetContent(j, i, tcell.RuneBlock, nil, cellStyle)
					} else {
						newGame[i][j] = 0
					}
				}
			}
			game = newGame
		}
	}()
	// Here's an example of how to inject a keystroke where it will
	// be picked up by the next PollEvent call.  Note that the
	// queue is LIFO, it has a limited length, and PostEvent() can
	// return an error.
	// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))
	for {
		// Update screen

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				s.Clear()
				writeFile("game.txt", game)
				return
			} else if ev.Rune() >= '0' || ev.Rune() <= '9' {
        interval = int64(ev.Rune() - '0')
			}
		case *tcell.EventMouse:
			mX, mY := ev.Position()
			if drawFlag {
				game[mY][mX] = 1
				s.SetContent(mX, mY, tcell.RuneBlock, nil, cellStyle)
				s.Show()
			}
			switch ev.Buttons() {
			case tcell.Button1:
				drawFlag = true
			case tcell.ButtonNone:
				drawFlag = false
			}

		}
	}
}
