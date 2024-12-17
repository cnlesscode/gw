package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cnlesscode/goWatcher/tools"
	"github.com/fsnotify/fsnotify"
)

func main() {
	// 安装到系统目录
	if len(os.Args) > 1 && os.Args[1] == "install" {
		tools.Install()
		return
	}

	// 文件监听
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("✘ 创建监听器失败:", err)
	}
	defer watcher.Close()

	// 启动监听协程
	go func() {
		// 防止重复编译
		lastBuild := time.Now()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 重命名
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					watcher.Add(event.Name)
				}

				// 新建目录
				if event.Op&fsnotify.Create == fsnotify.Create {
					// 添加监听
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						watcher.Add(event.Name)
						log.Printf("› 新增目录监听: %s\n", event.Name)
					}
				}

				// 只处理 .go 文件的变化
				if !strings.HasSuffix(event.Name, ".go") {
					continue
				}

				// 文件保存事件
				if event.Op&fsnotify.Write == fsnotify.Write {
					// 防止短时间内多次触发
					if time.Since(lastBuild).Seconds() < 1 {
						continue
					}
					lastBuild = time.Now()
					log.Printf("› 检测到文件变化: %s\n", event.Name)

					// 终止之前的进程
					tools.KillProcess()

					// 重新编译并运行
					tools.BuildAndRun()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("✘ 监听错误:", err)
			}
		}
	}()

	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 添加当前目录监听
	watcher.Add(dir)
	// 监听子目录
	// 递归监听当前目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// 排除一些不需要监听的目录
			if tools.ShouldIgnoreDir(path) {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}
		return nil
	})

	log.Println("› 开始监听文件变化 ...")

	// 首次运行
	tools.BuildAndRun()
	for {
		time.Sleep(time.Second * 10)
	}
}
