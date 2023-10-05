package main

import (
	"fmt"
	"regexp"
	"strings"
	// "log"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	re := regexp.MustCompile(`(?m)\/echo\/.*`)
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	// fmt.Printf("data is: %v", d)
	if err != nil {
		fmt.Println("failed:", err)
		return
	}
	data := string(buff)
	x := strings.Split(data, " ")[1]
	fmt.Printf("x is: %v", len(x))
	s := strings.Split(x, "/")
	found := false
	if x == "/" {
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n%v", s[len(s)-1])
		conn.Write([]byte(response))
		found = true
	}
	for i, match := range re.FindAllString(x, -1) {
		fmt.Println(match, "found at index", i)
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n%v", s[len(s)-1])
		conn.Write([]byte(response))
		found = true
		break
	}
	if !found {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	// if len(x) == 0 || x[1] == "echo" {
	// 	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n%v", x[len(x)-1])
	// 	conn.Write([]byte(response))
	// } else {
	// }
	// if x[len(x)-1] == "/" {
	// 	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	// } else {
	// 	fmt.Println("in else block")
	// 	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	// }
	// fmt.Printf("dat is %v", data)
	conn.Close()
}

func main() {
	// connections := make([]net.Conn, 0)
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
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
