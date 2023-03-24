package typeGopher

import "math/rand"

type stats struct {
	LevelsCompleted int
	LevelsAttempted int
	Dollars         int
	TotalEarned     int
	CPUUpgrades     int
	GoVersion       float32
	Lives           int
	Garbage         int
	GarbageFreq     int
}

// newStats creates and returns a new "stats" object with default values.
func newStats() stats {
	return stats{Lives: 3, CPUUpgrades: 1, GarbageFreq: 10, GoVersion: 1.0}
}

// GarbageCollect if Garbage is > 0 generate a random number between 0 and current garbage
// If the randomly generated integer is greater than the "GarbageFreq" field of the "stats" struct, the function returns true
func (s *stats) GarbageCollect() bool {
	if s.Garbage > 0 && rand.Intn(s.Garbage) > s.GarbageFreq {
		return true
	}
	return false
}
