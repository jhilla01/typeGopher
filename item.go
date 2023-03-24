package typeGopher

import (
	"fmt"
	"math/rand"
	"time"
)

type item interface {
	Name() string
	Desc() string
	Price() int
	PriceDesc() string
	SetID(int)
	Reset(gt *GopherTyper)
	Purchase(g *storeLevel) bool
	Dupe() item

	Tick(g *gameLevel)
}

type goroutineItem struct {
	wakeAt      time.Time
	baseWait    time.Duration
	waitRange   time.Duration
	currentWord *word
	id          int
	cpuUpgrades int
	price       int
}

// Name returns the name of the goroutineItem.
func (i *goroutineItem) Name() string {
	return "Goroutine"
}

// Desc returns the description of the goroutineItem.
func (i *goroutineItem) Desc() string {
	return "Add a goroutine to help type words for you"
}

// Price returns the price of the goroutineItem.
func (i *goroutineItem) Price() int {
	return i.price
}

// PriceDesc returns the price of the goroutineItem as a formatted string.
func (i *goroutineItem) PriceDesc() string {
	return fmt.Sprintf("$%d", i.Price())
}

// Tick handles the logic for the goroutineItem during each game tick.
func (i *goroutineItem) Tick(gl *gameLevel) {
	if time.Now().After(i.wakeAt) {
		if i.currentWord == nil {
			var possibleWords []int
			for i, w := range gl.words {
				if gl.currentWord != gl.words[i] && !w.Complete() && w.startedBy == 0 {
					possibleWords = append(possibleWords, i)
				}
			}
			if len(possibleWords) > 0 {
				i.currentWord = gl.words[possibleWords[rand.Intn(len(possibleWords))]]
				i.currentWord.completedChars++
				i.currentWord.startedBy = i.id
			}
		} else {
			i.currentWord.completedChars++
			gl.gt.stats.Garbage++
			if i.currentWord.Complete() {
				i.currentWord = nil
			}
		}

		i.sleep()
	}
}

// sleep sets the wakeAt time for the goroutineItem.
func (i *goroutineItem) sleep() {
	i.wakeAt = time.Now().Add(i.baseWait/time.Duration(i.cpuUpgrades) + time.Duration(rand.Intn(int(i.waitRange))))
}

// SetID sets the ID for the goroutineItem.
func (i *goroutineItem) SetID(id int) {
	i.id = id
}

// Reset resets the state of the goroutineItem.
func (i *goroutineItem) Reset(gt *GopherTyper) {
	i.currentWord = nil
	i.cpuUpgrades = gt.stats.CPUUpgrades
	i.price = 1000
	for _, itm := range gt.items {
		if itm.Name() == i.Name() {
			i.price *= 2
		}
	}
}

// Dupe creates a duplicate of the goroutineItem.
func (i *goroutineItem) Dupe() item {
	var dupe goroutineItem
	dupe = *i
	return &dupe
}

// Purchase handles the purchasing logic for the goroutine item and returns true if the purchase is successful.
func (i *goroutineItem) Purchase(l *storeLevel) bool {
	return true
}

// newGoroutineItem creates a new goroutine item.
func newGoroutineItem(waitRange, baseWait time.Duration) *goroutineItem {
	item := goroutineItem{waitRange: waitRange, baseWait: baseWait, cpuUpgrades: 1}
	item.sleep()
	return &item
}

type cpuUpgradeItem struct {
	id    int
	price int
}

// Name returns the name of the cpuUpgradeItem.
func (i *cpuUpgradeItem) Name() string {
	return "CPU Upgrade"
}

// Desc returns the description of the cpuUpgradeItem.
func (i *cpuUpgradeItem) Desc() string {
	return "Make your goroutines go faster"
}

// Price returns the price of the cpuUpgradeItem.
func (i *cpuUpgradeItem) Price() int {
	return i.price
}

// PriceDesc returns the price of the cpuUpgradeItem as a formatted string.
func (i *cpuUpgradeItem) PriceDesc() string {
	return fmt.Sprintf("$%d", i.Price())
}

// Tick handles the logic for the cpuUpgradeItem during each game tick.
func (i *cpuUpgradeItem) Tick(gl *gameLevel) {
}

// SetID sets the ID for the cpuUpgradeItem.
func (i *cpuUpgradeItem) SetID(id int) {
	i.id = id
}

// Reset resets the state of the cpuUpgradeItem.
func (i *cpuUpgradeItem) Reset(gt *GopherTyper) {
	i.price = 2000 * gt.stats.CPUUpgrades
}

// Purchase handles the purchasing logic for the cpuUpgradeItem and returns false if the purchase is not successful.
func (i *cpuUpgradeItem) Purchase(l *storeLevel) bool {
	l.gt.stats.CPUUpgrades++
	return false
}

// Dupe creates a duplicate of the cpuUpgradeItem.
func (i *cpuUpgradeItem) Dupe() item {
	var dupe cpuUpgradeItem
	dupe = *i
	return &dupe
}

type goUpgradeItem struct {
	id    int
	price int
}

// Name returns the name of the goUpgradeItem.
func (i *goUpgradeItem) Name() string {
	return "Go Upgrade"
}

// Desc returns the description of the goUpgradeItem.
func (i *goUpgradeItem) Desc() string {
	return "Improves garbage collection performance"
}

// Price returns the price of the goUpgradeItem.

func (i *goUpgradeItem) Price() int {
	return i.price
}

// PriceDesc returns the price of the goUpgradeItem as a formatted string.
func (i *goUpgradeItem) PriceDesc() string {
	return fmt.Sprintf("$%d", i.Price())
}

// Tick handles the logic for the goUpgradeItem during each game tick.
func (i *goUpgradeItem) Tick(gl *gameLevel) {
}

// SetID sets the ID for the goUpgradeItem.
func (i *goUpgradeItem) SetID(id int) {
	i.id = id
}

// Reset resets the state of the goUpgradeItem.
func (i *goUpgradeItem) Reset(gt *GopherTyper) {
	i.price = int(1000 * gt.stats.GoVersion)
}

// Purchase handles the purchasing logic for the goUpgradeItem and returns false if the purchase is not successful.
func (i *goUpgradeItem) Purchase(l *storeLevel) bool {
	l.gt.stats.GoVersion += 0.1
	l.gt.stats.GarbageFreq += 3
	return false
}

// Dupe creates a duplicate of the goUpgradeItem.
func (i *goUpgradeItem) Dupe() item {
	var dupe goUpgradeItem
	dupe = *i
	return &dupe

}
