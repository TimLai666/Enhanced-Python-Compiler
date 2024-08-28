package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type AST struct {
	Root *ast.File
}

func Parse(sourceCode []byte) (*AST, error) {
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, "", sourceCode, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source code: %v", err)
	}

	return &AST{Root: root}, nil
}
