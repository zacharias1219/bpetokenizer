package bpe

import (
	"regexp"
)

const GPT4SplitPattern = `'(?i:[sdmt]|ll|ve|re)|[^\r\n\p{L}\p{N}]?+\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]++[\r\n]*|\s*[\r\n]|\s+(?!\S)|\s+`

func CompilePattern(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}