.POSIX:
.PHONY: all lint clean
.SUFFIXES:

all: dist/dicebot

dist/dicebot:
	go build -o $@ .

lint:
	golangci-lint run .

clean:
	rm -rf dist tmp
