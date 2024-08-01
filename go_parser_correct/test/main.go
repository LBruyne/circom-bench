package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// 文件处理函数
func processFile(filePath string, done chan bool) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read file %s: %v", filePath, err)
		done <- false
		return
	}
	fmt.Printf("Processing file %s, size: %d bytes\n", filePath, len(data))
	done <- true
}

func main() {
	fileChannel := make(chan string)
	done := make(chan bool)
	exit := make(chan os.Signal, 1)

	// 捕捉系统中断信号
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 处理文件
	go func() {
		for {
			select {
			case filePath := <-fileChannel:
				processFile(filePath, done)
			case <-exit:
				fmt.Println("Shutting down...")
				return
			}
		}
	}()

	// 读取用户输入的文件路径
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter file path: ")
			filePath, _ := reader.ReadString('\n')
			filePath = strings.TrimSpace(filePath)
			if filePath == "exit" {
				close(exit)
				return
			}
			fileChannel <- filePath
			success := <-done
			if success {
				fmt.Printf("File %s processed successfully.\n", filePath)
			} else {
				fmt.Printf("Failed to process file %s.\n", filePath)
			}
		}
	}()

	// 挂起主程序，等待中断信号
	<-exit
	fmt.Println("Program terminated.")
}
