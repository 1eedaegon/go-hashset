# go-hashset

[![Go](https://pkg.go.dev/badge/github.com/1eedaegon/go-hashset.svg)](https://pkg.go.dev/github.com/1eedaegon/go-hashset)
[![CI](https://github.com/1eedaegon/go-hashset/actions/workflows/go.yml/badge.svg)](https://github.com/1eedaegon/go-hashset/actions/workflows/go.yml)
[![CodeQL](https://github.com/1eedaegon/go-hashset/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/1eedaegon/go-hashset/actions/workflows/codeql.yml)

A go library hashset for O(1)

## Example

```go
import (
	...
	port "github.com/1eedaegon/go-hashset"
	...
)

ports := port.Get(3)
// ports is something like []int{10000, 10001, 10002}
```

Or

```go
import (
	...
	port "github.com/1eedaegon/go-hashset"
	...
)

ports := port.GetS(3)
// ports is something like []string{"10000", "10001", "10002"}
```

## License

[MIT](LICENSE)
