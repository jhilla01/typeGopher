package typeGopher

import (
	"bufio"
	"io"
)

// newWordLoader reads words from an io.Reader and returns them as a slice of strings.
func newWordLoader(r io.Reader) []string {
	var words []string
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words
}
