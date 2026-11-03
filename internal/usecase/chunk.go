package usecase

import (
	"strings"
	"unicode/utf8"
)

type section struct {
	heading string
	body    string
}

func parseHeading(line string) (level int, title string, ok bool) {
	trimmed := strings.TrimLeft(line, "#")
	level = len(line) - len(trimmed)
	trimmed = strings.TrimRight(trimmed, " ")
	if level > 6 || level == 0 {
		return 0, "", false
	}
	if strings.HasPrefix(trimmed, " ") {
		trimmed = strings.TrimLeft(trimmed, " ")
		return level, trimmed, true
	}
	return 0, "", false
}

func headingPath(stack []string) string {
	newStack := []string{}
	for idx, s := range stack {
		if stack[idx] != "" {
			newStack = append(newStack, s)
		}
	}
	return strings.Join(newStack, " > ")
}

func splitSections(text string) []section {
	stackOfStrings := make([]string, 6)
	bodyLines := []string{}
	sections := []section{}

	inCodeBlock := false

	flush := func() {
		body := strings.Join(bodyLines, "\n")

		if len(strings.TrimSpace(body)) != 0 {
			sections = append(sections, section{heading: headingPath(stackOfStrings), body: body})
		}
		bodyLines = nil
	}

	for line := range strings.SplitSeq(text, "\n") {
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			bodyLines = append(bodyLines, line)
			continue
		}
		if !inCodeBlock {
			if level, title, ok := parseHeading(line); ok {
				flush()
				stackOfStrings[level-1] = title
				for i := level; i < len(stackOfStrings); i++ {
					stackOfStrings[i] = ""
				}
				continue
			}
		}
		bodyLines = append(bodyLines, line)
	}
	flush()

	return sections
}

func packParagraphs(paragraphs []string, maxLen int) []string {
	result := []string{}
	current := ""
	for _, paragraph := range paragraphs {
		if current == "" {
			current = paragraph
		} else if utf8.RuneCountInString(current)+2+utf8.RuneCountInString(paragraph) <= maxLen {
			current = current + "\n\n" + paragraph
		} else {
			result = append(result, current)
			current = paragraph
		}
	}
	if current != "" {
		result = append(result, current)
	}

	return result
}

type Chunk struct {
	Source  string
	Heading string
	Index   int
	Content string
}

func SplitMarkdown(source, text string, maxLen, overlap int) []Chunk {
	if maxLen <= 0 {
		maxLen = 1000
	}
	if overlap >= maxLen {
		overlap = maxLen / 2
	}
	res := []Chunk{}
	index := 0
	for _, sec := range splitSections(text) {
		var paragraphs []string
		for b := range strings.SplitSeq(sec.body, "\n\n") {
			trimmed := strings.TrimSpace(b)
			if trimmed == "" {
				continue
			}
			paragraphs = append(paragraphs, trimmed)
		}
		pieces := packParagraphs(paragraphs, maxLen)
		for i, piece := range pieces {
			text := piece
			if i > 0 && overlap > 0 {
				text = tailRunes(pieces[i-1], overlap) + "\n\n" + text
			}
			if sec.heading != "" {
				text = sec.heading + "\n\n" + text
			}
			res = append(res, Chunk{Source: source, Heading: sec.heading, Index: index, Content: text})
			index++
		}
	}

	return res
}

func tailRunes(text string, n int) string {

	if utf8.RuneCountInString(text) <= n {
		return text
	}
	byteIdx := len(text)
	for i := 0; i < n; i++ {
		_, size := utf8.DecodeLastRuneInString(text[:byteIdx])
		byteIdx -= size
	}
	return text[byteIdx:]
}

func tailRunesWithSpace(text string) string {
	_, after, ok := strings.Cut(text, " ")
	if !ok {
		return text
	}
	return after
}
