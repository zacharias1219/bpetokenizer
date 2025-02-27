package bpe

import (
	"bytes"
	"sync"

	"github.com/dlclark/regexp2"
)

// Tokenizer defines the BPE tokenizer.
type Tokenizer struct {
	Vocab          map[int][]byte
	Merges         map[Pair]int
	SpecialTokens  map[string]int
	Pattern        *regexp2.Regexp
	InverseMerges  map[int]Pair
	InverseVocab   map[string]int
	InverseSpecial map[int]string
	mu             sync.RWMutex // protects shared maps during concurrent Encode calls
}

// Pair represents a pair of token IDs.
type Pair struct {
	A, B int
}

// NewTokenizer creates a new Tokenizer instance.
func NewTokenizer(specialTokens map[string]int) *Tokenizer {
	t := &Tokenizer{
		Vocab:          make(map[int][]byte),
		Merges:         make(map[Pair]int),
		SpecialTokens:  specialTokens,
		InverseMerges:  make(map[int]Pair),
		InverseVocab:   make(map[string]int),
		InverseSpecial: make(map[int]string),
	}

	// Initialize base vocabulary for IDs 0–255.
	for i := 0; i < 256; i++ {
		t.Vocab[i] = []byte{byte(i)}
		t.InverseVocab[string([]byte{byte(i)})] = i
	}

	// Initialize special tokens.
	for token, id := range specialTokens {
		t.Vocab[id] = []byte(token)
		t.InverseSpecial[id] = token
	}

	// Compile the default pattern (assumed GPT4SplitPattern is defined elsewhere).
	t.Pattern = CompilePattern(GPT4SplitPattern)
	return t
}

// Train performs the BPE training using concurrent frequency counting and merging.
func (t *Tokenizer) Train(text string, vocabSize int) {
	if vocabSize <= 256 {
		panic("vocab_size must be >= 256")
	}

	textChunks := t.SplitText(text)
	ids := make([][]int, len(textChunks))

	// Convert each text chunk to a slice of byte IDs.
	for i, chunk := range textChunks {
		ids[i] = bytesToIDs([]byte(chunk))
	}

	// BPE Algorithm: iteratively merge until the vocabulary reaches the desired size.
	for i := 256; i < vocabSize; i++ {
		// Concurrently compute frequency statistics for each chunk.
		statsCh := make(chan map[Pair]int, len(ids))
		var wg sync.WaitGroup
		for _, chunkIDs := range ids {
			wg.Add(1)
			go func(idsChunk []int) {
				defer wg.Done()
				statsCh <- getStats(idsChunk)
			}(chunkIDs)
		}
		wg.Wait()
		close(statsCh)

		// Aggregate statistics.
		stats := make(map[Pair]int)
		for m := range statsCh {
			for p, count := range m {
				stats[p] += count
			}
		}

		if len(stats) == 0 {
			break // no more pairs to merge.
		}

		bestPair := getMaxPair(stats)
		newID := i

		// Concurrently perform the merge on each chunk.
		newIds := make([][]int, len(ids))
		var mergeWg sync.WaitGroup
		for j, chunkIDs := range ids {
			mergeWg.Add(1)
			go func(index int, idsChunk []int) {
				defer mergeWg.Done()
				newIds[index] = merge(idsChunk, bestPair, newID)
			}(j, chunkIDs)
		}
		mergeWg.Wait()
		ids = newIds

		// Update shared maps.
		t.mu.Lock()
		t.Merges[bestPair] = newID
		t.InverseMerges[newID] = bestPair
		t.Vocab[newID] = append(t.Vocab[bestPair.A], t.Vocab[bestPair.B]...)
		t.InverseVocab[string(t.Vocab[newID])] = newID
		t.mu.Unlock()
	}
}

// Encode converts the input text into a slice of token IDs using concurrency.
func (t *Tokenizer) Encode(text string) []int {
	// Use a thread-safe text splitter.
	chunks := t.splitTextThreadSafe(text)
	result := make([][]int, len(chunks))

	// Copy shared maps under read-lock.
	t.mu.RLock()
	localMerges := t.Merges
	localSpecial := t.SpecialTokens
	t.mu.RUnlock()

	var wg sync.WaitGroup
	for i, chunk := range chunks {
		wg.Add(1)
		go func(index int, chunk string) {
			defer wg.Done()
			// If the chunk is a special token, use its ID.
			if specialID, exists := localSpecial[chunk]; exists {
				result[index] = []int{specialID}
				return
			}
			byteSeq := []byte(chunk)
			chunkIDs := bytesToIDs(byteSeq)
			// Iteratively merge tokens within the chunk.
			for len(chunkIDs) >= 2 {
				pairs := getStats(chunkIDs)
				if len(pairs) == 0 {
					break
				}
				var bestPair Pair
				minID := int(^uint(0) >> 1)
				for p := range pairs {
					if id, exists := localMerges[p]; exists && id < minID {
						minID = id
						bestPair = p
					}
				}
				if minID == int(^uint(0)>>1) {
					break
				}
				chunkIDs = merge(chunkIDs, bestPair, minID)
			}
			result[index] = chunkIDs
		}(i, chunk)
	}
	wg.Wait()

	// Concatenate results in order.
	var ids []int
	for _, chunkIDs := range result {
		ids = append(ids, chunkIDs...)
	}
	return ids
}

// Decode reconstructs the original text from token IDs.
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

// splitTextThreadSafe creates a new regex instance from the stored pattern and uses it to split text.
func (t *Tokenizer) splitTextThreadSafe(text string) []string {
	localPattern := CompilePattern(t.Pattern.String())
	var chunks []string
	m, _ := localPattern.FindStringMatch(text)
	for m != nil {
		chunks = append(chunks, m.String())
		m, _ = localPattern.FindNextMatch(m)
	}
	return chunks
}

// SplitText is the original non-thread-safe text splitting method.
func (t *Tokenizer) SplitText(text string) []string {
	var chunks []string
	m, _ := t.Pattern.FindStringMatch(text)
	for m != nil {
		chunks = append(chunks, m.String())
		m, _ = t.Pattern.FindNextMatch(m)
	}
	return chunks
}
