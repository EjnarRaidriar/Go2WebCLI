package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
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

func toGoogleSearch(search []string) string {
	var searchBuilder strings.Builder
	for _, term := range search {
		if strings.Contains(term, "+") {
			var sb strings.Builder
			for _, char := range term {
				if char == '+' {
					sb.WriteString("%2B")
				} else {
					sb.WriteRune(char)
				}
			}
			searchBuilder.WriteString(sb.String())
			searchBuilder.WriteRune('+')
		} else {
			searchBuilder.WriteString(term)
			searchBuilder.WriteRune('+')
		}
	}
	result := searchBuilder.String()
	return result[:len(result)-1]
}

func search(arguments []string) {
	host := "www.google.com"
	path := "/search?q=" + toGoogleSearch(arguments)

	addr, err := net.ResolveTCPAddr("tcp", host+":80")
	if err != nil {
		log.Fatal("could not resolve tcp address:", err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal("could not dial:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Sending request to: " + host + path)
	httpReq := "GET " + path + " HTTP/1.0\r\n"
	httpReq += "\r\n"
	fmt.Fprintf(conn, "%s", httpReq)

	var buf bytes.Buffer
	io.Copy(&buf, conn)
	fmt.Printf("got this response:\n%s", buf.Bytes())
}
