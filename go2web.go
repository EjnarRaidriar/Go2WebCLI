package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
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
			searchUrl(args[1])
		}
	}
	if args[0] == SEARCH {
		if len(args[1:]) < 1 {
			fmt.Println("Error:\n\t-s expects at least one argument")
		} else {
			fmt.Printf("Searching term: %s\n", strings.Join(args[1:], " "))
			searchArg(args[1:])
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

func getResponse(urlStr string) (string, error) {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL: ", err)
		return "", err
	}

	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "https"
		// urlStr = parsedUrl.String()
		fmt.Println("No scheme provided, using https by default")
	}

	host := parsedUrl.Host
	path := parsedUrl.Path

	if path == "" {
		path = "/"
	}
	if parsedUrl.RawQuery != "" {
		path += "?" + parsedUrl.RawQuery
	}

	fmt.Printf("Host %s, Path: %s\n", host, path)

	port := "443"
	if parsedUrl.Scheme == "http" {
		port = "80"
	}
	if strings.Contains(host, ":") {
		hostParts := strings.Split(host, ":")
		host = hostParts[0]
		port = hostParts[1]
	}

	var conn net.Conn
	var connErr error

	if parsedUrl.Scheme == "https" {
		config := &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: false,
		}
		dialer := &net.Dialer{
			Timeout: 30 * time.Second,
		}
		conn, connErr = tls.DialWithDialer(dialer, "tcp", host+":"+port, config)
	} else {
		dialer := &net.Dialer{
			Timeout: 30 * time.Second,
		}
		conn, connErr = dialer.Dial("tcp", host+":"+port)
	}

	if connErr != nil {
		fmt.Println("Error connecting: ", connErr)
		return "", err
	}
	defer conn.Close()

	httpReq := "GET " + path + " HTTP/1.1\r\n"
	httpReq += "Host: " + host + "\r\n"
	httpReq += "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36\r\n"
	httpReq += "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\n"
	httpReq += "Accept-Language: en-US,en;q=0.5\r\n"
	httpReq += "Accept-Encoding: identity\r\n"
	httpReq += "Connection: close\r\n"
	httpReq += "Upgrade-Insecure-Requests: 1\r\n"
	httpReq += "Cache-Control: max-age=0\r\n"
	httpReq += "Sec-Fetch-Dest: document\r\n"
	httpReq += "Sec-Fetch-Mode: navigate\r\n"
	httpReq += "Sec-Fetch-Site: none\r\n"
	httpReq += "Sec-Fetch-User: ?1\r\n"
	httpReq += "Pragma: no-cache\r\n"
	httpReq += "DNT: 1\r\n"
	httpReq += "\r\n"

	fmt.Println("Sending request to: " + host + path)
	_, err = conn.Write([]byte(httpReq))
	if err != nil {
		fmt.Println("Error sendig request: ", err)
		return "", err
	}

	reader := bufio.NewReader(conn)

	statusLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading status: ", err)
		return "", err
	}

	fmt.Println("Status: ", statusLine)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading response body: ", err)
		return "", err
	}

	err = os.WriteFile("response.html", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error saving response:", err)
		return "", err
	}

	return buf.String(), nil
}

func searchUrl(urlStr string) {
	response, err := getResponse(urlStr)
	if err != nil {
		return
	}

	fmt.Println("\nResponse body:")
	fmt.Println(response)
}

func searchArg(arguments []string) {
	host := "www.bing.com"
	path := "/search?q=" + url.QueryEscape(strings.Join(arguments, " "))
	urlStr := "https://" + host + path
	response, err := getResponse(urlStr)
	if err != nil {
		return
	}

	printSearchResults(response)
}

func printSearchResults(body string) {
	regEx := regexp.MustCompile(`<h2><a href="([^"]+)"`)
	matches := regEx.FindAllStringSubmatch(body, -1)
	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, match[1])
		}
	}

	if len(results) > 0 {
		fmt.Printf("\nFound %d search results:\n", len(results))
		maxResults := 10
		if len(results) < maxResults {
			maxResults = len(results)
		}
		for i := 0; i < maxResults; i++ {
			fmt.Printf("%d. %s\n", i+1, results[i])
		}
	} else {
		fmt.Println(("No search results found in the response"))
	}
}
