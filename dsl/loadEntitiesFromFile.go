package dsl

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"example.com/mud/world/entities"
	participle "github.com/alecthomas/participle/v2"
)

func LoadEntitiesFromFile(directoryName string) (map[string]*entities.Entity, error) {
	parser := participle.MustBuild[DSL](
		participle.Lexer(dslLexer),
		participle.Elide("Whitespace"),
		participle.Unquote("String"),
	)

	// create container for file contents
	var ast = &DSL{}

	filepath.WalkDir(directoryName, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf("Something went wrong: %v", err)
		}

		// skip folders and files without .mud extension
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".mud") {
			return nil
		}

		// read file
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", path, err)
		}

		// parse file
		fileSyntaxTree, err := parser.ParseString("", string(data))
		if err != nil {
			log.Fatalf("Failed to parse %s: %v", path, err)
		}

		// append parsed dsl to the ast
		ast.Entities = append(ast.Entities, fileSyntaxTree.Entities...)
		return nil
	})

	entities, err := buildAll(ast)
	return entities, err
}
