package tools

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func PortForward() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <target_port> <127.0.0.1:local_port>")
		return
	}

	// 本地端口
	listenPort := ":" + os.Args[1]
	// 目标地址 (例如：本地 8080)
	targetAddr := os.Args[2]

	// 监听本地端口
	listener, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", listenPort, err)
	}
	defer listener.Close()

	log.Printf("Listening on %s, forwarding to %s", listenPort, targetAddr)

	for {
		// 接受外部连接
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// 连接到本地的 8080 端口
		targetConn, err := net.Dial("tcp", targetAddr)
		if err != nil {
			log.Printf("Failed to connect to target address %s: %v", targetAddr, err)
			clientConn.Close()
			continue
		}

		// 启动 goroutine 来转发数据
		go forward(clientConn, targetConn)
	}
}

// 转发数据
func forward(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()

	// 从 src 读取数据并写入 dst
	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Printf("Error copying from src to dst: %v", err)
		}
	}()

	// 从 dst 读取数据并写入 src
	if _, err := io.Copy(src, dst); err != nil {
		log.Printf("Error copying from dst to src: %v", err)
	}
}
