package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s - UDP/TCP echo server\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s <IP> <protocol> <ports> <protocol> <ports> ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "You have to specify an IP address, protocol (udp/tcp) and at least one port number for each protocol\n")
		fmt.Fprintf(os.Stderr, "Usage: ./port-server IP udp PORT,PORT tcp PORT,PORT\n")
		os.Exit(1)
	}

	ip := os.Args[1]
	var wg sync.WaitGroup

	// 解析命令行参数
	for i := 2; i < len(os.Args); i += 2 {
		protocol := os.Args[i]
		portStrs := os.Args[i+1]

		// 解析端口号（以逗号分隔）
		ports := strings.Split(portStrs, ",")
		var portNumbers []int
		for _, portStr := range ports {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Fatalf("Invalid port number: %s\n", portStr)
			}
			portNumbers = append(portNumbers, port)
		}

		// 启动对应协议的监听
		if protocol == "udp" {
			for _, port := range portNumbers {
				addr := net.UDPAddr{Port: port, IP: net.ParseIP(ip)}
				conn, err := net.ListenUDP("udp", &addr)
				if err != nil {
					log.Fatal(err)
				}
				wg.Add(1)
				go handleUDPConnection(conn, port, &wg)
			}
		} else if protocol == "tcp" {
			for _, port := range portNumbers {
				addr := net.TCPAddr{Port: port, IP: net.ParseIP(ip)}
				listen, err := net.ListenTCP("tcp", &addr)
				if err != nil {
					log.Fatal(err)
				}
				wg.Add(1)
				go handleTCPConnection(listen, port, &wg)
			}
		} else {
			log.Fatalf("Unknown protocol: %s. Only 'udp' and 'tcp' are supported.\n", protocol)
		}
	}

	// 等待所有的 goroutines 完成
	wg.Wait()

	fmt.Println("Server exiting...")
}

func handleUDPConnection(conn *net.UDPConn, port int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[UDP] Listening on port %d\n", port)

	b := make([]byte, 2048)

	for {
		// 等待接收数据
		cc, remote, rderr := conn.ReadFromUDP(b)
		if rderr != nil {
			fmt.Printf("net.ReadFromUDP() error: %s\n", rderr)
			return
		}

		// 处理收到的消息
		reqMsg := string(b[:cc])
		fmt.Printf("[REQ-UDP-PORT: %v]-[MSG: %q]-[REMOTE-ADDR: %v]\n", port, reqMsg, remote)

		// 根据不同端口返回不同的消息
		resMsg := "PONG"

		cc, wrerr := conn.WriteTo([]byte(resMsg), remote)
		if wrerr != nil {
			fmt.Printf("net.WriteTo() error: %s\n", wrerr)
		} else {
			fmt.Printf("%s\n", resMsg)
			fmt.Println("---------------------------------")
		}
	}
}

func handleTCPConnection(listen *net.TCPListener, port int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[TCP] Listening on TCP port %d\n", port)

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Printf("AcceptTCP() error: %s\n", err)
			return
		}

		// 使用 goroutine 处理每个客户端连接
		go handleTCPClientConnection(conn, port)
	}
}

func handleTCPClientConnection(conn *net.TCPConn, port int) {
	defer conn.Close()

	b := make([]byte, 2048)
	for {
		cc, err := conn.Read(b)
		if err != nil {
			if err.Error() == "EOF" {
				// 客户端断开连接
				fmt.Printf("Connection closed on port %d\n", port)
				return
			}
			fmt.Printf("TCP read error: %s\n", err)
			return
		}

		// 处理收到的消息
		reqMsg := string(b[:cc])
		fmt.Printf("[REQ-TCP-PORT: %v]-[MSG: %q]-[REMOTE-ADDR: %v]\n", port, reqMsg, conn.RemoteAddr())

		// 根据不同端口返回不同的消息
		resMsg := "PONG"

		_, wrerr := conn.Write([]byte(resMsg))
		if wrerr != nil {
			fmt.Printf("TCP write error: %s\n", wrerr)
		} else {
			fmt.Printf("%s\n", resMsg)
			fmt.Println("---------------------------------")
		}
	}
}
