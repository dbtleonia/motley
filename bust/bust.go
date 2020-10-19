package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	f, err := os.Create("out.csv")
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)

	var prevHash string
	const limit = 3760796
	for i := 0; i < limit; i++ {
		if i%100000 == 0 {
			log.Printf("Iteration %d...\n", i)
		}
		var hash string
		if prevHash != "" {
			sum := sha256.Sum256([]byte(prevHash))
			hash = hex.EncodeToString(sum[:])
		} else {
			hash = "071697d02522ea4136f41742433f644680c05bfd852e72bf8af503510eeb28c3"
		}
		seed, err := hex.DecodeString(hash)
		if err != nil {
			log.Fatal(err)
		}
		bust, err := gameResult(seed, "0000000000000000004d6ec16dafe9d8370958664c1dc422f452892264c59526")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%d,%.2f,%s\n", i, bust, hash)
		prevHash = hash
	}

	w.Flush()
}

func gameResult(seed []byte, salt string) (float64, error) {
	const nBits = 52

	// 1.
	mac := hmac.New(sha256.New, []byte(salt))
	mac.Write(seed)
	strseed := hex.EncodeToString(mac.Sum(nil))

	// 2.
	strseed = strseed[:nBits/4]
	r, err := strconv.ParseInt(strseed, 16, 64)
	if err != nil {
		return 0.0, err
	}

	// 3.
	x := float64(r) / (1 << nBits)

	// 4.
	x = 99 / (1 - x)

	// 5.
	result := math.Floor(x)
	return math.Max(1, result/100), nil
}
