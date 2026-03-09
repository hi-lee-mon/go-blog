.PHONY: generate dev

dev:
	templ generate --watch & air

generate:
	templ generate

build:
	go run build.go