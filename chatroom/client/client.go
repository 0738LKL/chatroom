package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin) // 从标准输入读取用户输入

	// 读取用户网名
	fmt.Print("请输入网名: ")
	nickname, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取网名失败:", err)
		return
	}
	nickname = strings.TrimSpace(nickname)

	// 发送网名到服务器
	_, err = conn.Write([]byte(nickname + "\n"))
	if err != nil {
		fmt.Println("发送网名失败:", err)
		return
	}

	// 启动 goroutine 接收服务器消息
	go receiveMessages(conn)

	// 循环读取用户输入，发送消息到服务器
	for {
		line, err := reader.ReadString('\n') // 读取用户输入，直到遇到换行符
		if err != nil {
			fmt.Println("读取用户输入失败:", err)
			return
		}

		line = strings.TrimSpace(line) // 消除 line 中的空格和换行符
		if line == "exit" {
			fmt.Println("客户端退出")
			break
		}

		// 将 line 发送给服务器
		_, err = conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Println("发送数据失败:", err)
			return
		}
	}
}

func receiveMessages(conn net.Conn) { // 接收服务器消息

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("接收消息失败:", err)
			return
		}
		fmt.Print(message)
	}
}
