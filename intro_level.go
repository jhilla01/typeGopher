package typeGopher

import (
	"os"
	"time"

	tl "github.com/JoelOtter/termloop"
)

type introLevel struct {
	tl.Level
	gt              *GopherTyper
	pressAKeyText   *tl.Text
	needsRefresh    bool
	swapMessageTime time.Time
	reverseText     bool
}

// Activate sets the intro level as the current level and marks it for refresh.
func (l *introLevel) Activate() {
	l.needsRefresh = true
	l.gt.g.Screen().SetLevel(l)
}

// refresh updates the intro level's display, adding the necessary entities and text.
func (l *introLevel) refresh() {
	l.gt.intro.AddEntity(&l.gt.console)
	l.gt.console.SetText("")
	w, h := l.gt.g.Screen().Size()
	quarterH := h / 4
	rect := tl.NewRectangle(10, 2, w-20, h-4, tl.ColorCyan)
	l.AddEntity(rect)

	logo, _ := os.ReadFile("data/logo.txt")
	c := tl.CanvasFromString(string(logo))
	logoEntity := tl.NewEntityFromCanvas(w/2-len(c)/2, quarterH, tl.CanvasFromString(string(logo)))
	l.AddEntity(logoEntity)

	msg := "Press any key to continue"
	l.pressAKeyText = tl.NewText(w/2-len(msg)/2, h/2, msg, tl.ColorBlue|tl.AttrReverse, tl.ColorDefault)
	l.AddEntity(l.pressAKeyText)

	instructions, _ := os.ReadFile("data/instructions.txt")
	c = tl.CanvasFromString(string(instructions))
	l.AddEntity(tl.NewEntityFromCanvas(w/2-len(c)/2, h/2+2, c))

	l.needsRefresh = false
}

// Draw refreshes the intro level's display if needed and updates the "Press any key" text's appearance.
func (l *introLevel) Draw(screen *tl.Screen) {
	if l.needsRefresh {
		l.refresh()
	}
	if time.Now().After(l.swapMessageTime) {
		if l.reverseText {
			l.pressAKeyText.SetColor(tl.ColorBlue, tl.ColorDefault)
		} else {
			l.pressAKeyText.SetColor(tl.ColorBlue|tl.AttrReverse, tl.ColorDefault)
		}
		l.reverseText = !l.reverseText
		l.swapMessageTime = time.Now().Add(500 * time.Millisecond)
	}
	l.Level.Draw(screen)
}

// Tick handles user input, transitioning to the game level when a key is pressed.
func (l *introLevel) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		l.gt.goToGame()
	}
}

// newIntroLevel creates a new intro level with the given GopherTyper, foreground, and background attributes.
func newIntroLevel(g *GopherTyper, fg, bg tl.Attr) introLevel {
	l := tl.NewBaseLevel(tl.Cell{Bg: bg, Fg: fg})
	return introLevel{Level: l, gt: g}
}
