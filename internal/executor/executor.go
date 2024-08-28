package executor

import (
	"enhanced_python_compiler/internal/runtime"
	"fmt"
	"io/ioutil"
)

func Execute(filePath string) (string, error) {
	// 讀取 Python 腳本
	sourceCode, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Error reading file: %v", err)
	}

	// 使用 CPython 執行 Python 腳本
	return runtime.ExecuteCPython(string(sourceCode))
}
