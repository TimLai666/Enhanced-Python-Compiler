package translator

import (
	"encoding/json"
	"enhanced_python_compiler/internal/parser"
	"fmt"
)

// TranslateASTToGo 将 Python 的 AST 转换为 Go 代码
func TranslateASTToGo(ast *parser.AST) (string, error) {
	goCode := "package main\n\nimport \"fmt\"\n\n"

	// 处理 AST 的顶层 Module 节点
	if rootMap, ok := ast.Root.(map[string]interface{}); ok {
		if body, ok := rootMap["body"].([]interface{}); ok {
			for _, stmt := range body {
				if stmtMap, ok := stmt.(map[string]interface{}); ok {
					if stmtType, ok := stmtMap["_type"].(string); ok {
						fmt.Printf("Processing statement type: %s\n", stmtType)
						switch stmtType {
						case "Import":
							goCode += handleImport(stmtMap)
						case "FunctionDef":
							goCode += handleFunctionDef(stmtMap)
						case "If":
							goCode += handleIf(stmtMap)
						default:
							stmtJson, _ := json.MarshalIndent(stmtMap, "", "  ")
							fmt.Printf("Unsupported statement type encountered: %s\nDetails: %s\n", stmtType, string(stmtJson))
							goCode += fmt.Sprintf("// Unsupported statement type: %s\n", stmtType)
						}
					} else {
						fmt.Println("Unknown statement type found")
						goCode += "// Unknown statement type\n"
					}
				} else {
					fmt.Println("Unknown statement structure found")
					goCode += "// Unknown statement structure\n"
				}
			}
		} else {
			return "", fmt.Errorf("unsupported AST structure: body not found")
		}
	} else {
		return "", fmt.Errorf("unsupported AST structure")
	}

	goCode += "func main() {\n\tmain()\n}\n"
	return goCode, nil
}

// 处理导入语句
func handleImport(stmtMap map[string]interface{}) string {
	return "// Skipping Python import\n"
}

// 处理函数定义
func handleFunctionDef(stmtMap map[string]interface{}) string {
	functionName := stmtMap["name"].(string)
	code := fmt.Sprintf("func %s() {\n", functionName)

	// 处理函数体
	if body, ok := stmtMap["body"].([]interface{}); ok {
		for _, stmt := range body {
			if stmtMap, ok := stmt.(map[string]interface{}); ok {
				if stmtType, ok := stmtMap["_type"].(string); ok {
					switch stmtType {
					case "Return":
						code += handleReturn(stmtMap)
					case "Expr":
						code += handleExpr(stmtMap)
					case "Assign":
						code += handleAssign(stmtMap)
					case "If":
						code += handleIf(stmtMap)
					default:
						stmtJson, _ := json.MarshalIndent(stmtMap, "", "  ")
						fmt.Printf("Unsupported function statement type encountered: %s\nDetails: %s\n", stmtType, string(stmtJson))
						code += fmt.Sprintf("// Unsupported function statement type: %s\n", stmtType)
					}
				} else {
					fmt.Println("Unknown function statement type found")
					code += "// Unknown function statement type\n"
				}
			} else {
				fmt.Println("Unknown function statement structure found")
				code += "// Unknown function statement structure\n"
			}
		}
	}
	code += "}\n"
	return code
}

// 处理返回语句
func handleReturn(stmtMap map[string]interface{}) string {
	if value, ok := stmtMap["value"].(map[string]interface{}); ok {
		if constValue, ok := value["value"].(string); ok {
			return fmt.Sprintf("\treturn \"%s\"\n", constValue)
		}
	}
	return "\treturn\n"
}

// 处理表达式语句
func handleExpr(stmtMap map[string]interface{}) string {
	if value, ok := stmtMap["value"].(map[string]interface{}); ok {
		if funcCall, ok := value["func"].(map[string]interface{}); ok {
			if funcName, ok := funcCall["id"].(string); ok && funcName == "print" {
				if args, ok := value["args"].([]interface{}); ok && len(args) > 0 {
					if arg, ok := args[0].(map[string]interface{}); ok {
						if constValue, ok := arg["s"].(string); ok {
							return fmt.Sprintf("\tfmt.Println(\"%s\")\n", constValue)
						}
					}
				}
			}
		}
	}
	return "// Skipping unsupported expression\n"
}

// 处理赋值语句
func handleAssign(stmtMap map[string]interface{}) string {
	if targets, ok := stmtMap["targets"].([]interface{}); ok {
		if len(targets) > 0 {
			if target, ok := targets[0].(map[string]interface{}); ok {
				if targetName, ok := target["id"].(string); ok {
					if value, ok := stmtMap["value"].(map[string]interface{}); ok {
						if constValue, ok := value["value"].(string); ok {
							return fmt.Sprintf("\t%s := \"%s\"\n", targetName, constValue)
						}
					}
				}
			}
		}
	}
	return "// Skipping unsupported assignment\n"
}

// 处理条件语句
func handleIf(stmtMap map[string]interface{}) string {
	test := stmtMap["test"].(map[string]interface{})
	body := stmtMap["body"].([]interface{})

	// 假设条件判断为简单的布尔值（此处仅作演示）
	if testValue, ok := test["id"].(string); ok {
		code := fmt.Sprintf("\tif %s {\n", testValue)
		for _, stmt := range body {
			if stmtMap, ok := stmt.(map[string]interface{}); ok {
				if stmtType, ok := stmtMap["_type"].(string); ok {
					switch stmtType {
					case "Return":
						code += handleReturn(stmtMap)
					default:
						code += fmt.Sprintf("\t// Unsupported statement in if body: %s\n", stmtType)
					}
				}
			}
		}
		code += "\t}\n"
		return code
	}

	return "// Skipping unsupported If statement\n"
}

func handleTry(stmtMap map[string]interface{}) string {
	body := stmtMap["body"].([]interface{})
	handlers := stmtMap["handlers"].([]interface{})

	code := "try {\n"
	for _, stmt := range body {
		if stmtMap, ok := stmt.(map[string]interface{}); ok {
			if stmtType, ok := stmtMap["_type"].(string); ok {
				switch stmtType {
				case "Assign":
					code += handleAssign(stmtMap)
				case "If":
					code += handleIf(stmtMap)
				// 添加其他可能出现的结构处理
				default:
					code += fmt.Sprintf("\t// Unsupported statement in try body: %s\n", stmtType)
				}
			}
		}
	}
	code += "} catch(Exception e) {\n"
	for _, handler := range handlers {
		if handlerMap, ok := handler.(map[string]interface{}); ok {
			code += handleExceptHandler(handlerMap)
		}
	}
	code += "}\n"
	return code
}

func handleExceptHandler(handlerMap map[string]interface{}) string {
	// 处理Except的代码
	body := handlerMap["body"].([]interface{})
	code := ""
	for _, stmt := range body {
		if stmtMap, ok := stmt.(map[string]interface{}); ok {
			if stmtType, ok := stmtMap["_type"].(string); ok {
				switch stmtType {
				case "Expr":
					code += handleExpr(stmtMap)
				default:
					code += fmt.Sprintf("\t// Unsupported statement in except body: %s\n", stmtType)
				}
			}
		}
	}
	return code
}
