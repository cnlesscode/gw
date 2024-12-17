package tools

import (
	"log"
	"os"
)

var CurrentProcess *os.Process

func KillProcess() {
	if CurrentProcess != nil {
		log.Println("› 终止旧进程 ...")
		CurrentProcess.Kill()
		CurrentProcess = nil
	}
}
