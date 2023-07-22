package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func isLegal(word []byte, side map[byte]int) bool {
	if len(word) < 3 {
		return false
	}
	for i, ch := range word {
		if s, ok := side[ch]; !ok || (i > 0 && s == side[word[i-1]]) {
			return false
		}
	}
	return true
}

func removeChars(need, word []byte) []byte {
	var result []byte
	for _, ch := range need {
		if !bytes.ContainsAny(word, string([]byte{ch})) {
			result = append(result, ch)
		}
	}
	return result
}

func solve(vocab [][]byte, start byte, need []byte, left int) [][]string {
	if len(need) == 0 {
		return [][]string{nil}
	}
	if left == 0 {
		return nil
	}
	var result [][]string
	for _, w := range vocab {
		if (start == 0 || start == w[0]) && bytes.ContainsAny(w, string(need)) {
			for _, r := range solve(vocab, w[len(w)-1], removeChars(need, w), left-1) {
				result = append(result, append(r, string(w)))
			}
		}
	}
	return result
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "one arg please\n")
		os.Exit(1)
	}

	box := []byte(os.Args[1])
	side := make(map[byte]int)
	for i, ch := range box {
		side[ch] = i / 3
	}

	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var vocab [][]byte
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		w := []byte(scanner.Text())
		if isLegal(w, side) {
			vocab = append(vocab, w)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for limit := 1; limit <= 4; limit++ {
		solutions := solve(vocab, byte(0), box, limit)
		for _, s := range solutions {
			for i := len(s) - 1; i >= 0; i-- {
				fmt.Printf("%s ", s[i])
			}
			fmt.Printf("\n")
		}
		if len(solutions) > 0 {
			break
		}
	}
}
