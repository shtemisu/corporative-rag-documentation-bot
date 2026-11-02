package usecase

import (
	"strings"
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
	sections := splitSections(text)
	if len(sections) <= maxLen {

	}

	return nil
}
