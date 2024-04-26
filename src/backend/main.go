package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

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
		if strings.HasPrefix(link, "/wiki/") && (!strings.HasPrefix(link, "/wiki/Category:") && !strings.HasPrefix(link, "/wiki/Template:") && !strings.HasPrefix(link, "/wiki/Special:") && !strings.HasSuffix(link, ".jpg") && !strings.HasSuffix(link, ".svg") && !strings.HasSuffix(link, ".png")) {
			if !node.isChild(link[6:]) {
				child := TreeNode{Parent: node, Root: link[6:]}
				node.Children = append(node.Children, &child)
			}
		}
	})
}

func GetTitle(linkName string) string {
	// Request the HTML page.
	var link string
	link = "https://en.wikipedia.org/wiki/"
	link += linkName
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
	var temp string = doc.Find("title").Text()
	idEnd := strings.LastIndex(temp, " - Wikipedia")
	return temp[0:idEnd]
	
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
	start := time.Now()
	for !found {
		if queue[0].Root == goal {
			var path []string
			path = queue[0].GetPath(path)

			var i int = 1
			fmt.Println("Path :")
			for _, page := range path {
				fmt.Printf("%d. %s\n", i, GetTitle(page))
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
	end := time.Since(start)
	fmt.Printf("Time elapsed: %s", end);
}

func boolHelper(b bool) *bool {
	return &b
}

func DLS(wg *sync.WaitGroup, start *TreeNode, curPath []*TreeNode, depth int, goal string, found *bool) {
	defer wg.Done()

	if start.Root == goal { //menemukan goal
		var path []string
		path = start.GetPath(path)

		var i int = 1
		fmt.Println("Path :")
		for _, page := range path {
			fmt.Printf("%d. %s\n", i, GetTitle(page))
			i++
		}
		found = boolHelper(true)
	} else if (depth <= 0){ //sudah mencapai kedalaman maksimum tetapi tidak menemukan goal
		found = boolHelper(false)
	} else { //lanjutkan pencarian
		if visited[start.Root] == 0 { //menandai current node sudah dikunjungi
			visited[start.Root] = 1
		}
		start.AddChildren()
		for _, v := range start.Children { //melanjutkan DFS pada children dari node start
			if (visited[v.Root]==0) {DLS(wg, v, append(curPath, v), depth-1, goal, found)}
			if (found == boolHelper(true)) {
				break
			}
		}
	}
}

func IDS(goal string) {
	runtime.GOMAXPROCS(5)
	var wg sync.WaitGroup
	var curPath []*TreeNode
	curPath = append(curPath, queue[0])
	var depth int = 1
	start := time.Now()

	for (!found) {
		wg.Add(1)
		go DLS(&wg, queue[0], curPath, depth, goal, &found)
		depth++
	}
	end := time.Since(start)
	fmt.Printf("Time elapsed: %s", end);
}

func main() {
	var awal, akhir, metode string
	fmt.Println("Masukkan page awal : ")
	fmt.Scanln(&awal)
	fmt.Println("Masukkan page akhir : ")
	fmt.Scanln(&akhir)
	if awal == akhir {
		fmt.Printf("Kedua page sama\n")
	} else {
		fmt.Printf("BFS/IDS\n")
		fmt.Scanln(&metode)
		awal := TreeNode{Root: awal}
		Enqueue(&awal)
		fmt.Printf("Pencarian dimulai!\n")
		if (metode == "BFS") {
			BFS(akhir)
		} else {
			IDS(akhir)
		}
	}
}
