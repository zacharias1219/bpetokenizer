package bpe

import (
    "github.com/dlclark/regexp2"
)

// Revised pattern with proper group ordering
const GPT4SplitPattern = `((?:'(?:[sdmt]|ll|ve|re))|\w+|\p{N}{1,3}| ?[^\s\w\p{N}]+|\s*[\r\n]|\s+(?!\S)|\s+)`

func CompilePattern(pattern string) *regexp2.Regexp {
    return regexp2.MustCompile(pattern, regexp2.IgnoreCase)
}
