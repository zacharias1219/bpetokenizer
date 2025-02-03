package bpe

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBasicTokenization(t *testing.T) {
	specials := map[string]int{
		"<|endoftext|>": 1001,
	}
	tokenizer := NewTokenizer(specials)
	text := "aaabdaaabac"

	tokenizer.Train(text, 259)
	
	assert.Equal(t, 259, len(tokenizer.Vocab), "Vocabulary size mismatch")
	
	encoded := tokenizer.Encode("aaabdaaabac")
	assert.Equal(t, []int{258, 100, 258, 97, 99}, encoded, "Encoding mismatch")
	
	decoded := tokenizer.Decode(encoded)
	assert.Equal(t, "aaabdaaabac", decoded, "Decoding mismatch")
}

func TestSpecialTokens(t *testing.T) {
	specials := map[string]int{
		"<|startoftext|>": 1001,
		"<|endoftext|>":   1002,
	}
	
	tokenizer := NewTokenizer(specials)
	text := "<|startoftext|>Hello World<|endoftext|>"
	
	tokenizer.Train(text, 300)
	
	encoded := tokenizer.Encode(text)
	assert.Contains(t, encoded, 1001, "Missing start token")
	assert.Contains(t, encoded, 1002, "Missing end token")
	
	decoded := tokenizer.Decode(encoded)
	assert.Equal(t, text, decoded, "Special token decoding mismatch")
}

func TestSaveLoad(t *testing.T) {
	tokenizer := NewTokenizer(nil)
	text := "The quick brown fox jumps over the lazy dog"
	tokenizer.Train(text, 300)
	
	err := tokenizer.Save("test_model.json")
	assert.Nil(t, err, "Saving failed")
	
	newTokenizer := NewTokenizer(nil)
	err = newTokenizer.Load("test_model.json")
	assert.Nil(t, err, "Loading failed")
	
	assert.Equal(t, tokenizer.Vocab, newTokenizer.Vocab, "Vocab mismatch after load")
	assert.Equal(t, tokenizer.Merges, newTokenizer.Merges, "Merges mismatch after load")
}