package generator

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// goTemplate 定义了生成的 Go 代码的模板
const goTemplate = `
package main

import (
	"fmt"
	"enhanced_python_compiler/internal/runtime"
)

func main() {
	{{.GoCode}}

	// 嵌入并执行 Python 代码
	result, err := runtime.ExecuteCPython({{.PythonCode}})
	if err != nil {
		fmt.Println("Error executing Python code:", err)
	} else {
		fmt.Println(result)
	}
}
`

// GenerateGoCode 生成最终的 Go 代码文件
func GenerateGoCode(goCode, pythonCode string) error {
	tmpl, err := template.New("goCode").Parse(goTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse Go template: %v", err)
	}

	outputFile, err := os.Create("output.go")
	if err != nil {
		return fmt.Errorf("failed to create output Go file: %v", err)
	}
	defer outputFile.Close()

	// 清理并移除 Go 代码中的多余结构
	cleanGoCode := cleanUpGeneratedGoCode(goCode)

	data := map[string]string{
		"GoCode":     cleanGoCode,
		"PythonCode": fmt.Sprintf("%q", pythonCode),
	}

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute Go template: %v", err)
	}

	fmt.Println("Go code has been successfully generated in output.go")
	return nil
}

// cleanUpGeneratedGoCode 清理生成的 Go 代码，移除多余的 package 和 import 结构
func cleanUpGeneratedGoCode(goCode string) string {
	// 移除不必要的 package 声明
	cleanedCode := strings.Replace(goCode, "package main", "", -1)

	// 移除不必要的 import 语句
	cleanedCode = strings.Replace(cleanedCode, `import "fmt"`, "", -1)

	// 移除不必要的 func main() 结构
	cleanedCode = strings.Replace(cleanedCode, "func main()", "", -1)

	// 移除多余的花括号
	cleanedCode = strings.TrimSpace(cleanedCode)
	if strings.HasPrefix(cleanedCode, "{") && strings.HasSuffix(cleanedCode, "}") {
		cleanedCode = strings.TrimPrefix(cleanedCode, "{")
		cleanedCode = strings.TrimSuffix(cleanedCode, "}")
		cleanedCode = strings.TrimSpace(cleanedCode)
	}

	return cleanedCode
}
