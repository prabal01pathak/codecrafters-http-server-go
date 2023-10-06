package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func getHeader(data string) (string, error) {
	fmt.Printf("data is: %v", data)
	regex := regexp.MustCompile(`User-Agent:\s*(.*)`)
	headMatch := regex.FindStringSubmatch(data)
	if len(headMatch) > 0 {
		s := headMatch[len(headMatch)-1]
		s = strings.TrimSpace(s)
		return s, nil
	}
	return "", errors.New("no value matches")
}

func handleConnection(conn net.Conn) {
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("failed:", err)
		return
	}
	data := string(buff)
	header, _ := getHeader(data)
	fmt.Println(header)
	x := strings.Split(data, " ")[1]
	endpointRegex := regexp.MustCompile(`/echo/([^/]+.*)$`)
	match := endpointRegex.FindStringSubmatch(x)
	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v"
	switch {
	case x == "/":
		s := fmt.Sprintf(
			response,
			0, "",
		)
		conn.Write([]byte(s))
	case x == "/user-agent":
		s := fmt.Sprintf(response, len(header), header)
		conn.Write([]byte(s))
	case strings.HasPrefix(x, "/echo"):
		s := fmt.Sprintf(
			response,
			len(match[len(match)-1]), match[len(match)-1],
		)
		conn.Write([]byte(s))
	default:
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
	count := 1
	for {
		conn, err := l.Accept()
		count++
		go handleConnection(conn)
		if err != nil {
			fmt.Println("--- Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Printf("count is: %v", count)
	}
	// }
}
