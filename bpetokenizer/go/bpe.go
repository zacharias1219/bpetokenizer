package bpe

import (
	"bytes"

	"github.com/dlclark/regexp2"
)

type Tokenizer struct {
	Vocab          map[int][]byte
	Merges         map[Pair]int
	SpecialTokens  map[string]int
	Pattern        *regexp2.Regexp
	InverseMerges  map[int]Pair
	InverseVocab   map[string]int
	InverseSpecial map[int]string
}

type Pair struct {
	A, B int
}

func NewTokenizer(specialTokens map[string]int) *Tokenizer {
	t := &Tokenizer{
		Vocab:          make(map[int][]byte),
		Merges:         make(map[Pair]int),
		SpecialTokens:  specialTokens,
		InverseMerges:  make(map[int]Pair),
		InverseVocab:   make(map[string]int),
		InverseSpecial: make(map[int]string),
	}

	// Initialize base vocabulary
	for i := 0; i < 256; i++ {
		t.Vocab[i] = []byte{byte(i)}
		t.InverseVocab[string([]byte{byte(i)})] = i
	}

	// Initialize special tokens
	for token, id := range specialTokens {
		t.Vocab[id] = []byte(token)
		t.InverseSpecial[id] = token
	}

	t.Pattern = CompilePattern(GPT4SplitPattern)
	return t
}

func (t *Tokenizer) Train(text string, vocabSize int) {
	if vocabSize <= 256 {
		return
	}

	textChunks := t.SplitText(text)
	ids := make([][]int, len(textChunks))

	// Convert text chunks to byte IDs
	for i, chunk := range textChunks {
		ids[i] = bytesToIDs([]byte(chunk))
	}

	// BPE Algorithm
	for i := 256; i < vocabSize; i++ {
		stats := make(map[Pair]int)
		for _, chunkIDs := range ids {
			chunkStats := getStats(chunkIDs)
			for p, count := range chunkStats {
				stats[p] += count
			}
		}

		if len(stats) == 0 {
			break
		}

		bestPair := getMaxPair(stats)
		newID := i

		// Perform merge
		for j := range ids {
			ids[j] = merge(ids[j], bestPair, newID)
		}

		// Update vocab and merges
		t.Merges[bestPair] = newID
		t.InverseMerges[newID] = bestPair
		t.Vocab[newID] = append(t.Vocab[bestPair.A], t.Vocab[bestPair.B]...)
		t.InverseVocab[string(t.Vocab[newID])] = newID
	}
}

func (t *Tokenizer) Encode(text string) []int {
	var ids []int
	chunks := t.SplitText(text)

	for _, chunk := range chunks {
		if specialID, exists := t.SpecialTokens[chunk]; exists {
			ids = append(ids, specialID)
			continue
		}

		byteSeq := []byte(chunk)
		chunkIDs := bytesToIDs(byteSeq)

		for len(chunkIDs) >= 2 {
			pairs := getStats(chunkIDs)
			if len(pairs) == 0 {
				break
			}

			var bestPair Pair
			var minID = int(^uint(0) >> 1)
			for p := range pairs {
				if id, exists := t.Merges[p]; exists && id < minID {
					minID = id
					bestPair = p
				}
			}

			if minID == int(^uint(0)>>1) {
				break
			}

			chunkIDs = merge(chunkIDs, bestPair, minID)
		}
		ids = append(ids, chunkIDs...)
	}
	return ids
}

func (t *Tokenizer) Decode(ids []int) string {
	var buffer bytes.Buffer
	for _, id := range ids {
		if token, exists := t.InverseSpecial[id]; exists {
			buffer.WriteString(token)
			continue
		}
		buffer.Write(t.Vocab[id])
	}
	return buffer.String()
}

func (t *Tokenizer) SplitText(text string) []string {
	var chunks []string
	m, _ := t.Pattern.FindStringMatch(text)
	for m != nil {
		chunks = append(chunks, m.String())
		m, _ = t.Pattern.FindNextMatch(m)
	}
	return chunks
}
