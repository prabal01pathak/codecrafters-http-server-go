package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func handleConnection(conn net.Conn) {
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("failed:", err)
		return
	}
	data := string(buff)
	x := strings.Split(data, " ")[1]
	fmt.Printf("x is: %v", len(x))
	regex := regexp.MustCompile(`/echo/([^/]+.*)$`)
	match := regex.FindStringSubmatch(x)
	fmt.Printf("matches are: %v", match)

	if x == "/" {
		response := fmt.Sprintf(
			"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v",
			0, "",
		)
		conn.Write([]byte(response))
	} else if len(match) > 1 {
		response := fmt.Sprintf(
			"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v",
			len(match[len(match)-1]), match[len(match)-1],
		)
		conn.Write([]byte(response))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	conn.Close()
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	// for {
	conn, err := l.Accept()
	handleConnection(conn)
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	// }
}
