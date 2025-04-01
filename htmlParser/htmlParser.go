package htmlparser

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func Test(connection *strings.Reader) {

	doc, err := html.Parse(connection)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
	parse(doc)
}

func parse(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "id" && strings.Contains(a.Val, "rso") {
				parseNode(n)
				return
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse(c)
	}
}

func parseNode(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "data-snf" && strings.Contains(a.Val, "x5WNvb") {
				parseSearchHeader(n)
			}
			if a.Key == "data-snf" && strings.Contains(a.Val, "nke7rc") {
				parseSearchDescription(n)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNode(c)
	}
}

func parseSearchHeader(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				fmt.Println("link: ", a.Val)
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "h3" {
		fmt.Println("title: ", n.FirstChild.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseSearchHeader(c)
	}
}

func parseSearchDescription(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "VwiC3b yXK7lf p4wth r025kc hJNv6b Hdw6tb") {
				fmt.Println("description: ", n.FirstChild.FirstChild.Data)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseSearchDescription(c)
	}
}
