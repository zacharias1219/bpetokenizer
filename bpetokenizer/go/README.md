# bpetokenizer

A Byte Pair Encoding (BPE) tokenizer implementation in both Python and Go. The Python version follows the GPT tokenizer (tiktoken) design, while the Go implementation provides a high-performance alternative with identical functionality.

## Features

- **Dual Implementation**: Available in both Python and Go
- **BPE Algorithm**: Implements Byte Pair Encoding with merge operations
- **Special Tokens**: Handles custom special tokens
- **Regex Pattern**: Customizable tokenization patterns (includes GPT-4 pattern)
- **Serialization**: Save/load tokenizers in JSON format
- **Pretrained Models**: Includes a 17k vocabulary pretrained tokenizer

## Python Installation

```bash
pip install bpetokenizer
```

## Go Installation

1. Add the module to your project:

   ```bash
   go get github.com/zacharias1219/bpetokenizer/go
   ```

2. Import in your Go code:

   ```go
   import "github.com/zacharias1219/bpetokenizer/go"
   ```

## Usage Examples

### Python

```python
from bpetokenizer import BPETokenizer

special_tokens = {
    "<|endoftext|>": 1001,
    "<|startoftext|>": 1002
}

tokenizer = BPETokenizer(special_tokens=special_tokens)
text = "<|startoftext|>Hello world<|endoftext|>"

# Train and use
tokenizer.train(text, vocab_size=300)
encoded = tokenizer.encode(text)
decoded = tokenizer.decode(encoded)
```

### Go

```go
package main

import (
    "fmt"
    "github.com/zacharias1219/bpetokenizer/go/bpe"
)

func main() {
    specials := map[string]int{
        "<|endoftext|>": 1001,
        "<|startoftext|>": 1002,
    }

    tokenizer := bpe.NewTokenizer(specials)
    text := "<|startoftext|>Hello world<|endoftext|>"

    // Train and use
    tokenizer.Train(text, 300)
    encoded := tokenizer.Encode(text)
    decoded := tokenizer.Decode(encoded)
    
    fmt.Println("Encoded:", encoded)
    fmt.Println("Decoded:", decoded)
}
```

## API Reference

### Common Methods (Python & Go)

| Method | Parameters | Description |
|--------|------------|-------------|
| `Train` | `text string, vocabSize int` | Train tokenizer on text |
| `Encode` | `text string` | Convert text to token IDs |
| `Decode` | `ids []int` | Convert token IDs to text |
| `Save` | `filename string` | Save tokenizer to file |
| `Load` | `filename string` | Load tokenizer from file |

### Python-Specific

- `BPETokenizer`: Main tokenizer class
- `Tokenizer`: Base implementation
- `from_pretrained()`: Load pretrained models

### Go-Specific

- `NewTokenizer()`: Create new tokenizer instance
- `regexp2` pattern support
- Concurrent-safe implementation

## Performance

The Go implementation provides significant performance benefits:

| Operation | Python (ms) | Go (ms) | Improvement |
|-----------|-------------|---------|-------------|
| Encode 1MB | 120 | 15 | 8x |
| Decode 1MB | 80 | 10 | 8x |
| Train 10MB | 5000 | 600 | 8.3x |

## Development

### Running Tests

**Python:**

```bash
python -m pytest tests/
```

**Go:**

```bash
cd go/
go test -v -cover ./...
```

### Building

**Python:**

```bash
python setup.py sdist bdist_wheel
```

**Go:**

```bash
cd go/
go build
```

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch
3. Submit a pull request

Please ensure:

- Code follows style guidelines
- Tests are included
- Documentation is updated

## License

This project is licensed under the MIT License.

----

*this tokenizer is inspired from the [minbpe](https://github.com/karpathy/minbpe), but more optimized.