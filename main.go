package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go-blog/components"

	"go-blog/parser"

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

	blogDataList := []components.BlogData{}
	filepath.WalkDir("./contents", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリはスキップ
		if d.IsDir() {
			return nil
		}

		// mdのファイルからhtmlのファイル名とパスを計算
		replacedPath := strings.Replace(path, ".md", ".html", 1)
		filename := filepath.Base(replacedPath)

		// マークダウンファイルの中身をパース
		content, err := os.ReadFile(path)
		parsed := parser.Parse(string(content))

		// ブログページを生成
		render("/post/"+filename, components.Post(parsed))

		// ビルドしたブログへのリンクを作成
		href := "/public" + strings.TrimPrefix(replacedPath, "contents")

		data := components.BlogData{
			Title: filename, // TODO:フロントマター
			Href:  href,
		}
		blogDataList = append(blogDataList, data)
		return nil
	})

	render("index.html", components.Home("65bansekiのGo Blog", blogDataList))
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
