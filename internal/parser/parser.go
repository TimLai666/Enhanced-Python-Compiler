package parser

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type AST struct {
	Root string // 使用 string 來存儲 AST 的根節點
}

func ParsePythonCode(sourceCode []byte) (*AST, error) {
	// 使用 Python 代碼生成 JSON 格式的 AST
	cmd := exec.Command("python", "-c", fmt.Sprintf(`
import ast, json, sys
source_code = '''%s'''
parsed_ast = ast.parse(source_code)

def ast_to_json(node):
    if isinstance(node, ast.AST):
        fields = {field: ast_to_json(getattr(node, field)) for field in node._fields}
        return {'_type': node.__class__.__name__, **fields}
    elif isinstance(node, list):
        return [ast_to_json(item) for item in node]
    else:
        return node

json_ast = json.dumps(ast_to_json(parsed_ast))
print(json_ast)
`, sourceCode))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to parse Python source code: %v\nOutput: %s", err, string(output))
	}

	return &AST{Root: string(output)}, nil
}

// ConvertASTStringToJSON 將字串模式的 AST 轉換為 JSON
func ConvertASTStringToJSON(astStr string) (string, error) {
	// 替換 Python 的 None, True, False 為 JSON 兼容的 null, true, false
	replacements := map[string]string{
		"None":  "null",
		"True":  "true",
		"False": "false",
	}

	// 逐個進行替換
	for pyVal, jsonVal := range replacements {
		astStr = strings.Replace(astStr, pyVal, jsonVal, -1)
	}

	// 將單引號替換為雙引號
	astStr = strings.Replace(astStr, "'", "\"", -1)

	// 使用正則表達式在關鍵位置加入逗號和大括號，使其成為合法的 JSON 格式
	re := regexp.MustCompile(`(?m)(\w+)=`)
	astStr = re.ReplaceAllString(astStr, "\"$1\":")

	return astStr, nil
}
