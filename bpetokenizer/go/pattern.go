package bpe

import (
    "github.com/dlclark/regexp2"
)

// Final optimized pattern matching Python's behavior
const GPT4SplitPattern = `((?:'\w+|\w+'\w+|\w+)(?:[sdmt]|ll|ve|re)?)|\p{N}{1,3}|([^\s\p{L}\p{N}]+)|\s+`

func CompilePattern(pattern string) *regexp2.Regexp {
    return regexp2.MustCompile(pattern, regexp2.IgnoreCase)
}