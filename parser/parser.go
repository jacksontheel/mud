package parser

import (
	"regexp"
	"strings"

	"example.com/mud/models"
)

const (
	CommandMove   = "move"
	CommandTake   = "take"
	CommandLook   = "look"
	CommandAttack = "attack"
	CommandKiss   = "kiss"
)

type Command struct {
	Kind   string
	Params map[string]string
}

var directionAliases = map[string]string{
	"n":     models.DirectionNorth,
	"north": models.DirectionNorth,

	"s":     models.DirectionSouth,
	"south": models.DirectionSouth,

	"e":    models.DirectionEast,
	"east": models.DirectionEast,

	"w":    models.DirectionWest,
	"west": models.DirectionWest,

	"u":  models.DirectionUp,
	"up": models.DirectionUp,

	"d":    models.DirectionDown,
	"down": models.DirectionDown,
}

var verbAliases = map[string]string{
	"go":   CommandMove,
	"walk": CommandMove,
	"move": CommandMove,

	"take":    CommandTake,
	"get":     CommandTake,
	"grab":    CommandTake,
	"pickup":  CommandTake,
	"collect": CommandTake,

	"look":    CommandLook,
	"examine": CommandLook,
	"inspect": CommandLook,
	"l":       CommandLook,

	"attack": CommandAttack,
	"kill":   CommandAttack,
	"hit":    CommandAttack,

	"kiss":    CommandKiss,
	"smooch":  CommandKiss,
	"makeout": CommandKiss,
}

var multiWordVerbMerges = [][]string{
	{"pick", "up"},
	{"make", "out"},
}

type patToken struct {
	literal    string
	slotName   string
	slotIsRest bool
	slotType   string
}

type pattern struct {
	kind   string
	tokens []patToken
}

func lit(word string) patToken {
	return patToken{literal: word}
}
func slot(slotType, name string) patToken {
	return patToken{slotName: name, slotType: slotType}
}
func slotRest(slotType, name string) patToken {
	return patToken{slotName: name, slotType: slotType, slotIsRest: true}
}

var patterns = []pattern{
	{kind: CommandMove, tokens: []patToken{
		slot("direction", "direction"),
	}},
	{kind: CommandMove, tokens: []patToken{
		lit(CommandMove),
		slot("direction", "direction"),
	}},

	{kind: CommandTake, tokens: []patToken{
		lit(CommandTake),
		slotRest("object", "object"),
	}},

	{kind: CommandLook, tokens: []patToken{
		lit(CommandLook),
	}},
	{kind: CommandLook, tokens: []patToken{
		lit(CommandLook),
		slotRest("object", "object"),
	}},
	{kind: CommandAttack, tokens: []patToken{
		lit(CommandAttack),
		slotRest("object", "object"),
	}},
	{kind: CommandKiss, tokens: []patToken{
		lit(CommandKiss),
		slotRest("object", "object"),
	}},
}

var punctRe = regexp.MustCompile(`[[:punct:]]`)

func tokenize(input string) []string {
	s := strings.ToLower(strings.TrimSpace(input))
	// remove punctuation
	s = punctRe.ReplaceAllString(s, " ")
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

func tryMatch(p pattern, tokens []string) (ok bool, params map[string]string) {
	params = map[string]string{}
	ti := 0
	for pi := 0; pi < len(p.tokens); pi++ {
		pt := p.tokens[pi]

		// if pattern expects a literal, e.g. "take"
		if pt.literal != "" {
			if ti >= len(tokens) {
				return false, nil
			}
			if tokens[ti] != pt.literal {
				return false, nil
			}
			ti++
			continue
		}

		// if pattern expects slot to be the remaining tokens, e.g. "take big orange key"
		if pt.slotIsRest {
			if ti > len(tokens) {
				return false, nil
			}
			rest := tokens[ti:]
			val, ok := validateSlot(pt.slotType, rest)
			if !ok {
				return false, nil
			}
			params[pt.slotName] = val
			ti = len(tokens)
			continue
		}

		// consume the next single token
		if ti >= len(tokens) {
			return false, nil
		}
		val, ok := validateSlot(pt.slotType, []string{tokens[ti]})
		if !ok {
			return false, nil
		}
		params[pt.slotName] = val
		ti++
	}

	// must consume all tokens
	if ti != len(tokens) {
		return false, nil
	}

	return true, params
}

func validateSlot(slotType string, toks []string) (string, bool) {
	switch slotType {
	case "direction":
		if len(toks) != 1 {
			return "", false
		}

		if canon, ok := directionAliases[toks[0]]; ok {
			return canon, true
		}

		return "", false
	case "object":
		if len(toks) == 0 {
			return "", false
		}

		// the object name is the remaining tokens joined
		return strings.Join(toks, " "), true
	default:
		return "", false
	}
}

func Parse(input string) Command {
	toks := tokenize(input)
	if len(toks) == 0 {
		return Command{Kind: "", Params: nil}
	}

	for _, p := range patterns {
		if ok, params := tryMatch(p, toks); ok {
			return Command{Kind: p.kind, Params: params}
		}
	}

	return Command{Kind: "", Params: nil}
}
