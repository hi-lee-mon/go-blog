# templ ガイド (v0.3.x)

templ は Go 用の HTML テンプレート言語です。`.templ` ファイルを書いて `templ generate` を実行すると、型安全な Go コードが生成されます。

---

## 1. セットアップ

```bash
# CLI のインストール
go install github.com/a-h/templ/cmd/templ@latest

# go.mod に追加
go get github.com/a-h/templ
```

---

## 2. ファイル構成

```
components/
  layout.templ       <- ソース（こちらを編集する）
  layout_templ.go    <- 自動生成（編集しない）
```

`templ generate`（またはウォッチモードの `templ generate --watch`）を実行すると、`.templ` から `.go` ファイルが再生成されます。

---

## 3. 基本的なコンポーネント

```templ
package components

templ Hello(name string) {
  <p>こんにちは、{ name }！</p>
}
```

- ファイルの先頭に `package <名前>` が必要です。
- コンポーネント名は PascalCase の Go 関数として定義します。
- `{ expr }` で Go の式を埋め込めます（HTML エスケープは自動）。

---

## 4. Children（スロット）

親コンポーネント内で `{ children... }` をスロットとして使います。

```templ
// layout.templ
templ Layout(title string) {
  <html>
    <body>
      <h1>{ title }</h1>
      { children... }
    </body>
  </html>
}
```

子コンポーネントを渡すには `@コンポーネント名 { }` ブロック構文を使います。

```templ
// home.templ
templ Home(title string) {
  @Layout(title) {
    <article>Home のコンテンツ</article>
  }
}
```

このプロジェクトでは `components/layout.templ` と `components/home.templ` でこのパターンを使用しています。

---

## 5. 別コンポーネントの呼び出し

```templ
templ Page() {
  @Header()
  <main>コンテンツ</main>
  @Footer()
}
```

---

## 6. If / Else

```templ
templ Greeting(isLoggedIn bool) {
  if isLoggedIn {
    <p>おかえりなさい！</p>
  } else {
    <a href="/login">ログイン</a>
  }
}
```

---

## 7. ループ

```templ
templ PostList(posts []Post) {
  <ul>
    for _, p := range posts {
      <li>{ p.Title }</li>
    }
  </ul>
  if len(posts) == 0 {
    <p>記事が見つかりませんでした。</p>
  }
}
```

---

## 8. 動的属性

```templ
templ Button(id string, label string) {
  <button id={ id }>{ label }</button>
}
```

`class` を動的に設定する場合は、文字列リテラル・CSS コンポーネント・`templ.KV` を組み合わせます。

```templ
templ Card(isActive bool) {
  <div class={ "card", templ.KV("card--active", isActive) }>
    { children... }
  </div>
}
```

---

## 9. CSS コンポーネント

スコープ付き CSS をインラインで定義できます。クラス名は衝突を避けるためハッシュ化されます。

```templ
package main

css loading(percent int) {
  width: { fmt.Sprintf("%d%%", percent) };
}

templ ProgressBar(percent int) {
  <div class={ loading(percent) }></div>
}
```

---

## 10. HTTP ハンドラでのレンダリング

```go
// main.go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  components.Home("マイブログ").Render(r.Context(), w)
})
```

`templ.Handler` を使うとより簡潔に書けます。

```go
http.Handle("/", templ.Handler(components.Home("マイブログ")))
```

---

## 11. クイックリファレンス

| 構文 | 用途 |
|---|---|
| `templ Name(args) { }` | コンポーネントの定義 |
| `{ expr }` | Go 式の埋め込み |
| `{ children... }` | 呼び出し元から渡された子要素のレンダリング |
| `@OtherComponent(args)` | 別コンポーネントの呼び出し（子なし） |
| `@OtherComponent(args) { }` | 子要素付きで呼び出し |
| `if / else` | 条件分岐 |
| `for range` | ループ |
| `css ClassName() { }` | スコープ付き CSS コンポーネント |
| `templ.KV(class, bool)` | 条件付きクラスの適用 |

---

## 12. よく使うコマンド

```bash
templ generate          # 一度だけ生成
templ generate --watch  # ウォッチモードで生成
templ fmt .             # 全 .templ ファイルをフォーマット
templ lsp               # 言語サーバーを起動（エディタ連携用）
```
