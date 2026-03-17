.PHONY: build run clean install uninstall status

VERSION ?= dev

build:
	go build -ldflags="-X main.Version=$(VERSION)" -o claude-usage-tray-go .

run: build
	./claude-usage-tray-go

clean:
	rm -f claude-usage-tray-go claude-usage-tray-go.exe

install: build
	./claude-usage-tray-go install

uninstall: build
	./claude-usage-tray-go uninstall

status: build
	./claude-usage-tray-go status
