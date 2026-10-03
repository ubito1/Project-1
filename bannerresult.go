package main

import (
	"fmt"
	"os"
	"strings"
)


//func to prepare the shadow file for use.
func LoadBanner(filename string) ([]string, error) {
	in, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(in), "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}

	if len(lines) < 9*95 {
		return nil, fmt.Errorf("banner file %s is not valid", filename)
	}
	return lines, nil
}

//func to accept the 
func Render(banner []string, input string) (string, error) {
	if input == "" {
		return "", nil
	}

	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\\n", "\n")

	parts := strings.Split(input, "\n")

	isOnlyNewlines := true
	for _, part := range parts {
		if part != "" {
			isOnlyNewlines = false
			break
		}
	}

	var result strings.Builder

	if isOnlyNewlines {
		for i := 0; i < len(parts)-1; i++ {
			result.WriteString("\n")
		}
		return result.String(), nil
	}
//,pas
	for i, part := range parts {
		if part == "" {
			if i > 0 {
				result.WriteString("\n")
			}
			continue
		}

		renderedPart, err := renderLine(banner, part)
		if err != nil {
			return "", err
		}
		result.WriteString(renderedPart)
	}

	return result.String(), nil
}

func getCharLines(lines []string, c rune) []string {
	startindx := int(c-' ') * 9
	if startindx < 0 || startindx+9 > len(lines) {
		return make([]string, 9)
	}
	var charLines []string
	charLines = append(charLines, lines[startindx:startindx+9]...)
	return charLines
}

func renderLine(banner []string, text string) (string, error) {
	var lineResult strings.Builder

	for row := 1; row < 9; row++ {
		for _, c := range text {
			if c < ' ' || c > '~' {
				return "", fmt.Errorf("character %q is not supported", c)
			}
			charLines := getCharLines(banner, c)
			lineResult.WriteString(charLines[row])
		}
		lineResult.WriteString("\n")
	}
	return lineResult.String(), nil
}	