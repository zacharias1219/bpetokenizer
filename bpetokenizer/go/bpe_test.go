package bpe

import (
	"strings"
	"sync"
	"testing"

	"github.com/dlclark/regexp2"
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

	expectedChunks := []string{"Hello's", " ", "world", " ", "123"}

	chunks := tokenizer.SplitText(text)
	assert.Equal(t, expectedChunks, chunks, "Pattern splitting mismatch")
}

func TestInvalidVocabSize(t *testing.T) {
	tokenizer := NewTokenizer(nil)
	text := "valid text"

	t.Run("vocab too small", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for vocab_size < 256")
			}
		}()
		tokenizer.Train(text, 100)
	})
}

func TestConcurrentEncode(t *testing.T) {
	tokenizer := NewTokenizer(nil)
	text := "concurrency test"
	tokenizer.Train(text, 300)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			encoded := tokenizer.Encode(text)
			decoded := tokenizer.Decode(encoded)
			assert.Equal(t, text, decoded)
		}()
	}
	wg.Wait()
}

func TestFuzzRoundtrip(t *testing.T) {
	tokenizer := NewTokenizer(nil)
	texts := []string{
		"Hello world! 123",
		"日本語テスト",
		"😀🤖✨",
		"Special chars: ~!@#$%^&*()",
		"",
		"   ",
	}

	for _, text := range texts {
		encoded := tokenizer.Encode(text)
		decoded := tokenizer.Decode(encoded)
		assert.Equal(t, text, decoded)
	}
}

func TestBPETokenizerSpecialTokens(t *testing.T) {
	specials := map[string]int{
		"<|start|>": 1001,
		"<|end|>":   1002,
	}

	tokenizer := NewTokenizer(specials)
	text := "<|start|>Hello world<|end|>"

	t.Run("special token preservation", func(t *testing.T) {
		tokenizer.Train(text, 300)
		encoded := tokenizer.Encode(text)
		decoded := tokenizer.Decode(encoded)

		assert.Contains(t, decoded, "<|start|>")
		assert.Contains(t, decoded, "<|end|>")
		assert.Equal(t, text, decoded)
	})
}

func TestExactEncodedIDs(t *testing.T) {
	tokenizer := NewTokenizer(nil)
	text := "aaabdaaabac"
	tokenizer.Train(text, 259)

	expected := []int{258, 100, 258, 97, 99}
	assert.Equal(t, expected, tokenizer.Encode(text))
}

func TestPatternConfiguration(t *testing.T) {
    customPattern := `([a-zA-Z]+|\d+)`

    tokenizer := &Tokenizer{
        Pattern: regexp2.MustCompile(customPattern, regexp2.None),
    }

    text := "Hello123World"
    expected := []string{"Hello", "123", "World"}
    assert.Equal(t, expected, tokenizer.SplitText(text))
}

func TestLargeTextHandling(t *testing.T) {
	tokenizer := NewTokenizer(nil)
	text := strings.Repeat("The quick brown fox jumps over the lazy dog ", 1000)

	tokenizer.Train(text, 500)
	encoded := tokenizer.Encode(text)
	decoded := tokenizer.Decode(encoded)

	assert.Equal(t, text, decoded)
	assert.True(t, len(encoded) < len(text)/2, "Compression ratio check")
}
