.PHONY: build run clean

VERSION ?= dev

build:
	go build -ldflags="-X main.Version=$(VERSION)" -o claude-usage-tray-go .

run: build
	./claude-usage-tray-go

clean:
	rm -f claude-usage-tray-go claude-usage-tray-go.exe
