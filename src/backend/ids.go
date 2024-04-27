package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"	
	"time"

	"github.com/PuerkitoBio/goquery"
)

var found bool = false

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
			var title string = GetTitle(link[6:])
			if !node.isChild(link[6:]) {
				child := TreeNode{Parent: node, Root: title}
				node.Children = append(node.Children, &child)
			}
		}
	})
}

func DLS(start *TreeNode, depth int, goal string) {

	if start.Root == goal { //menemukan goal
		var path []string
		path = start.GetPath(path)

		var i int = 1
		fmt.Println("Path :")
		for _, page := range path {
			fmt.Printf("%d. %s\n", i, GetTitle(page))
			i++
		}
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
				DLS(v, depth-1, goal)
			}
			if found {
				return
			}
		}
	}
}

func IDS(goal string) {
	var depth int = 1
	start := time.Now()

	for !found {
		go DLS(queue[0], depth, goal)
		depth++
	}
	end := time.Since(start)
	fmt.Printf("Time elapsed: %s", end)
}
