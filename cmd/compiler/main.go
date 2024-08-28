package main

import (
	"enhanced_python_compiler/internal/executor"
	"enhanced_python_compiler/internal/parser"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: compiler <script.py>")
		return
	}

	filename := os.Args[1]
	sourceCode, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// 解析 Python 代碼，生成 AST
	ast, err := parser.Parse(sourceCode)
	if err != nil {
		fmt.Printf("Error parsing script: %v\n", err)
		return
	}

	// 執行代碼
	result, err := executor.Execute(ast)
	if err != nil {
		fmt.Printf("Error executing script: %v\n", err)
		return
	}

	fmt.Println("Execution result:", result)
}
