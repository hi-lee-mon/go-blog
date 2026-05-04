.PHONY: generate dev build

dev:
	templ generate --watch & air

# templからgoファイルを作成する
generate:
	templ generate

build:
	go run main.go