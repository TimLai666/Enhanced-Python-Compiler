package parser

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type AST struct {
	Root interface{} // 使用 interface{} 来处理复杂的结构
}

func ParsePythonCode(sourceCode []byte) (*AST, error) {
	// 使用 Python 代码直接生成 JSON 格式的 AST
	cmd := exec.Command("python3", "-c", fmt.Sprintf(`
import ast, json, sys
parsed = ast.parse('''%s''')
print(json.dumps(ast.dump(parsed, annotate_fields=False, include_attributes=False)))
`, sourceCode))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to parse Python source code: %v", err)
	}

	var astRoot interface{}
	err = json.Unmarshal(output, &astRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal AST JSON: %v", err)
	}

	// 打印 AST 结构
	fmt.Printf("Parsed AST: %+v\n", astRoot)

	return &AST{Root: astRoot}, nil
}
