.PHONY: generate dev build

dev:
	templ generate --watch & air

generate:
	templ generate

build:
	go run main.go