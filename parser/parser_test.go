package parser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	input := "# タイトル\nこれは **太字** です\n- item1\n- item2"
	got := Parse(input)
	// got を出力して確認してみる
	t.Log(got)
	fmt.Println(got)
}
