package parser

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type AST struct {
	Root interface{} // 使用 interface{} 以便处理复杂结构
}

func ParsePythonCode(sourceCode []byte) (*AST, error) {
	// 使用 Python 代码生成 JSON 格式的 AST
	cmd := exec.Command("python3", "-c", fmt.Sprintf(`
import ast, json, sys
parsed = ast.parse('''%s''')
print(json.dumps(ast.dump(parsed, annotate_fields=True, include_attributes=True)))
`, sourceCode))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to parse Python source code: %v", err)
	}

	// 打印出原始的AST
	fmt.Println("Generated AST:", string(output))

	var astRoot interface{}
	err = json.Unmarshal(output, &astRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal AST JSON: %v", err)
	}

	return &AST{Root: astRoot}, nil
}
