package main

import (
	"enhanced_python_compiler/internal/executor"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: compiler <script.py>")
		return
	}

	filename := os.Args[1]

	// 執行 Python 腳本
	result, err := executor.Execute(filename)
	if err != nil {
		fmt.Printf("Error executing script: %v\n", err)
		return
	}

	fmt.Println("Execution result:", result)
}
