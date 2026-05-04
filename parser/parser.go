package parser

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	prefixHeadingOne   = "# "
	prefixHeadingTwo   = "## "
	prefixHeadingThree = "### "
	prefixBulletPoint  = "- "
)

var regexpPrefixOrder = regexp.MustCompile(`^\d+\. `)

func replaceStrong(input string) string {
	var b strings.Builder
	splits := strings.Split(input, "**")
	for i, s := range splits {
		if i%2 != 0 {
			// 奇数
			content := fmt.Sprintf("<strong>%s</strong>", s)
			b.WriteString(content)
		} else {
			b.WriteString(s)
		}
	}
	return b.String()
}

func writeString(input string, b *strings.Builder) {
	bold := replaceStrong(input)
	b.WriteString(bold)
}

func Parse(input string) string {
	// HTMLビルダー
	var b strings.Builder

	// リスト要素 状態
	inBulletPoint := false
	var bulletPointBuilder strings.Builder
	inOrder := false
	var orderBuilder strings.Builder

	/*
		パース処理
	*/
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		/*
			状態毎の処理
		*/
		switch {
		case inBulletPoint:
			/*
				liビルド処理
			*/
			if strings.HasPrefix(line, prefixBulletPoint) {
				// 連続したli処理
				replacedLine := strings.Replace(line, prefixBulletPoint, "", 1)
				content := fmt.Sprintf("<li>%s</li>", replacedLine)
				bulletPointBuilder.WriteString(content)
				continue
			} else {
				// liの次にli以外が来た場合はulで閉じる
				contents := bulletPointBuilder.String()
				content := fmt.Sprintf("<ul>%s</ul>", contents)
				writeString(content, &b)
				// 状態リセット
				inBulletPoint = false
				bulletPointBuilder.Reset()
			}
		case inOrder:
			if regexpPrefixOrder.MatchString(line) {
				replacedLine := regexpPrefixOrder.ReplaceAllString(line, "")
				content := fmt.Sprintf("<li>%s</li>", replacedLine)
				orderBuilder.WriteString(content)
				continue
			} else {
				contents := orderBuilder.String()
				content := fmt.Sprintf("<ol>%s</ol>", contents)
				writeString(content, &b)
				// 状態リセット
				inOrder = false
				orderBuilder.Reset()
			}
		}

		/*
			状態なしの処理
		*/
		switch {
		case strings.HasPrefix(line, prefixHeadingOne):
			replacedLine := strings.Replace(line, prefixHeadingOne, "", 1)
			content := fmt.Sprintf("<h1>%s</h1>", replacedLine)
			writeString(content, &b)

		case strings.HasPrefix(line, prefixHeadingTwo):
			replacedLine := strings.Replace(line, prefixHeadingTwo, "", 1)
			content := fmt.Sprintf("<h2>%s</h2>", replacedLine)
			writeString(content, &b)

		case strings.HasPrefix(line, prefixHeadingThree):
			replacedLine := strings.Replace(line, prefixHeadingThree, "", 1)
			content := fmt.Sprintf("<h3>%s</h3>", replacedLine)
			writeString(content, &b)

		case strings.HasPrefix(line, prefixBulletPoint):
			replacedLine := strings.Replace(line, prefixBulletPoint, "", 1)
			content := fmt.Sprintf("<li>%s</li>", replacedLine)
			bulletPointBuilder.WriteString(content)
			inBulletPoint = true
			continue

		case regexpPrefixOrder.MatchString(line):
			replacedLine := regexpPrefixOrder.ReplaceAllString(line, "")
			content := fmt.Sprintf("<li>%s</li>", replacedLine)
			orderBuilder.WriteString(content)
			inOrder = true
			continue

		default:
			content := fmt.Sprintf("<p>%s</p>", line)
			writeString(content, &b)
		}
	}

	/*
		入れ子解決処理
		最終行が<li>担っている場合<ul>で閉じる必要がある。それを解決するための処理
	*/
	switch {
	case inBulletPoint:
		contents := bulletPointBuilder.String()
		content := fmt.Sprintf("<ul>%s</ul>", contents)
		writeString(content, &b)
	case inOrder:
		contents := orderBuilder.String()
		content := fmt.Sprintf("<ol>%s</ol>", contents)
		writeString(content, &b)
	}

	return b.String()
}
