package translator

import (
	"encoding/json"
	"enhanced_python_compiler/internal/parser"
	"fmt"
	"strings"
)

func TranslateASTToGo(rootMap map[string]interface{}) (string, error) {
	goCode := "package main\n\nimport \"fmt\"\n\n"

	mainDefined := false // 用於避免重複定義 main 函數

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
						funcCode := handleFunctionDef(stmtMap)
						if strings.Contains(funcCode, "func main()") {
							if !mainDefined {
								goCode += funcCode
								mainDefined = true
							} else {
								fmt.Println("Skipping additional main function.")
							}
						} else {
							goCode += funcCode
						}
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

	// 確保只有一個 main 函數
	if !mainDefined {
		goCode += generateMainFunction()
	}

	fmt.Println("Generated Go code:")
	fmt.Println(goCode)

	return goCode, nil
}

func handleImport(stmtMap map[string]interface{}) string {
	return "// Skipping Python import\n"
}

func handleFunctionDef(stmtMap map[string]interface{}) string {
	functionName := stmtMap["name"].(string)
	code := fmt.Sprintf("func %s() {\n", functionName)

	if body, ok := stmtMap["body"].([]interface{}); ok {
		for _, stmt := range body {
			if stmtMap, ok := stmt.(map[string]interface{}); ok {
				if stmtType, ok := stmtMap["_type"].(string); ok {
					switch stmtType {
					case "Assign":
						code += handleAssign(stmtMap)
					case "Return":
						code += handleReturn(stmtMap)
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

	// 針對函數名稱為 main 的特殊處理，避免重複定義
	if functionName == "main" {
		return generateMainFunction()
	}

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
		default:
			stmtJson, _ := json.MarshalIndent(stmtMap, "", "  ")
			return fmt.Sprintf("// Unsupported statement type: %s\nDetails: %s\n", stmtType, string(stmtJson))
		}
	}
	return "// Unknown statement type\n"
}

func handleReturn(stmtMap map[string]interface{}) string {
	if value, ok := stmtMap["value"]; ok && value != nil {
		return "\t_ = response\n" // 假設返回值被存儲在 response 中
	}
	return ""
}

func handleAssign(stmtMap map[string]interface{}) string {
	targetName := "_"
	if targets, ok := stmtMap["targets"].([]interface{}); ok && len(targets) > 0 {
		if target, ok := targets[0].(map[string]interface{}); ok {
			if tn, ok := target["id"].(string); ok {
				targetName = tn
			}

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
	return fmt.Sprintf("\tvar %s interface{}\n_ = %s\n", targetName, targetName)
}

func handleIf(stmtMap map[string]interface{}) string {
	test := stmtMap["test"].(map[string]interface{})
	body := stmtMap["body"].([]interface{})

	if testValue, ok := test["id"].(string); ok {
		code := fmt.Sprintf("\tif %s {\n", testValue)
		for _, stmt := range body {
			if stmtMap, ok := stmt.(map[string]interface{}); ok {
				code += processStatement(stmtMap)
			}
		}
		code += "\t}\n"
		return code
	}

	return "// Skipping unsupported If statement\n"
}

func handleTry(stmtMap map[string]interface{}) string {
	code := "defer func() {\n"
	code += "\tif r := recover(); r != nil {\n"
	code += "\t\tfmt.Println(\"Recovered from error\")\n"
	code += "\t}\n"
	code += "}()\n"

	if body, ok := stmtMap["body"].([]interface{}); ok {
		for _, stmt := range body {
			if stmtMap, ok := stmt.(map[string]interface{}); ok {
				code += processStatement(stmtMap)
			}
		}
	}
	return code
}

func handleFor(stmtMap map[string]interface{}) string {
	iter := stmtMap["iter"].(map[string]interface{})["id"].(string)
	target := stmtMap["target"].(map[string]interface{})["id"].(string)

	code := fmt.Sprintf("for _, %s := range %s {\n", target, iter)
	code += fmt.Sprintf("\t_ = %s\n", target)

	if body, ok := stmtMap["body"].([]interface{}); ok {
		for _, stmt := range body {
			if stmtMap, ok := stmt.(map[string]interface{}); ok {
				code += processStatement(stmtMap)
			}
		}
	}
	code += "}\n"
	return code
}

func generateMainFunction() string {
	return "func main() {\n\t// 程序入口點\n}\n"
}

func TranslateASTToBinary(ast *parser.AST) ([]byte, bool, error) {
	var rootMap map[string]interface{}

	if err := json.Unmarshal([]byte(ast.Root), &rootMap); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal AST root: %v", err)
	}

	goCode, err := TranslateASTToGo(rootMap)
	if err != nil {
		return nil, false, fmt.Errorf("error translating AST to Go: %v", err)
	}

	return []byte(goCode), true, nil
}
