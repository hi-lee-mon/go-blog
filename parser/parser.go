package parser

import (
	"fmt"
	"strings"
)

const (
	PrefixHeadingOne   = "# "
	PrefixHeadingTwo   = "## "
	PrefixHeadingThree = "### "
	PrefixBulletPoint  = "- "
)

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
	var b strings.Builder

	// 状態
	inBulletPoint := false
	var bulletPointBuilder strings.Builder

	// パース処理
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// 状態毎の処理
		if inBulletPoint {
			if strings.HasPrefix(line, PrefixBulletPoint) {
				replacedLine := strings.Replace(line, PrefixBulletPoint, "", 1)
				content := fmt.Sprintf("<li>%s</li>", replacedLine)
				bulletPointBuilder.WriteString(content)
				continue
			} else {
				contents := bulletPointBuilder.String()
				content := fmt.Sprintf("<ul>%s</ul>", contents)
				writeString(content, &b)
				// 状態リセット
				inBulletPoint = false
				bulletPointBuilder.Reset()
			}
		}

		// 状態なし処理
		if strings.HasPrefix(line, PrefixHeadingOne) {
			replacedLine := strings.Replace(line, PrefixHeadingOne, "", 1)
			content := fmt.Sprintf("<h1>%s</h1>", replacedLine)
			writeString(content, &b)
		} else if strings.HasPrefix(line, PrefixHeadingTwo) {
			replacedLine := strings.Replace(line, PrefixHeadingTwo, "", 1)
			content := fmt.Sprintf("<h2>%s</h2>", replacedLine)
			writeString(content, &b)
		} else if strings.HasPrefix(line, PrefixHeadingThree) {
			replacedLine := strings.Replace(line, PrefixHeadingThree, "", 1)
			content := fmt.Sprintf("<h3>%s</h3>", replacedLine)
			writeString(content, &b)
		} else if strings.HasPrefix(line, PrefixBulletPoint) {
			replacedLine := strings.Replace(line, PrefixBulletPoint, "", 1)
			content := fmt.Sprintf("<li>%s</li>", replacedLine)
			bulletPointBuilder.WriteString(content)
			inBulletPoint = true
			continue
		} else {
			content := fmt.Sprintf("<p>%s</p>", line)
			writeString(content, &b)
		}
	}

	if inBulletPoint {
		contents := bulletPointBuilder.String()
		content := fmt.Sprintf("<ul>%s</ul>", contents)
		writeString(content, &b)
	}

	return b.String()
}
