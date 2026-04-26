package main

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"go-blog/components"

	"github.com/a-h/templ"
)

func main() {
	// public/static ディレクトリを作成(cloudflare pagesはpublicディレクトリをルートとしているため)
	if err := os.MkdirAll("public/static", 0o755); err != nil {
		panic(err)
	}

	if err := copyDir("static", "public/static"); err != nil {
		panic(err)
	}

	// index.htmlにブログの内容をレンダリングする
	render("index.html", components.Home(
		"65bansekiのGo Blog",
	))

	render("post/index.html", components.Post())

}

func render(path string, c templ.Component) error {
	prefix := "public/"
	withPrefixPath := prefix + path
	// フォルダ作成
	os.MkdirAll(filepath.Dir(withPrefixPath), 0o755)
	// ファイル作成
	f, _ := os.Create(withPrefixPath)
	defer f.Close()
	ctx := context.Background()

	c.Render(ctx, f)

	return nil
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
