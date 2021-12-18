package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type wordset struct {
	isPangram bool
	words     []string
}

func (ws *wordset) points() int {
	points := 0
	for _, w := range ws.words {
		if len(w) == 4 {
			points++
		} else {
			points += len(w)
		}
		if ws.isPangram {
			points += 7
		}
	}
	return points
}

func mask(word string) (wmask uint32, uniq int) {
	for i := 0; i < len(word); i++ {
		c := word[i]
		if c < 'a' || c > 'z' {
			return 0, 0
		}
		lmask := uint32(1) << (c - 'a')
		if wmask&lmask == 0 {
			uniq++
		}
		wmask |= lmask
	}
	return wmask, uniq
}

func main() {
	var centers [26]map[uint32]*wordset
	for i := 0; i < 26; i++ {
		centers[i] = make(map[uint32]*wordset)
	}

	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) < 4 {
			continue
		}
		wmask, uniq := mask(word)
		if uniq == 0 || uniq > 7 {
			continue
		}
		for i := uint32(0); i < 26; i++ {
			lmask := uint32(1) << i
			if wmask&lmask > 0 {
				if centers[i][wmask] == nil {
					centers[i][wmask] = &wordset{isPangram: uniq == 7}
				}
				centers[i][wmask].words = append(centers[i][wmask].words, word)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for center, wordsets := range centers {
		for wmask, wset := range wordsets {
			if wset.isPangram {
				var totalCount, totalPoints int
				for vmask, vset := range wordsets {
					if vmask|wmask == wmask {
						totalCount += len(vset.words)
						totalPoints += vset.points()
					}
				}
				for i := uint32(0); i < 26; i++ {
					lmask := uint32(1) << i
					if wmask&lmask > 0 {
						fmt.Printf("%c", i+'a')
					}
				}
				fmt.Printf(" %c %2d %3d %4d %s\n", center+'a', len(wset.words), totalCount, totalPoints, wset.words)
			}
		}
	}
}
