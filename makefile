.PHONY: install
install:
	go mod tidy

.PHONY: test
test: 
	go test ./... -race -v

