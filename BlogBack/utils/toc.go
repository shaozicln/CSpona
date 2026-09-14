package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	headingLineRe = regexp.MustCompile(`(?m)^(#{1,6})[ \t]+(.+?)[ \t]*$`)
	urlInTextRe   = regexp.MustCompile(`https?://\S+`)
)

// TocItem 文章目录项（id 与前端 marked heading 渲染保持一致）
type TocItem struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Level int    `json:"level"`
}

// BuildMarkdownTOC 从 Markdown 抽取标题目录；跳过 fenced code；含 URL 的标题忽略。
func BuildMarkdownTOC(md string) []TocItem {
	plain := stripFencedCode(md)
	matches := headingLineRe.FindAllStringSubmatch(plain, -1)
	if len(matches) == 0 {
		return []TocItem{}
	}

	counter := map[string]int{}
	toc := make([]TocItem, 0, len(matches))
	for _, m := range matches {
		level := len(m[1])
		text := strings.TrimSpace(m[2])
		text = strings.TrimSpace(strings.TrimSuffix(text, "#"))
		text = strings.TrimSpace(text)
		if text == "" || urlInTextRe.MatchString(text) {
			continue
		}
		id := headingID(text, counter)
		toc = append(toc, TocItem{ID: id, Text: text, Level: level})
	}
	return toc
}

func stripFencedCode(md string) string {
	var b strings.Builder
	b.Grow(len(md))
	lines := strings.Split(md, "\n")
	inFence := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "```") {
			inFence = !inFence
			b.WriteByte('\n')
			continue
		}
		if inFence {
			b.WriteByte('\n')
			continue
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// headingID 对齐前端：NFD 去音标、空格变 -、保留字母数字下划线汉字与 -，重复则 -2 -3...
func headingID(text string, counter map[string]int) string {
	slug := slugifyHeading(text)
	if slug == "" {
		slug = "section"
	}
	n := counter[slug] + 1
	counter[slug] = n
	if n > 1 {
		return slug + "-" + itoa(n)
	}
	return slug
}

func slugifyHeading(text string) string {
	// 近似 NFD：去掉常见 Combining marks
	var b strings.Builder
	b.Grow(len(text))
	prevHyphen := false
	for _, r := range text {
		if unicode.In(r, unicode.Mn) {
			continue
		}
		if unicode.IsSpace(r) {
			if !prevHyphen && b.Len() > 0 {
				b.WriteByte('-')
				prevHyphen = true
			}
			continue
		}
		if isWordChar(r) || (r >= 0x4e00 && r <= 0x9fa5) || r == '-' {
			b.WriteRune(unicode.ToLower(r))
			prevHyphen = r == '-'
			continue
		}
		// 其它字符丢弃（对齐 JS replace 为 ""）
	}
	s := strings.Trim(b.String(), "-")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return s
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [12]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
