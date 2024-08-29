package translator

import (
	"encoding/json"
	"enhanced_python_compiler/internal/parser"
	"fmt"
	"strings"
)

// TranslateASTToGo 將 Python 的 AST 轉換為 Go 代碼
func TranslateASTToGo(ast *parser.AST) (string, error) {
	var rootMap map[string]interface{}

	// 將 AST 的根節點轉換為 Go 的 map 結構
	rootJson, err := json.Marshal(ast.Root)
	if err != nil {
		return "", fmt.Errorf("failed to marshal AST root: %v", err)
	}

	err = json.Unmarshal(rootJson, &rootMap)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal AST root: %v", err)
	}

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

func TranslateASTToBinary(ast *parser.AST) ([]byte, bool, error) {
	goCode, err := TranslateASTToGo(ast)
	if err != nil {
		return nil, false, fmt.Errorf("error translating AST to Go: %v", err)
	}
	return []byte(goCode), false, nil
}
