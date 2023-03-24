package main

import (
	typeGopher "gopher_typer"
	"log"
	"math/rand"
)

var (
	r *rand.Rand
)

func init() {
	// Seed the random number generator with a fixed seed
	r = rand.New(rand.NewSource(179876954235))
}
func main() {
	gt, err := typeGopher.NewGopherTyper()
	if err != nil {
		log.Fatal(err)
	}
	gt.Run()

}
