package executor

import (
	"enhanced_python_compiler/internal/parser"
	"enhanced_python_compiler/internal/runtime"
	"enhanced_python_compiler/internal/translator"
	"fmt"
)

func Execute(ast *parser.AST) (string, error) {
	// 判斷是否需要轉換為 Go 代碼
	// 初期我們可以選擇所有代碼都交給 CPython
	useGo := false // 在此設定初始策略

	if useGo {
		goCode, err := translator.TranslateToGo(ast)
		if err != nil {
			return "", fmt.Errorf("translation to Go failed: %v", err)
		}
		return runtime.ExecuteGo(goCode)
	} else {
		return runtime.ExecuteCPython(ast)
	}
}
