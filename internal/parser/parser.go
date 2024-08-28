package parser

import (
	"fmt"
	"os/exec"
)

type AST struct {
	Root string // 將 AST 結構保存為字符串
}

func ParsePythonCode(sourceCode []byte) (*AST, error) {
	// 使用 Python 的 ast 模組來解析代碼並返回 AST
	cmd := exec.Command("python3", "-c", fmt.Sprintf("import ast; print(ast.dump(ast.parse('''%s'''))) ", sourceCode))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to parse Python source code: %v", err)
	}

	return &AST{Root: string(output)}, nil
}
