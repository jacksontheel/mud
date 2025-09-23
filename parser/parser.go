package parser

import (
	"strings"

	"example.com/mud/commands"
	"example.com/mud/models"
)

// TODO: register multi-word merges at the command level instead of in parser
var multiWordVerbMerges = [][]string{
	{"pick", "up"},
	{"make", "out"},
}

var directionAliases = commands.DirectionAliases()

var verbAliases = commands.AllAliases()

var patterns = commands.AllPatterns()

func tokenize(input string) []string {
	s := strings.ToLower(strings.TrimSpace(input))
	parts := splitCompact(s)

	// merge multi-word verbs
	parts = mergeSequences(parts, multiWordVerbMerges)

	// normalize verb aliases
	if len(parts) > 0 {
		if base, ok := verbAliases[parts[0]]; ok {
			parts[0] = base
		}
	}
	return parts
}

func splitCompact(s string) []string {
	out := []string{}
	for _, p := range strings.Fields(s) {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func mergeSequences(tokens []string, merges [][]string) []string {
	if len(tokens) == 0 {
		return tokens
	}
	i := 0
	var out []string
	for i < len(tokens) {
		merged := false
		for _, seq := range merges {
			if i+len(seq) <= len(tokens) {
				match := true
				for j, w := range seq {
					if tokens[i+j] != w {
						match = false
						break
					}
				}
				if match {
					out = append(out, strings.Join(seq, ""))
					i += len(seq)
					merged = true
					break
				}
			}
		}
		if !merged {
			out = append(out, tokens[i])
			i++
		}
	}
	return out
}

func tryMatch(p models.Pattern, tokens []string) (ok bool, params map[string]string) {
	params = map[string]string{}
	ti := 0
	for pi := 0; pi < len(p.Tokens); pi++ {
		pt := p.Tokens[pi]

		// if pattern expects a Literal, e.g. "take"
		if pt.Literal != "" {
			if ti >= len(tokens) {
				return false, nil
			}
			if tokens[ti] != pt.Literal {
				return false, nil
			}
			ti++
			continue
		}

		// if pattern expects slot to be the remaining tokens, e.g. "take big orange key"
		if pt.SlotIsRest {
			if ti > len(tokens) {
				return false, nil
			}
			rest := tokens[ti:]
			val, ok := validateSlot(pt.SlotType, rest)
			if !ok {
				return false, nil
			}
			params[pt.SlotName] = val
			ti = len(tokens)
			continue
		}

		// consume the next single token
		if ti >= len(tokens) {
			return false, nil
		}
		val, ok := validateSlot(pt.SlotType, []string{tokens[ti]})
		if !ok {
			return false, nil
		}
		params[pt.SlotName] = val
		ti++
	}

	// must consume all tokens
	if ti != len(tokens) {
		return false, nil
	}

	return true, params
}

func validateSlot(SlotType string, toks []string) (string, bool) {
	switch SlotType {
	case "direction":
		if len(toks) != 1 {
			return "", false
		}

		if canon, ok := directionAliases[toks[0]]; ok {
			return canon, true
		}

		return "", false
	default:
		if len(toks) == 0 {
			return "", false
		}
		return strings.Join(toks, " "), true
	}
}

func Parse(input string) commands.Command {
	toks := tokenize(input)
	if len(toks) == 0 {
		return commands.Command{Kind: "", Params: nil}
	}

	for _, p := range patterns {
		if ok, params := tryMatch(p, toks); ok {
			return commands.Command{Kind: p.Kind, Params: params}
		}
	}

	return commands.Command{Kind: "", Params: nil}
}
