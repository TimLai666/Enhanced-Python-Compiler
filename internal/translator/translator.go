package translator

import (
	"encoding/json"
	"enhanced_python_compiler/internal/parser"
	"fmt"
	"strings"
)

// TranslateASTToGo 將 Python 的 AST 轉換為 Go 代碼
func TranslateASTToGo(rootMap map[string]interface{}) (string, error) {
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
					case "Assign":
						goCode += handleAssign(stmtMap)
					case "Return":
						goCode += handleReturn(stmtMap)
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

// 處理 Import 語句
func handleImport(stmtMap map[string]interface{}) string {
	return "// Skipping Python import\n"
}

// 處理函數定義
func handleFunctionDef(stmtMap map[string]interface{}) string {
	functionName := stmtMap["name"].(string)
	code := fmt.Sprintf("func %s() {\n", functionName)

	// 處理函數體
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
					case "Assign":
						code += handleAssign(stmtMap)
					case "Return":
						code += handleReturn(stmtMap)
					case "Try":
						code += handleTry(stmtMap) // 添加這一行
					case "For":
						code += handleFor(stmtMap) // 添加這一行
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

func processStatement(stmtMap map[string]interface{}) string {
	if stmtType, ok := stmtMap["_type"].(string); ok {
		switch stmtType {
		case "Assign":
			return handleAssign(stmtMap)
		case "Return":
			return handleReturn(stmtMap)
		case "If":
			return handleIf(stmtMap)
		// 添加其他語句類型的處理
		default:
			stmtJson, _ := json.MarshalIndent(stmtMap, "", "  ")
			return fmt.Sprintf("// Unsupported statement type: %s\nDetails: %s\n", stmtType, string(stmtJson))
		}
	}
	return "// Unknown statement type\n"
}

// 處理返回語句
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
		}
	}
	return "\treturn\n"
}

// 處理表達式語句
func handleExpr(stmtMap map[string]interface{}) string {
	if value, ok := stmtMap["value"].(map[string]interface{}); ok {
		switch valueType := value["_type"].(string); valueType {
		case "Call":
			if funcCall, ok := value["func"].(map[string]interface{}); ok {
				if funcName, ok := funcCall["id"].(string); ok && funcName == "print" {
					if args, ok := value["args"].([]interface{}); ok && len(args) > 0 {
						var printArgs []string
						for _, arg := range args {
							if argMap, ok := arg.(map[string]interface{}); ok {
								if constValue, ok := argMap["value"].(string); ok {
									printArgs = append(printArgs, fmt.Sprintf("\"%s\"", constValue))
								} else if varName, ok := argMap["id"].(string); ok {
									printArgs = append(printArgs, varName)
								}
							}
						}
						return fmt.Sprintf("\tfmt.Println(%s)\n", strings.Join(printArgs, ", "))
					}
				}
			}
		}
	}
	return "// Skipping unsupported expression\n"
}

// 處理賦值語句
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
						case "Call":
							if funcCall, ok := value["func"].(map[string]interface{}); ok {
								if funcName, ok := funcCall["id"].(string); ok {
									return fmt.Sprintf("\t%s := %s()\n", targetName, funcName)
								}
							}
						case "Name":
							if varName, ok := value["id"].(string); ok {
								return fmt.Sprintf("\t%s := %s\n", targetName, varName)
							}
						case "List":
							return fmt.Sprintf("\t%s := []interface{}{}\n", targetName)
						}
					}
				}
			}
		}
	}
	return "// Skipping unsupported assignment\n"
}

// 處理條件語句
func handleIf(stmtMap map[string]interface{}) string {
	test := stmtMap["test"].(map[string]interface{})
	body := stmtMap["body"].([]interface{})

	// 假設條件判斷為簡單的布爾值（此處僅作演示）
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
	// 基本模板：Go 中沒有原生的 Try-Catch 結構，因此需要用 defer 和 recover 模擬
	code := "defer func() {\n"
	code += "\tif r := recover(); r != nil {\n"
	// 這裡可以加上處理 `handlers` 部分的邏輯
	code += "\t\tfmt.Println(\"Recovered from error\")\n"
	code += "\t}\n"
	code += "}()\n"

	if body, ok := stmtMap["body"].([]interface{}); ok {
		for _, stmt := range body {
			if stmtMap, ok := stmt.(map[string]interface{}); ok {
				code += processStatement(stmtMap) // 假設 processStatement 是您用來處理各種語句的通用函數
			}
		}
	}
	return code
}

func handleFor(stmtMap map[string]interface{}) string {
	// 基本模板：將 Python 的 For 轉換為 Go 的 for range 結構
	iter := stmtMap["iter"].(map[string]interface{})["id"].(string)
	target := stmtMap["target"].(map[string]interface{})["id"].(string)

	code := fmt.Sprintf("for _, %s := range %s {\n", target, iter)
	if body, ok := stmtMap["body"].([]interface{}); ok {
		for _, stmt := range body {
			if stmtMap, ok := stmt.(map[string]interface{}); ok {
				code += processStatement(stmtMap) // 假設 processStatement 是您用來處理各種語句的通用函數
			}
		}
	}
	code += "}\n"
	return code
}

func TranslateASTToBinary(ast *parser.AST) ([]byte, bool, error) {
	var rootMap map[string]interface{}

	// 将 ast.Root（JSON 字符串）转换为 Go 结构
	if err := json.Unmarshal([]byte(ast.Root), &rootMap); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal AST root: %v", err)
	}

	// 此處進行相應的轉換處理，例如將 AST 轉換為 Go 代碼
	goCode, err := TranslateASTToGo(rootMap)
	if err != nil {
		return nil, false, fmt.Errorf("error translating AST to Go: %v", err)
	}

	// 返回轉換後的 Go 代碼
	return []byte(goCode), true, nil
}
