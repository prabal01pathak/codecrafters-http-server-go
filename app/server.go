package main

import (
	"fmt"
	// "log"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	buff := make([]byte, 1)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("failed:", err)
		return
	}
	fmt.Println("Got: ", string(buff))
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\rn"))
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
	// connections = append(connections, conn)
	go handleConnection(conn)
	// fmt.Printf("didn't get uo: %v %v", err, connections)
	// conn.Write("hello")
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		// os.Exit(1)
	}
	// }
}
