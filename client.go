package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <IP> <protocol> <ports> <protocol> <ports> ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "You must specify an IP, protocol and at least one port for each protocol.\n")
		fmt.Fprintf(os.Stderr, "Usage: ./port-client IP udp PORT,PORT tcp PORT,PORT\n")
		os.Exit(1)
	}

	ip := os.Args[1]
	var wg sync.WaitGroup

	// 从命令行中解析协议和端口
	for i := 2; i < len(os.Args); i += 2 {
		protocol := os.Args[i]
		ports := os.Args[i+1]

		// 解析端口号，使用逗号分隔
		portStrings := strings.Split(ports, ",")
		var portNumbers []int
		for _, portStr := range portStrings {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Fatalf("Invalid port number: %s\n", portStr)
			}
			portNumbers = append(portNumbers, port)
		}

		// 启动指定协议的处理
		if protocol == "udp" {
			for _, port := range portNumbers {
				wg.Add(1)
				go handleConnection("udp", ip, port, &wg)
			}
		} else if protocol == "tcp" {
			for _, port := range portNumbers {
				wg.Add(1)
				go handleConnection("tcp", ip, port, &wg)
			}
		} else {
			log.Fatalf("Unknown protocol: %s. Only 'udp' and 'tcp' are supported.\n", protocol)
		}
	}

	// 等待所有 goroutines 完成
	wg.Wait()

	fmt.Println("Server exiting...")
}

// 处理 UDP 和 TCP 连接
func handleConnection(proto, ip string, port int, wg *sync.WaitGroup) {
	defer wg.Done()

	nameport := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.Dial(proto, nameport)
	if err != nil {
		// log.Printf("%s connection to %s failed: %v\n", proto, nameport, err)
		fmt.Printf("[ACK-%s:%v]-[STATUS: FAILED]\n", strings.ToUpper(proto), port)
		return
	}

	// 设置读取超时时间（2秒）
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	// 发送消息
	sendMsgStr := "PING"
	sendMsg := []byte(sendMsgStr)
	_, wrerr := conn.Write(sendMsg)

	if wrerr != nil {
		fmt.Printf("conn.Write() error: %s\n", wrerr)
	} else {
		// fmt.Printf("[REQ-%s-PORT: %v]-[MSG: %s]-[LOCAL-ADDR: %v]=>[REMOTE-ADDR: %v]\n", proto, port, sendMsgStr, conn.LocalAddr(), conn.RemoteAddr())
		// 接收响应
		recvMsg := make([]byte, 2048)
		cc, rderr := conn.Read(recvMsg)
		if rderr != nil {
			// 超时处理
			if opErr, ok := rderr.(net.Error); ok && opErr.Timeout() {
				fmt.Printf("[ACK-%s:%v]-[STATUS: FAILED]\n", strings.ToUpper(proto), port)
			} else {
				fmt.Printf("[ACK-%s:%v]-[STATUS: FAILED]\n", strings.ToUpper(proto), port)
			}
		} else {
			recvMsgStr := string(recvMsg[:cc])
			fmt.Printf("[ACK-%s:%v]-[STATUS: CONNECTED]-[MSG: %s]\n", strings.ToUpper(proto), port, recvMsgStr)
		}
	}

	// 关闭连接
	if err = conn.Close(); err != nil {
		log.Fatal(err)
	}
}
