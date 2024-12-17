package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Install() {
	// 获取当前可执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("✘ 获取程序路径失败: %v\n", err)
		os.Exit(1)
	}

	// 获取目标安装目录
	var installDir string
	if runtime.GOOS == "windows" {
		// Windows 系统通常安装到 C:\Windows\System32
		installDir = os.Getenv("SYSTEMROOT") + "\\System32"
	} else {
		// Linux/MacOS 通常安装到 /usr/local/bin
		installDir = "/usr/local/bin"
	}

	// 获取程序名称
	programName := filepath.Base(execPath)
	targetPath := filepath.Join(installDir, programName)

	// 复制文件到目标目录
	if runtime.GOOS == "windows" {
		// Windows 需要管理员权限
		cmd := exec.Command("cmd", "/C", "copy", execPath, targetPath)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Printf("✘ 安装失败，请确保以管理员权限运行: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Linux/MacOS 使用 cp 命令
		cmd := exec.Command("sudo", "cp", execPath, targetPath)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Printf("✘ 安装失败，请确保有 sudo 权限: %v\n", err)
			os.Exit(1)
		}

		// 设置可执行权限
		cmd = exec.Command("sudo", "chmod", "+x", targetPath)
		if err := cmd.Run(); err != nil {
			fmt.Printf("✘ 设置执行权限失败: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("✔ 程序已成功安装到: %s\n", targetPath)
}
