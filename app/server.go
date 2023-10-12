package main

import (
	"bufio"
	"strconv"
	// "errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

func readRequestBody(reader *bufio.Reader, contentLength int) ([]byte, error) {
	data := make([]byte, contentLength)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ParseHeaders(reader *bufio.Reader) (map[string]string, error) {
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request line:", err)
		return nil, err
	}
	requestLine = strings.TrimSpace(requestLine)
	fmt.Println("\nRequest Line:", requestLine)
	headers := make(map[string]string)
	requestMethod := regexp.MustCompile(`^[^\s]*`)
	lineParse := requestMethod.FindStringSubmatch(requestLine)
	method := lineParse[0]
	headers["Method"] = method
	headers["Endpoint"] = getEndpoint(requestLine)
	fmt.Print(headers)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		line = strings.TrimSpace(line)
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			headerName := strings.TrimSpace(headerParts[0])
			headerValue := strings.TrimSpace(headerParts[1])
			headers[headerName] = headerValue
		}
	}
	return headers, nil
}

func getEndpoint(requestLine string) string {
	value := strings.Split(requestLine, " ")[1]
	return value
}

func getFileName(requestLine string) string {
	fileNameRegx := regexp.MustCompile(`^/files/(.*)$`)
	match := fileNameRegx.FindStringSubmatch(requestLine)
	fmt.Printf("matches are: %v\n", match)
	// fileName := fmt.Sprintf("%s/%s", dir, match[len(match)-1])
	return match[len(match)-1]
}

func handleConnection(conn net.Conn, dir string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	headerValues, error := ParseHeaders(reader)
	method := headerValues["Method"]
	if error != nil {
		return
	}
	x := headerValues["Endpoint"]
	endpointRegex := regexp.MustCompile(`/echo/([^/]+.*)$`)
	match := endpointRegex.FindStringSubmatch(x)
	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v"
	notFoundRes := "HTTP/1.1 404 Not Found\r\n\r\n"
	if method == "GET" {
		switch {
		case x == "/":
			s := fmt.Sprintf(
				response,
				0, "",
			)
			conn.Write([]byte(s))
		case x == "/user-agent":
			s := fmt.Sprintf(response, len(headerValues["User-Agent"]), headerValues["User-Agent"])
			conn.Write([]byte(s))
		case strings.HasPrefix(x, "/echo"):
			s := fmt.Sprintf(
				response,
				len(match[len(match)-1]), match[len(match)-1],
			)
			conn.Write([]byte(s))
		case strings.HasPrefix(x, "/files"):
			fn := getFileName(x)
			fileName := fmt.Sprintf("%s/%s", dir, fn)
			data, err := os.ReadFile(fileName)
			if err != nil {
				conn.Write([]byte(notFoundRes))
			} else {
				fmt.Printf("data is %s and err  is %s", data, err)
				response := "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %v\r\n\r\n%v"
				s := fmt.Sprintf(response, len(data), string(data))
				conn.Write([]byte(s))
			}
		default:
			conn.Write([]byte(notFoundRes))
		}
	} else if method == "POST" {
		switch {
		case strings.HasPrefix(x, "/files"):
			fmt.Printf("files endpoint is called")
			contentLength, er := headerValues["Content-Length"]
			cl, _ := strconv.Atoi(contentLength)
			if er {
				fmt.Printf("ERRor while reading the contentLength: %v\n%v\n\nnnnoooo", er, headerValues)
			}
			file, err := readRequestBody(reader, cl)
			if err != nil {
				conn.Write([]byte(notFoundRes))
			} else {
				fn := getFileName(x)
				fileName := fmt.Sprintf("%s/%s", dir, fn)
				os.WriteFile(fileName, file, 0644)
				fmt.Printf("file name is: %s\n", fn)
				// fmt.Printf("file si: %v", file)
				response := "HTTP/1.1 201 OK"
				conn.Write([]byte(response))
			}

		default:
			conn.Write([]byte(notFoundRes))
		}
		fmt.Printf("method is: %v\n", method)
	} else {
		conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n\r\n"))
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	dir := ""
	flag.StringVar(&dir, "directory", "", "Directory")
	flag.Parse()
	fmt.Printf("direcotry is: %v\n", dir)
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	count := 1
	for {
		conn, err := l.Accept()
		go handleConnection(conn, dir)
		if err != nil {
			fmt.Println("--- Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		count++
		fmt.Printf("count is: %v", count)
	}
}
