package shdoc

import (
	"strings"
)

// trimExampleLeadingTabs removes up to count leading tabs from each line that
// begins with a tab, but only if the first line of the example begins with at
// least count tabs. Lines that begin with any other character are left as-is.
func trimExampleLeadingTabs(text string, count int) string {
	if count <= 0 || text == "" {
		return text
	}

	first := text
	if i := strings.IndexByte(text, '\n'); i >= 0 {
		first = text[:i]
	}
	if !strings.HasPrefix(first, strings.Repeat("\t", count)) {
		return text
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line == "" || line[0] != '\t' {
			continue
		}
		trimmed := 0
		for trimmed < count && trimmed < len(line) && line[trimmed] == '\t' {
			trimmed++
		}
		lines[i] = line[trimmed:]
	}

	return strings.Join(lines, "\n")
}

// unindent removes common leading whitespace from text lines.
// Matches the awk implementation precisely.
func unindent(text string) string {
	lines := strings.Split(text, "\n")

	// Find first non-empty line and minimum indent in a single pass.
	start := -1
	indent := -1 // -1 = not yet set
	for i, line := range lines {
		if line == "" {
			continue
		}
		if start == -1 {
			start = i
		}
		spaces := 0
		for _, ch := range line {
			if ch == ' ' {
				spaces++
			} else {
				break
			}
		}
		if indent == -1 || spaces < indent {
			indent = spaces
		}
	}

	// If no non-empty lines found, return empty
	if start == -1 {
		return ""
	}

	// Remove indent and join from start
	var result strings.Builder
	for i := start; i < len(lines); i++ {
		if i > start {
			result.WriteString("\n")
		}
		if len(lines[i]) > indent {
			result.WriteString(lines[i][indent:])
		} else {
			result.WriteString(lines[i])
		}
	}

	return result.String()
}

func formatExample(text string, trimTabs int) string {
	return unindent(trimExampleLeadingTabs(text, trimTabs))
}
