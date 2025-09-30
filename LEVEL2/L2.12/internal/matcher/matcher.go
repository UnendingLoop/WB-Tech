// Package matcher ckecks input line for matching the pattern - string or regexp, returns bool
package matcher

import (
	"strings"

	"grepClone/internal/parser"
)

func FindMatch(line string) bool {
	if parser.SP.IgnoreCase { //-i
		line = strings.ToLower(line)
	}

	switch {
	case parser.SP.ExactMatch: //-F
		return strings.Contains(line, parser.SP.Pattern)
	default:
		return parser.SP.RegexpPattern.MatchString(line)
	}
}
