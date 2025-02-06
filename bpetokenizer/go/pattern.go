package bpe

import (
	"github.com/dlclark/regexp2"
)

const GPT4SplitPattern = `('(?:[sdmt]|ll|ve|re))|([^\r\n\p{L}\p{N}]?\p{L}+)|(\p{N}{1,3})|( ?[^\s\p{L}\p{N}]+[\r\n]*)|(\s*[\r\n])|(\s+(?!\S))|(\s+)`

func CompilePattern(pattern string) *regexp2.Regexp {
	return regexp2.MustCompile(pattern, regexp2.IgnoreCase)
}
