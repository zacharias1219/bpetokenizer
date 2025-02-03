package bpe

import (
	"bytes"
	"unicode/utf8"
	"fmt"
)

func bytesToIDs(b []byte) []int {
	ids := make([]int, len(b))
	for i, v := range b {
		ids[i] = int(v)
	}
	return ids
}

func getStats(ids []int) map[Pair]int {
	stats := make(map[Pair]int)
	for i := 0; i < len(ids)-1; i++ {
		p := Pair{ids[i], ids[i+1]}
		stats[p]++
	}
	return stats
}

func getMaxPair(stats map[Pair]int) Pair {
	var maxPair Pair
	maxCount := 0
	for p, count := range stats {
		if count > maxCount || (count == maxCount && p.A < maxPair.A) {
			maxPair = p
			maxCount = count
		}
	}
	return maxPair
}

func merge(ids []int, p Pair, newID int) []int {
	newIDs := make([]int, 0, len(ids))
	i := 0
	for i < len(ids) {
		if i < len(ids)-1 && ids[i] == p.A && ids[i+1] == p.B {
			newIDs = append(newIDs, newID)
			i += 2
		} else {
			newIDs = append(newIDs, ids[i])
			i++
		}
	}
	return newIDs
}

func replaceControlCharacters(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		if utf8.RuneLen(r) == -1 || (r < 0x20 && r != 0x0A && r != 0x0D) || r == 0x7F {
			buf.WriteString(fmt.Sprintf("\\u%04x", r))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}