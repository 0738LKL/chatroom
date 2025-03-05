package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]bool) // 存储所有连接的客户端
	clientsMu sync.Mutex                // 保护 clients 的互斥锁
	broadcast = make(chan string)       // 广播消息的通道
)

func main() {

	fmt.Println("服务器开始监听...")
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("监听失败:", err)
		return
	}
	defer listen.Close()

	// 启动广播消息的 goroutine
	go broadcaster()

	// 循环等待客户端连接
	for {
		fmt.Println("等待客户端连接...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("连接失败:", err)
			continue
		}

		clientsMu.Lock()
		clients[conn] = true
		clientsMu.Unlock()

		fmt.Printf("连接成功: %v 客户端IP：%v\n", conn, conn.RemoteAddr().String())

		// 创建一个 goroutine 处理连接
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 读取客户端的网名
	reader := bufio.NewReader(conn)
	nickname, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取网名失败:", err)
		return
	}
	nickname = strings.TrimSpace(nickname) //

	// 广播新用户加入的消息
	broadcast <- fmt.Sprintf("%s 加入了聊天室\n", nickname)

	// 循环接受客户端的数据
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("客户端退出:", err)
			break
		}

		message := strings.TrimSpace(string(buf[:n]))
		if message == "exit" {
			broadcast <- fmt.Sprintf("%s 离开了聊天室\n", nickname)
			break
		}

		// 广播消息
		broadcast <- fmt.Sprintf("%s: %s\n", nickname, message)
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
}

func broadcaster() {

	for {
		message := <-broadcast        // 从广播消息通道中读取消息
		clientsMu.Lock()              // 锁定客户端列表
		for client := range clients { // 遍历客户端列表
			_, err := client.Write([]byte(message)) // 发送消息给客户端
			if err != nil {
				fmt.Println("发送消息失败:", err)
				client.Close()
				delete(clients, client) // 删除客户端
			}
		}
		clientsMu.Unlock()
	}
}
