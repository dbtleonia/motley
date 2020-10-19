package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	letters = []byte("abcdefghijklmnopqrstuvwxyz")
)

func main() {
	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read all words.
	all := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		all[scanner.Text()] = true
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Remove some words.
	for w := range all {
		last := len(w) - 1
		if w[last] == 's' && all[w[:last]] {
			all[w] = false
		}
		if strings.Contains(w, "'") || strings.HasSuffix(w, "ed") || strings.HasSuffix(w, "ing") {
			all[w] = false
		}
	}

	words := make([]map[string]bool, 30)
	for i := range words {
		words[i] = make(map[string]bool)
	}
	for w, ok := range all {
		if ok {
			words[len(w)][strings.ToLower(w)] = true
		}
	}

	for n := 3; n <= 10; n++ {
		for i := 0; i < 5; i++ {
			bestWord, bestMatches := findBest(words[n])
			fmt.Printf("%s %d\n", bestWord, len(bestMatches))
			delete(words[n], bestWord)
			for _, w := range bestMatches {
				delete(words[n], w)
			}
		}
	}
}

func findBest(words map[string]bool) (string, []string) {
	var bestWord string
	var bestMatches []string
	for w := range words {
		var matches []string
		b := []byte(w)
		for i := range b {
			save := b[i]
			for _, l := range letters {
				if l == save {
					continue
				}
				b[i] = l
				s := string(b)
				if words[s] {
					matches = append(matches, s)
				}
			}
			b[i] = save
		}
		if len(matches) >= len(bestMatches) {
			bestWord = w
			bestMatches = matches
		}
	}
	return bestWord, bestMatches
}
