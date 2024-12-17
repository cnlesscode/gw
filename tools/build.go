package tools

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func BuildAndRun() {

	var exeSuffix string
	switch runtime.GOOS {
	case "windows":
		exeSuffix = ".exe"
	case "darwin", "linux":
		exeSuffix = ""
	default:
		log.Println("✘ 不支持的操作系统")
		return
	}

	// 编译
	buildCmd := exec.Command("go", "build", "-o", "main"+exeSuffix, "main.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	log.Println("› 正在编译...")
	if err := buildCmd.Run(); err != nil {
		log.Println("✘ 编译失败:", err)
		return
	}

	log.Println("› 编译成功，正在运行...")
	// 运行编译后的程序
	cmd := exec.Command("./main" + exeSuffix)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Println("启动程序失败:", err)
		return
	}

	CurrentProcess = cmd.Process
}
