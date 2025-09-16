package dsl

import (
	"log"
	"os"

	"example.com/mud/world/entities"
	participle "github.com/alecthomas/participle/v2"
)

func LoadEntitiesFromFile(fileName string) (map[string]*entities.Entity, error) {
	parser := participle.MustBuild[DSL](
		participle.Lexer(dslLexer),
		participle.Elide("Whitespace"),
		participle.Unquote("String"),
	)

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	ast, err := parser.ParseString("", string(data))
	if err != nil {
		log.Fatal(err)
	}

	entities, err := buildAll(ast)
	return entities, err
}
