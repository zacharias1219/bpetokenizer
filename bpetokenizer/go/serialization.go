package bpe

import (
	"encoding/json"
	"fmt"
	"os"
)

type serializableTokenizer struct {
	Vocab         map[int]string
	Merges        map[string]int
	SpecialTokens map[string]int
	Pattern       string
}

func (t *Tokenizer) Save(filename string) error {
	data := serializableTokenizer{
		Vocab:         make(map[int]string),
		Merges:        make(map[string]int),
		SpecialTokens: t.SpecialTokens,
		Pattern:       t.Pattern.String(),
	}

	for id, bytes := range t.Vocab {
		data.Vocab[id] = replaceControlCharacters(string(bytes))
	}

	for pair, id := range t.Merges {
		key := fmt.Sprintf("%d,%d", pair.A, pair.B)
		data.Merges[key] = id
	}

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, file, 0644)
}

func (t *Tokenizer) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var data serializableTokenizer
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}

	t.Vocab = make(map[int][]byte)
	t.Merges = make(map[Pair]int)
	t.SpecialTokens = data.SpecialTokens
	t.Pattern = CompilePattern(data.Pattern)

	for id, str := range data.Vocab {
		t.Vocab[id] = []byte(str)
		t.InverseVocab[str] = id
	}

	for key, id := range data.Merges {
		var a, b int
		fmt.Sscanf(key, "%d,%d", &a, &b)
		pair := Pair{a, b}
		t.Merges[pair] = id
		t.InverseMerges[id] = pair
	}

	return nil
}
