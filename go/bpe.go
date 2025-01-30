// bpe.go
package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"sort"
)

// TrainBPE implements the BPE algorithm in Go
func TrainBPE(text string, vocabSize int) map[[2]int]int {
	tokens := []byte(text)
	merges := make(map[[2]int]int)
	vocab := make(map[int][]byte)

	// Initialize vocab with byte tokens
	for i := 0; i < 256; i++ {
		vocab[i] = []byte{byte(i)}
	}

	// BPE algorithm
	for i := 256; i < vocabSize; i++ {
		pair := findMostFrequentPair(tokens)
		if len(pair) == 0 {
			break
		}
		merges[pair] = i
		tokens = mergePair(tokens, pair, i)
		vocab[i] = append(vocab[pair[0]], vocab[pair[1]]...)
	}

	return merges
}

// Helper function to find the most frequent pair
func findMostFrequentPair(tokens []byte) [2]int {
	counts := make(map[[2]int]int)
	for i := 0; i < len(tokens)-1; i++ {
		pair := [2]int{int(tokens[i]), int(tokens[i+1])}
		counts[pair]++
	}

	if len(counts) == 0 {
		return [2]int{}
	}

	var maxPair [2]int
	maxCount := 0
	for pair, count := range counts {
		if count > maxCount {
			maxPair = pair
			maxCount = count
		}
	}

	return maxPair
}

// Helper function to merge a pair in the tokens
func mergePair(tokens []byte, pair [2]int, idx int) []byte {
	var newTokens []byte
	i := 0
	for i < len(tokens) {
		if i < len(tokens)-1 && int(tokens[i]) == pair[0] && int(tokens[i+1]) == pair[1] {
			newTokens = append(newTokens, byte(idx))
			i += 2
		} else {
			newTokens = append(newTokens, tokens[i])
			i++
		}
	}
	return newTokens
}

//export TrainBPEWrapper
func TrainBPEWrapper(text *C.char, vocabSize C.int) *C.char {
	goText := C.GoString(text)
	merges := TrainBPE(goText, int(vocabSize))

	// Convert merges to JSON string
	jsonStr, _ := json.Marshal(merges)
	return C.CString(string(jsonStr))
}

func main() {}
