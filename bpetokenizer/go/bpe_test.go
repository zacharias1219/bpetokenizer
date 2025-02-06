package bpe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicTokenization(t *testing.T) {
    t.Run("simple merge operations", func(t *testing.T) {
        text := "aaabdaaabac"
        
        tokenizer := NewTokenizer(nil)
        tokenizer.Train(text, 259)
        
        encoded := tokenizer.Encode(text)
        decoded := tokenizer.Decode(encoded)
        
        assert.Equal(t, text, decoded)
        assert.Equal(t, []int{258, 100, 258, 97, 99}, encoded)
    })
}

func TestSpecialTokens(t *testing.T) {
	specials := map[string]int{
		"<|startoftext|>": 1001,
		"<|endoftext|>":   1002,
	}

	text := "<|startoftext|>Hello World<|endoftext|>"

	tokenizer := NewTokenizer(specials)
	tokenizer.Train(text, 300)

	encoded := tokenizer.Encode(text)
	decoded := tokenizer.Decode(encoded)

	assert.Contains(t, decoded, "<|startoftext|>", "Start token missing")
	assert.Contains(t, decoded, "<|endoftext|>", "End token missing")
	assert.Equal(t, text, decoded, "Decoding mismatch")
}

func TestPatternSplitting(t *testing.T) {
    tokenizer := NewTokenizer(nil)
    text := "Hello's world 123"
    
    expectedChunks := []string{"Hello's", " world", "123"}
    
    chunks := tokenizer.SplitText(text)
    assert.Equal(t, expectedChunks, chunks, "Pattern splitting mismatch")
}
