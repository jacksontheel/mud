package parser

import (
	"strings"

	"example.com/mud/models"
	"example.com/mud/parser/commands"
)

// TODO: register multi-word merges at the command level instead of in parser
var multiWordVerbMerges = [][]string{
	{"pick", "up"},
	{"make", "out"},
}

func tokenize(input string) []string {
	s := strings.ToLower(strings.TrimSpace(input))
	parts := splitCompact(s)

	// merge multi-word verbs
	parts = mergeSequences(parts, multiWordVerbMerges)

	// normalize verb aliases
	if len(parts) > 0 {
		if base, ok := commands.VerbAliases[parts[0]]; ok {
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

func tryMatch(p models.Pattern, tokens []string) (*models.Command, bool) {
	params := map[string]string{}
	ti := 0
	for pi := 0; pi < len(p.Tokens); pi++ {
		pt := p.Tokens[pi]

		// if pattern expects a Literal, e.g. "take"
		if pt.Literal != "" {
			if ti >= len(tokens) {
				return nil, false
			}
			if tokens[ti] != pt.Literal {
				return nil, false
			}
			ti++
			continue
		}

		// if pattern expects slot to be the remaining tokens, e.g. "take big orange key"
		if pt.SlotIsRest {
			if ti > len(tokens) {
				return nil, false
			}
			rest := tokens[ti:]
			val, ok := validateSlot(pt.SlotName, rest)
			if !ok {
				return nil, false
			}
			params[pt.SlotName] = val
			ti = len(tokens)
			continue
		}

		// consume the next single token
		if ti >= len(tokens) {
			return nil, false
		}
		val, ok := validateSlot(pt.SlotName, []string{tokens[ti]})
		if !ok {
			return nil, false
		}
		params[pt.SlotName] = val
		ti++
	}

	// must consume all tokens
	if ti != len(tokens) {
		return nil, false
	}

	return &models.Command{
		Kind:           p.Kind,
		Params:         params,
		NoMatchMessage: p.NoMatchMessage,
	}, true
}

func validateSlot(SlotType string, toks []string) (string, bool) {
	switch SlotType {
	case "direction":
		if len(toks) != 1 {
			return "", false
		}

		if canon, ok := commands.DirectionAliases[toks[0]]; ok {
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

func Parse(input string) *models.Command {
	toks := tokenize(input)
	if len(toks) == 0 {
		return nil
	}

	for _, p := range commands.Patterns {
		if cmd, ok := tryMatch(p, toks); ok {
			return cmd
		}
	}

	return nil
}
