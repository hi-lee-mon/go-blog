package main

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"go-blog/components"
)

func main() {
	if err := os.MkdirAll("public/static", 0o755); err != nil {
		panic(err)
	}

	if err := copyDir("static", "public/static"); err != nil {
		panic(err)
	}

	f, _ := os.Create("public/index.html")
	defer f.Close()
	ctx := context.Background()

	components.Post(
		"Goでブログを作ってみた",
		"Goでブログを作ってみました。Goはとても速いです。",
	).Render(ctx, f)
}

func copyDir(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
