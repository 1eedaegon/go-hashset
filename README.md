# go-hashset

[![Go](https://pkg.go.dev/badge/github.com/1eedaegon/go-hashset.svg)](https://pkg.go.dev/github.com/1eedaegon/go-hashset)
[![CI](https://github.com/1eedaegon/go-hashset/actions/workflows/go.yml/badge.svg)](https://github.com/1eedaegon/go-hashset/actions/workflows/go.yml)
[![CodeQL](https://github.com/1eedaegon/go-hashset/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/1eedaegon/go-hashset/actions/workflows/codeql.yml)

A go library hashset for O(1)

## Example

Initialize set and Add element

```go
import (
	...
	hashset "github.com/1eedaegon/go-hashset"
	...
)

	s := hashset.New(1, 2, "3")
	s.Add("3") // Since "3" already exists in the Set, its size remains 3.
```

Remove and Distinguish between different types

```go
import (
	...
	hashset "github.com/1eedaegon/go-hashset"
	...
)
    s := hashset.New("1", "2", 3)
	s.Remove("1")
	s.Remove("multiple")
	s.Remove("3") // The length of s is 2, because due to the difference in types.
```

Other case doc: 

## License

[MIT](LICENSE)
