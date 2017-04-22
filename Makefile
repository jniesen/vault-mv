default: test

test:
	go test || exit 1

PHONY: test
