package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TreeNode struct {
	Parent   *TreeNode
	Root     string
	Children []*TreeNode
}

func (node *TreeNode) isChild(str string) bool {
	for _, child := range node.Children {
		if child.Root == str {
			return true
		}
	}
	return false
}

var queue []*TreeNode

func Enqueue(node *TreeNode) {
	queue = append(queue, node)
}

func Dequeue() {
	queue = queue[1:]
}

func (node *TreeNode) AddChildToQueue() {
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
		if strings.HasPrefix(link, "/wiki/") && (!strings.HasPrefix(link, "/wiki/Category:") && !strings.HasPrefix(link, "/wiki/Template:") && !strings.HasPrefix(link, "/wiki/Special:") && !strings.HasSuffix(link, ".jpg") && !strings.HasSuffix(link, ".svg") && !strings.HasSuffix(link, ".png")) {
			if !node.isChild(link[6:]) {
				child := TreeNode{Parent: node, Root: link[6:]}
				node.Children = append(node.Children, &child)
				Enqueue(&child)
			}
		}
	})
}

func (node *TreeNode) GetPath(path []string) []string {
	if node.Parent != nil {
		path = node.Parent.GetPath(path)
		path = append(path, node.Root)
		return path
	} else {
		path = append(path, node.Root)
		return path
	}
}

var found bool = false
var visited = make(map[string]int)

func BFS(goal string) {
	for !found {
		if queue[0].Root == goal {
			var path []string
			path = queue[0].GetPath(path)

			var i int = 1
			fmt.Println("Path :")
			for _, page := range path {
				fmt.Printf("%d. %s\n", i, page)
				i++
			}
			found = true
		}
		if visited[queue[0].Root] == 0 {
			visited[queue[0].Root] = 1
			queue[0].AddChildToQueue()
		}
		Dequeue()
	}
}

func main() {
	var awal, akhir string
	fmt.Println("Masukkan page awal : ")
	fmt.Scanln(&awal)
	fmt.Println("Masukkan page akhir : ")
	fmt.Scanln(&akhir)
	if awal == akhir {
		fmt.Printf("Kedua page sama\n")
	} else {
		fmt.Printf("BFS/DFS\n")
		awal := TreeNode{Root: awal}
		Enqueue(&awal)
		BFS(akhir)
	}
}
