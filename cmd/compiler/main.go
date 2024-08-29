package main

import (
	"enhanced_python_compiler/internal/packager"
	"enhanced_python_compiler/internal/parser"
	"enhanced_python_compiler/internal/translator"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: compiler <python_file>")
		return
	}

	pythonFile := os.Args[1]
	sourceCode, err := os.ReadFile(pythonFile)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return
	}

	// Parse the Python code into an AST
	ast, err := parser.ParsePythonCode(sourceCode)
	if err != nil {
		fmt.Printf("Error parsing Python code: %v\n", err)
		return
	}

	// Translate the AST to binary data
	binaryData, _, err := translator.TranslateASTToBinary(ast)
	if err != nil {
		fmt.Printf("Error translating AST to binary: %v\n", err)
		return
	}

	// Package the binary data into a single executable
	err = packager.PackageExecutable(string(binaryData))
	if err != nil {
		fmt.Printf("Error packaging executable: %v\n", err)
		return
	}

	fmt.Println("Packaging successful!")
}
