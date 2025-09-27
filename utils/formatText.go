package utils

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// ANSI SGR codes
var SGR = map[string]string{
	"reset":     "\x1b[0m",
	"bold":      "\x1b[1m",
	"italic":    "\x1b[3m",
	"underline": "\x1b[4m",

	"black":   "\x1b[30m",
	"red":     "\x1b[31m",
	"green":   "\x1b[32m",
	"yellow":  "\x1b[33m",
	"blue":    "\x1b[34m",
	"magenta": "\x1b[35m",
	"cyan":    "\x1b[36m",
	"white":   "\x1b[37m",
}

// FormatText replaces takes s, and replaces matches to vars
// also applies text with control codes.
func FormatText(s string, vars map[string]string) (string, error) {
	// double brackets are preserved
	const lbraceSentinel = "\x00LBRACE\x00"
	const rbraceSentinel = "\x00RBRACE\x00"
	s = strings.ReplaceAll(s, "{{", lbraceSentinel)
	s = strings.ReplaceAll(s, "}}", rbraceSentinel)

	var b strings.Builder
	for i := 0; i < len(s); {
		r, w := utf8.DecodeRuneInString(s[i:])
		if r != '{' {
			b.WriteRune(r)
			i += w
			continue
		}

		content, next, ok := scanToken(s, i+1)
		if !ok {
			return "", errors.New("unclosed '{' in format string")
		}
		i = next // position after the closing '}'

		// split content on '|' outside quotes.
		parts, err := smartPipeSplit(content)
		if err != nil {
			return "", err
		}
		if len(parts) == 0 {
			continue
		}

		// first segment: either a literal or a var name.
		val, isLiteral, err := resolveHead(parts[0], vars)
		if err != nil {
			return "", err
		}

		// remaining segments: styles.
		if len(parts) > 1 {
			var prefix strings.Builder
			for _, fx := range parts[1:] {
				key := strings.ToLower(strings.TrimSpace(fx))
				if key == "" {
					continue
				}
				code, ok := SGR[key]
				if !ok {
					return "", fmt.Errorf("unknown style %q in {%s}", fx, content)
				}
				prefix.WriteString(code)
			}
			b.WriteString(prefix.String())
			b.WriteString(val)
			// only append reset if we actually applied any styles.
			if prefix.Len() > 0 {
				b.WriteString(SGR["reset"])
			}
		} else {
			b.WriteString(val)
		}

		// Optional: if you want literals to be required to be quoted when they
		// contain spaces or punctuation, enforce it here. Currently this is permissive.
		_ = isLiteral
	}

	out := b.String()
	out = strings.ReplaceAll(out, lbraceSentinel, "{")
	out = strings.ReplaceAll(out, rbraceSentinel, "}")
	return out, nil
}

// scanToken reads until the matching '}' while respecting single quotes and escapes.
func scanToken(s string, start int) (content string, next int, ok bool) {
	var b strings.Builder
	inQuote := false
	escaped := false
	for i := start; i < len(s); i++ {
		ch := s[i]

		if inQuote {
			if escaped {
				// accept any escaped char inside quotes: \' \\ \}
				b.WriteByte(ch)
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '\'' {
				inQuote = false
				b.WriteByte(ch)
				continue
			}
			// Normal quoted char.
			b.WriteByte(ch)
			continue
		}

		// Not in quotes
		if ch == '\'' {
			inQuote = true
			b.WriteByte(ch)
			continue
		}
		if ch == '}' {
			return b.String(), i + 1, true
		}
		b.WriteByte(ch)
	}
	return "", 0, false
}

// smartPipeSplit splits on '|' that are not inside single quotes.
func smartPipeSplit(s string) ([]string, error) {
	parts := []string{}
	var b strings.Builder
	inQuote := false
	escaped := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inQuote {
			if escaped {
				b.WriteByte(ch)
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '\'' {
				inQuote = false
			}
			b.WriteByte(ch)
			continue
		}
		if ch == '\'' {
			inQuote = true
			b.WriteByte(ch)
			continue
		}
		if ch == '|' {
			parts = append(parts, strings.TrimSpace(b.String()))
			b.Reset()
			continue
		}
		b.WriteByte(ch)
	}
	if inQuote {
		return nil, errors.New("unterminated single quote in token")
	}
	parts = append(parts, strings.TrimSpace(b.String()))
	// drop empty segments (e.g., stray pipes).
	filtered := parts[:0]
	for _, p := range parts {
		if p != "" {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

// resolveHead resolves the first token segment as either a quoted literal or a variable.
func resolveHead(head string, vars map[string]string) (val string, isLiteral bool, err error) {
	head = strings.TrimSpace(head)
	if isSingleQuoted(head) {
		txt, err := unquoteSingle(head)
		return txt, true, err
	}
	// variable lookup; case insensitive
	if v, ok := vars[strings.ToLower(head)]; ok {
		return v, false, nil
	}
	return "", false, fmt.Errorf("unknown variable %q", head)
}

func isSingleQuoted(s string) bool {
	return len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\''
}

func unquoteSingle(s string) (string, error) {
	if len(s) < 2 || s[0] != '\'' || s[len(s)-1] != '\'' {
		return "", fmt.Errorf("not a single-quoted string: %q", s)
	}
	return s[1 : len(s)-1], nil
}
