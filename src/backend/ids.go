package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var found bool
var artikelDiperiksa int64
var path []string
var depth int

func (node *TreeNode) AddChildren() {
	// Request the HTML page.
	var link string
	link = "https://en.wikipedia.org/wiki/"
	link += node.Root
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("%s\n", link)
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the links in the main content section
	doc.Find("#mw-content-text").Find("a").Each(func(i int, s *goquery.Selection) {
		// For each link found, get the href attribute and text
		link, _ := s.Attr("href")
		// Check if the link stays within Wikipedia domain
		if strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":") {
			var title string = link[6:]
			if !node.isChild(link[6:]) {
				child := TreeNode{Parent: node, Root: title}
				node.Children = append(node.Children, &child)
			}
		}
	})
}

func DLS(start *TreeNode, depth int, goal string) {

	if start.Root == goal { //menemukan goal
		path = start.GetPath(path)
		found = true
	} else if depth <= 0 { //sudah mencapai kedalaman maksimum tetapi tidak menemukan goal
		found = false
	} else { //lanjutkan pencarian
		if visited[start.Root] == 0 { //menandai current node sudah dikunjungi
			visited[start.Root] = 1
		}
		start.AddChildren()
		for _, v := range start.Children { //melanjutkan DFS pada children dari node start
			if visited[v.Root] == 0 {
				artikelDiperiksa++
				DLS(v, depth-1, goal)
			}
			if found {
				return
			}
		}
	}
}

func IDS(begin string, goal string) (int64, int, []string, int64) {
	ClearVisited()
	depth = 1
	artikelDiperiksa = 0
	path = []string{}
	found = false

	start := time.Now()

	var startNode TreeNode
	startNode.Root = begin
	startNode.AddChildren()
	for !found {
		DLS(&startNode, depth, goal)
		ClearVisited()
		depth++
	}

	end := time.Since(start).Milliseconds()
	// Return
	return artikelDiperiksa, len(path), path, end
}
