package runtime

import (
	"fmt"
	"os"
	"os/exec"
)

func ExecuteGo(goCode string) (string, error) {
	// 這裡我們可以先將 Go 代碼寫入臨時文件，然後使用 `go run` 執行
	tmpFile := "/tmp/temp_go_code.go"
	err := os.WriteFile(tmpFile, []byte(goCode), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write Go code to file: %v", err)
	}

	cmd := exec.Command("go", "run", tmpFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute Go code: %v", err)
	}

	return string(output), nil
}
