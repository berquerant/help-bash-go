package main

import (
	"regexp"
	"strings"
)

func IsTopLevelComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

var (
	regexTopLevelCommentHead = regexp.MustCompile(`^#+ |#`)
	regexFunctionHead        = regexp.MustCompile(`^[^#][^)]+\(\)`)
	regexVariableDeclHead    = regexp.MustCompile(`^[^#].*=.*`)
)

func ExtractTopLevelComment(line string) (string, bool) {
	if strings.HasPrefix(line, "#") {
		if line == "#" {
			return "", true
		}
		loc := regexTopLevelCommentHead.FindStringIndex(line)
		return line[loc[1]:], true
	}
	return "", false
}

func ExtractFunctionHead(line string) (string, bool) {
	if x := regexFunctionHead.FindString(line); x != "" {
		return x, true
	}
	return "", false
}

func ExtractVariableDeclHead(line string) (string, bool) {
	if x := regexVariableDeclHead.FindString(line); x != "" {
		return x, true
	}
	return "", false
}
