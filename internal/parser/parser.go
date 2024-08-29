package parser

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type AST struct {
	Root interface{} // 這裡使用 interface{} 來處理 AST 的各種結構
}

// ParsePythonCode 解析 Python 代碼並返回 AST 結構
func ParsePythonCode(sourceCode []byte) (*AST, error) {
	// 使用 Python 的 ast 模塊來生成 JSON 格式的 AST
	cmd := exec.Command("python3", "-c", fmt.Sprintf(`
import ast, json, sys
parsed = ast.parse('''%s''')
print(json.dumps(ast.dump(parsed, annotate_fields=True, include_attributes=True)))
`, sourceCode))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to parse Python source code: %v\nOutput: %s", err, output)
	}

	var astRoot interface{}
	err = json.Unmarshal(output, &astRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal AST JSON: %v", err)
	}

	return &AST{Root: astRoot}, nil
}
