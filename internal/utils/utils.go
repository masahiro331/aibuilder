package utils

import (
    "bytes"
    "os/exec"
    "strings"
)

func ExecuteCommand(command string) (string, error) {
    cmd := exec.Command("bash", "-c", command)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out

    err := cmd.Run()
    if err != nil {
        return "", err
    }
    return out.String(), nil
}

func IsSafeCommand(command string) bool {
    // 許可されたコマンドのみ実行
    allowedCommands := []string{"mkdir", "cat", "echo"}
    for _, cmd := range allowedCommands {
        if strings.HasPrefix(strings.TrimSpace(command), cmd) {
            return true
        }
    }
    return false
}
