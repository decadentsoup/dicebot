.POSIX:
.PHONY: all lint clean
.SUFFIXES:

all: dist/dicebot

dist/dicebot:
	go build -o $@ .

lint:
	cp go.mod go.mod.tmp
	cp go.sum go.sum.tmp
	go mod tidy
	if ! diff go.mod go.mod.tmp; then echo "go.mod is not up to date; please run \"make tidy\"" 2>&1; exit 1; fi
	if ! diff go.sum go.sum.tmp; then echo "go.sum is not up to date; please run \"make tidy\"" 2>&1; exit 1; fi
	rm go.mod.tmp go.sum.tmp

	golangci-lint run .

clean:
	rm -rf dist tmp
