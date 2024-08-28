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
	// 使用 Python 生成 JSON 格式的 AST
	cmd := exec.Command("python3", "-c", fmt.Sprintf(`
import ast, json, sys
parsed = ast.parse('''%s''')
def ast_to_dict(node):
    if isinstance(node, list):
        return [ast_to_dict(item) for item in node]
    elif isinstance(node, ast.AST):
        result = {"_type": node.__class__.__name__}
        for field in node._fields:
            result[field] = ast_to_dict(getattr(node, field))
        return result
    else:
        return node
print(json.dumps(ast_to_dict(parsed), indent=2))
`, sourceCode))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to parse Python source code: %v\n%s", err, output)
	}

	// 将生成的 JSON 解析为 Go 数据结构
	var astRoot interface{}
	err = json.Unmarshal(output, &astRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal AST JSON: %v", err)
	}

	return &AST{Root: astRoot}, nil
}
