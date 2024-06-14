package CodeWriter

import "strings"

type GoFile struct {
	FileName string
	Content  string
}

func WriteBacktickString(lines []string) string {
	content := strings.Join(lines, "\n")

	var builder strings.Builder

	builder.WriteRune('`')
	builder.WriteString(content)
	builder.WriteRune('`')

	return builder.String()
}

func ForWriter(index, variable, list, content string) []string {
	var builder strings.Builder

	builder.WriteString("for ")
	builder.WriteString(index)
	builder.WriteString(", ")
	builder.WriteString(variable)
	builder.WriteString(" := range ")
	builder.WriteString(list)
	builder.WriteString(" {")

	lines := []string{
		builder.String(),
	}

	subLines := strings.Split(content, "\n")

	for _, line := range subLines {
		builder.Reset()

		builder.WriteString("\t")
		builder.WriteString(line)

		lines = append(lines, builder.String())
	}

	lines = append(lines, "}")

	return lines
}
