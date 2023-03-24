package typeGopher

import (
	"os"

	tl "github.com/JoelOtter/termloop"
)

// GopherTyper handles the local state of the game.
type GopherTyper struct {
	g        *tl.Game
	wordList []string
	intro    introLevel
	game     gameLevel
	store    storeLevel
	end      endLevel
	console  tl.Text
	level    tl.Level
	stats    stats
	items    []item
}

// NewGopherTyper gets the game ready to run.
func NewGopherTyper() (*GopherTyper, error) {
	wReader, err := os.Open("data/words.txt")
	if err != nil {
		return nil, err
	}

	gt := GopherTyper{}
	gt.g = tl.NewGame()
	gt.g.Screen().SetFps(30)
	gt.wordList = newWordLoader(wReader)
	gt.intro = newIntroLevel(&gt, tl.ColorBlack, tl.ColorBlue)
	gt.game = newGameLevel(&gt, tl.ColorBlack, tl.ColorRed)
	gt.store = newStoreLevel(&gt, tl.ColorBlack, tl.ColorCyan)
	gt.end = newEndLevel(&gt, tl.ColorBlack, tl.ColorGreen)

	gt.stats = newStats()

	return &gt, nil
}

// Run starts the game, and blocks forever.
func (gt *GopherTyper) Run() {
	gt.goToIntro()
	gt.g.Start()
}

// goToIntro sets the current level to the intro and activates it.
func (gt *GopherTyper) goToIntro() {
	gt.level = &gt.intro
	gt.intro.Activate()
}

// goToGame sets the current level to the game level and activates it.
func (gt *GopherTyper) goToGame() {
	if gt.stats.Lives == 0 {
		gt.stats = newStats()
		gt.items = []item{}
	}
	gt.level = &gt.game
	gt.game.Activate()
}

// goToStore sets the current level to the store level and activates it.
func (gt *GopherTyper) goToStore() {
	gt.level = &gt.store
	gt.store.Activate()
}

// goToEndWin sets the current level to the end level with a win condition and activates it.
func (gt *GopherTyper) goToEndWin() {
	gt.level = &gt.end
	gt.end.ActivateWin()
}

// goToEndFail sets the current level to the end level with a fail condition and activates it.
func (gt *GopherTyper) goToEndFail() {
	gt.level = &gt.end
	gt.end.ActivateFail()
}
