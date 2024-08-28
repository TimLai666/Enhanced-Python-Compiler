package translator

import (
	"enhanced_python_compiler/internal/parser"
	"fmt"
	"strings"
)

// TranslateASTToGo 將 Python 的 AST 轉換為 Go 代碼
func TranslateASTToGo(ast *parser.AST) (string, error) {
	// 初始 Go 代碼模板
	goCode := "package main\n\nimport \"fmt\"\n\nfunc main() {\n"

	// 處理 AST 的字符串表示
	if strings.Contains(ast.Root, "Call(func=Name(id='print'") {
		// 假設我們可以將 print 語句簡單轉換為 fmt.Println
		goCode += "\tfmt.Println(\"Hello from Go!\")\n"
	} else {
		return "", fmt.Errorf("unsupported AST structure: %v", ast.Root)
	}

	goCode += "}\n"
	return goCode, nil
}
