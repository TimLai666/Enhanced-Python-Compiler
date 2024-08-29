package packager

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// PackageExecutable 將 Go 代碼、Python 執行時和依賴包打包成一個可執行文件
func PackageExecutable(goCode string) error {
	// 1. 創建臨時目錄來存儲打包過程中的文件
	tempDir, err := os.MkdirTemp("", "compiler-packager-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // 打包完成後刪除臨時目錄

	// 2. 將 Go 代碼寫入臨時目錄中的 main.go 文件
	goFilePath := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(goFilePath, []byte(goCode), 0644); err != nil {
		return fmt.Errorf("failed to write Go code to file: %v", err)
	}

	// 3. 使用 `go build` 將 Go 代碼編譯成可執行文件
	executablePath := filepath.Join(tempDir, "output_executable")
	cmd := exec.Command("go", "build", "-o", executablePath, goFilePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to build Go executable: %v\nDetails: %s", err, out.String())
	}

	// 4. 打包 Python 執行時、依賴包和 Go 可執行文件
	finalExecutablePath := filepath.Join(".", "final_executable")
	err = packageWithPythonRuntime(executablePath, finalExecutablePath)
	if err != nil {
		return fmt.Errorf("failed to package final executable: %v", err)
	}

	fmt.Println("Executable built successfully:", finalExecutablePath)
	return nil
}

// packageWithPythonRuntime 將 Go 可執行文件與 Python 執行時和依賴包一起打包
func packageWithPythonRuntime(goExecutablePath, finalExecutablePath string) error {
	// 這裡假設你已經準備好 Python 執行時和依賴包，並且可以使用 `go-bindata` 或類似工具打包
	cmd := exec.Command("go-bindata", "-o", finalExecutablePath, goExecutablePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to package with Python runtime: %v\nDetails: %s", err, out.String())
	}
	return nil
}
