package main

import (
	"enhanced_python_compiler/internal/generator"
	"enhanced_python_compiler/internal/parser"
	"enhanced_python_compiler/internal/translator"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: compiler <script.py>")
		return
	}

	filename := os.Args[1]

	// 读取 Python 脚本
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading script: %v\n", err)
		return
	}

	// 解析 Python 代码
	ast, err := parser.ParsePythonCode(code)
	if err != nil {
		fmt.Printf("Error parsing script: %v\n", err)
		return
	}

	// 将 AST 转换为 Go 代码
	goCode, err := translator.TranslateASTToGo(ast)
	if err != nil {
		fmt.Printf("Error translating AST to Go: %v\n", err)
		return
	}

	// 保存无法转换的部分为临时 Python 文件
	err = ioutil.WriteFile("temp_unconverted.py", code, 0644)
	if err != nil {
		fmt.Printf("Error writing unconverted Python code: %v\n", err)
		return
	}

	// 调用 PyInstaller 打包无法转换的部分
	err = packagePythonCode()
	if err != nil {
		fmt.Printf("Error packaging Python code: %v\n", err)
		return
	}

	// 生成最终的 Go 代码文件
	err = generator.GenerateGoCode(goCode, string(code))
	if err != nil {
		fmt.Printf("Error generating Go code: %v\n", err)
		return
	}

	fmt.Println("Go code has been successfully generated in output.go")
}

func packagePythonCode() error {
	// 使用 PyInstaller 打包无法转换的 Python 代码
	cmd := exec.Command("pyinstaller", "--onefile", "temp_unconverted.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error packaging Python code: %v", err)
	}

	// 将打包后的可执行文件嵌入到 Go 程序中
	err = exec.Command("go-bindata", "-o", "bindata.go", "dist/temp_unconverted").Run()
	if err != nil {
		return fmt.Errorf("Error embedding Python executable: %v", err)
	}

	return nil
}
