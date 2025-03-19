package main

import (
	"fmt"
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
