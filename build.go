package main

import (
	"context"
	"os"

	"go-blog/components"
)

func main() {
	f, _ := os.Create("public/index.html")
	defer f.Close()
	ctx := context.Background()

	components.Post(
		"Goでブログを作ってみた",
		"Goでブログを作ってみました。Goはとても速いです。",
	).Render(ctx, f)
}
