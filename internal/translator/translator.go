package translator

import (
	"encoding/json"
	"enhanced_python_compiler/internal/parser"
	"fmt"
)

// TranslateASTToGo 将 Python 的 AST 转换为 Go 代码
func TranslateASTToGo(ast *parser.AST) (string, error) {
	// 這裡 ast.Root 已經是一個 map[string]interface{} 結構，無需再次進行反序列化
	rootMap := ast.Root.(map[string]interface{})

	goCode := "package main\n\nimport \"fmt\"\n\n"

	// 處理 AST 的頂層 Module 節點
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
					case "Try":
						goCode += handleTry(stmtMap)
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
					case "Import":
						code += handleImport(stmtMap)
					case "FunctionDef":
						code += handleFunctionDef(stmtMap)
					case "If":
						code += handleIf(stmtMap)
					case "Try":
						code += handleTry(stmtMap)
					case "For":
						code += handleFor(stmtMap)
					default:
						stmtJson, _ := json.MarshalIndent(stmtMap, "", "  ")
						fmt.Printf("Unsupported statement type encountered: %s\nDetails: %s\n", stmtType, string(stmtJson))
						code += fmt.Sprintf("// Unsupported statement type: %s\n", stmtType)
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

// 处理表达式语句
func handleExpr(stmtMap map[string]interface{}) string {
	if value, ok := stmtMap["value"].(map[string]interface{}); ok {
		switch valueType := value["_type"].(string); valueType {
		case "Call":
			if funcCall, ok := value["func"].(map[string]interface{}); ok {
				if funcName, ok := funcCall["id"].(string); ok && funcName == "print" {
					if args, ok := value["args"].([]interface{}); ok && len(args) > 0 {
						if arg, ok := args[0].(map[string]interface{}); ok {
							if constValue, ok := arg["value"].(string); ok {
								return fmt.Sprintf("\tfmt.Println(\"%s\")\n", constValue)
							}
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
						switch valueType := value["_type"].(string); valueType {
						case "Constant":
							if constValue, ok := value["value"].(string); ok {
								return fmt.Sprintf("\t%s := \"%s\"\n", targetName, constValue)
							}
						case "Name":
							if varName, ok := value["id"].(string); ok {
								return fmt.Sprintf("\t%s := %s\n", targetName, varName)
							}
						case "List":
							return fmt.Sprintf("\t%s := []interface{}{}\n", targetName)
						case "Call":
							if funcCall, ok := value["func"].(map[string]interface{}); ok {
								if funcName, ok := funcCall["id"].(string); ok {
									argsStr := handleCallArguments(value["args"].([]interface{}))
									return fmt.Sprintf("\t%s := %s(%s)\n", targetName, funcName, argsStr)
								}
							}
						}
					}
				}
			}
		}
	}
	return "// Skipping unsupported assignment\n"
}

// 处理返回语句
func handleReturn(stmtMap map[string]interface{}) string {
	if value, ok := stmtMap["value"].(map[string]interface{}); ok {
		switch valueType := value["_type"].(string); valueType {
		case "Constant":
			if constValue, ok := value["value"].(string); ok {
				return fmt.Sprintf("\treturn \"%s\"\n", constValue)
			}
		case "Name":
			if varName, ok := value["id"].(string); ok {
				return fmt.Sprintf("\treturn %s\n", varName)
			}
		case "Call":
			if funcCall, ok := value["func"].(map[string]interface{}); ok {
				if funcName, ok := funcCall["id"].(string); ok {
					argsStr := handleCallArguments(value["args"].([]interface{}))
					return fmt.Sprintf("\treturn %s(%s)\n", funcName, argsStr)
				}
			}
		}
	}
	return "\treturn\n"
}

// 处理函数调用的参数
func handleCallArguments(args []interface{}) string {
	var argsList []string
	for _, arg := range args {
		if argMap, ok := arg.(map[string]interface{}); ok {
			if argType, ok := argMap["_type"].(string); ok {
				switch argType {
				case "Constant":
					if constValue, ok := argMap["value"].(string); ok {
						argsList = append(argsList, fmt.Sprintf("\"%s\"", constValue))
					}
				case "Name":
					if varName, ok := argMap["id"].(string); ok {
						argsList = append(argsList, varName)
					}
				}
			}
		}
	}
	return fmt.Sprintf("%s", argsList)
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

// 处理 Try 语句
func handleTry(stmtMap map[string]interface{}) string {
	body := stmtMap["body"].([]interface{})
	handlers := stmtMap["handlers"].([]interface{})

	code := "\ttry {\n"
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
					code += fmt.Sprintf("\t\t// Unsupported statement in try body: %s\n", stmtType)
				}
			}
		}
	}
	code += "\t} catch(Exception e) {\n"
	for _, handler := range handlers {
		if handlerMap, ok := handler.(map[string]interface{}); ok {
			code += handleExceptHandler(handlerMap)
		}
	}
	code += "\t}\n"
	return code
}

// 处理 ExceptHandler 语句
func handleExceptHandler(handlerMap map[string]interface{}) string {
	body := handlerMap["body"].([]interface{})
	code := ""
	for _, stmt := range body {
		if stmtMap, ok := stmt.(map[string]interface{}); ok {
			if stmtType, ok := stmtMap["_type"].(string); ok {
				switch stmtType {
				case "Expr":
					code += handleExpr(stmtMap)
				default:
					code += fmt.Sprintf("\t\t// Unsupported statement in except body: %s\n", stmtType)
				}
			}
		}
	}
	return code
}

// 处理 For 语句
func handleFor(stmtMap map[string]interface{}) string {
	target := stmtMap["target"].(map[string]interface{})
	iter := stmtMap["iter"].(map[string]interface{})
	body := stmtMap["body"].([]interface{})

	// 假设 target 是简单变量，iter 是简单变量
	targetName := target["id"].(string)
	iterName := iter["id"].(string)

	code := fmt.Sprintf("for _, %s := range %s {\n", targetName, iterName)

	for _, stmt := range body {
		if stmtMap, ok := stmt.(map[string]interface{}); ok {
			if stmtType, ok := stmtMap["_type"].(string); ok {
				switch stmtType {
				case "If":
					code += handleIf(stmtMap)
				case "Expr":
					code += handleExpr(stmtMap)
				default:
					stmtJson, _ := json.MarshalIndent(stmtMap, "", "  ")
					fmt.Printf("Unsupported statement type in For body: %s\nDetails: %s\n", stmtType, string(stmtJson))
					code += fmt.Sprintf("\t// Unsupported statement type: %s\n", stmtType)
				}
			}
		}
	}

	code += "}\n"
	return code
}
