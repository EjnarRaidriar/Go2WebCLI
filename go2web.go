package main

import (
	"bufio"
	"bytes"
	"fmt"
	htmlparser "go2web/htmlParser"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
)

const (
	SEARCH string = "-s"
	VIEW   string = "-u"
)

func main() {
	args := os.Args[1:]
	if args[0] == VIEW {
		if len(args[1:]) != 1 {
			fmt.Println("Error:\n\t-u expects single url argument")
		} else {
			fmt.Printf("Making request to %s\n", args[1])
		}
	}
	if args[0] == SEARCH {
		if len(args[1:]) < 1 {
			fmt.Println("Error:\n\t-s expects at least one argument")
		} else {
			fmt.Printf("Searching term: %s\n", strings.Join(args[1:], " "))
			search(args[1:])
		}
	}
	if args[0] == "-h" {
		fmt.Println("Make an HTTP request to the specified URL")
		fmt.Println("\t-u <URL>")
		fmt.Println("Make an HTTP request to search the term")
		fmt.Println("\t-s <search-term>")
		fmt.Println("Show this help")
		fmt.Println("\t-h")
	}
}

func search(arguments []string) {
	host := "www.google.com"
	path := "/search?q=" + url.QueryEscape(strings.Join(arguments, " "))
	// path = "/"

	conn, err := net.Dial("tcp", host+":80")
	if err != nil {
		fmt.Println("Error connecting: ", err)
		return
	}
	defer conn.Close()

	httpReq := "GET " + path + " HTTP/1.1\r\n"
	httpReq += "Host: " + host + "\r\n"
	httpReq += "User-Agent: Go2WebSearch/1.0\r\n"
	httpReq += "Accept: text/html\r\n"
	httpReq += "Connection: close\r\n"
	httpReq += "\r\n"

	fmt.Println("Sending request to: " + host + path)
	_, err = conn.Write([]byte(httpReq))
	if err != nil {
		fmt.Println("Error sendig request: ", err)
	}

	reader := bufio.NewReader(conn)
	fmt.Println("Response headers:")
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		fmt.Print(line)
	}

	fmt.Println("\nResponse body:")
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading response body: ", err)
		return
	}
	fmt.Printf("Body: %d bytes\n", buf.Len())
	// fmt.Println(buf.String())
	err = os.WriteFile("google_response.html", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error saving response:", err)
		return
	}
	htmlparser.Test(strings.NewReader(buf.String()))
}
