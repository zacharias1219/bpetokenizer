package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Tokenizer struct {
	Vocab  map[string]int
	Merges map[string]string
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		Vocab:  make(map[string]int),
		Merges: make(map[string]string),
	}
}

func (t *Tokenizer) Train(text string, vocabSize int) {
	// Tokenization logic goes here
	words := strings.Split(text, " ")
	for i, word := range words {
		t.Vocab[word] = i + 1
	}
	fmt.Println("Training completed.")
}

func (t *Tokenizer) Encode(text string) []int {
	tokens := strings.Split(text, " ")
	var ids []int
	for _, token := range tokens {
		if id, exists := t.Vocab[token]; exists {
			ids = append(ids, id)
		} else {
			ids = append(ids, 0) // Unknown token
		}
	}
	return ids
}

func (t *Tokenizer) Decode(ids []int) string {
	var words []string
	for _, id := range ids {
		for word, wordID := range t.Vocab {
			if wordID == id {
				words = append(words, word)
				break
			}
		}
	}
	return strings.Join(words, " ")
}

func (t *Tokenizer) Save(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (t *Tokenizer) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, t)
}

func main() {
	tokenizer := NewTokenizer()
	text := "Hello world this is a test"
	tokenizer.Train(text, 100)
	encoded := tokenizer.Encode("Hello world")
	fmt.Println("Encoded:", encoded)
	fmt.Println("Decoded:", tokenizer.Decode(encoded))
}
