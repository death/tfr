package main

import (
	"log"
	"os"

	"github.com/death/tfr/db"
	"github.com/nsf/termbox-go"
)

const (
	textfilesDir    = "/media/1984/Documents/textfiles/"
	textfilesDBFile = "/media/1984/Documents/textfiles/tfr.db"
)

var currentView View

func main() {
	if _, err := os.Stat(textfilesDBFile); os.IsNotExist(err) {
		if err := db.Build(textfilesDir, textfilesDBFile); err != nil {
			log.Fatal(err)
		}
	}

	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	currentView = NewSectionsView()

	draw()
	eventLoop()
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	currentView.Draw()

	termbox.Flush()
}

func eventLoop() {
	for {
		ev := termbox.PollEvent()

		if ev.Type != termbox.EventKey {
			continue
		}

		if ev.Key == termbox.KeyEsc {
			return
		}

		if currentView.HandleEvent(ev) {
			draw()
		}
	}
}
