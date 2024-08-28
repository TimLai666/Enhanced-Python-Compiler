package translator

import (
	"enhanced_python_compiler/internal/parser"
	"fmt"
)

// TranslateASTToGo 将 Python 的 AST 转换为 Go 代码
func TranslateASTToGo(ast *parser.AST) (string, error) {
	// 初始 Go 代码模板
	goCode := "package main\n\nimport \"fmt\"\n\nfunc main() {\n"

	// 假设 AST 是一个字典，包含 "body" 字段
	if rootMap, ok := ast.Root.(map[string]interface{}); ok {
		if body, ok := rootMap["body"].([]interface{}); ok {
			for _, stmt := range body {
				if stmtMap, ok := stmt.(map[string]interface{}); ok {
					if stmtType, ok := stmtMap["_type"].(string); ok {
						switch stmtType {
						case "Expr":
							// 处理表达式语句，例如 print()
							if value, ok := stmtMap["value"].(map[string]interface{}); ok {
								if funcCall, ok := value["func"].(map[string]interface{}); ok {
									if funcName, ok := funcCall["id"].(string); ok && funcName == "print" {
										// 处理 print 函数
										if args, ok := value["args"].([]interface{}); ok && len(args) > 0 {
											if arg, ok := args[0].(map[string]interface{}); ok {
												if constValue, ok := arg["s"].(string); ok {
													goCode += fmt.Sprintf("\tfmt.Println(\"%s\")\n", constValue)
												}
											}
										}
									}
								}
							}
						default:
							fmt.Printf("Unsupported statement type: %s\n", stmtType)
						}
					}
				}
			}
		} else {
			return "", fmt.Errorf("unsupported AST structure: body not found")
		}
	} else {
		return "", fmt.Errorf("unsupported AST structure")
	}

	goCode += "}\n"
	return goCode, nil
}
