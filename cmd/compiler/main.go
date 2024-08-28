package main

import (
	"enhanced_python_compiler/internal/generator"
	"enhanced_python_compiler/internal/parser"
	"enhanced_python_compiler/internal/translator"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: compiler <script.py>")
		return
	}

	filename := os.Args[1]

	// 讀取 Python 腳本
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading script: %v\n", err)
		return
	}

	// 解析 Python 代碼
	ast, err := parser.ParsePythonCode(code)
	if err != nil {
		fmt.Printf("Error parsing script: %v\n", err)
		return
	}

	// 將 AST 轉換為 Go 代碼
	goCode, err := translator.TranslateASTToGo(ast)
	if err != nil {
		fmt.Printf("Error translating AST to Go: %v\n", err)
		return
	}

	// 生成最終的 Go 代碼文件
	err = generator.GenerateGoCode(goCode, string(code))
	if err != nil {
		fmt.Printf("Error generating Go code: %v\n", err)
		return
	}

	fmt.Println("Go code has been successfully generated in output.go")
}
