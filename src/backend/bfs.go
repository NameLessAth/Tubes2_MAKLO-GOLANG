package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var queue []*TreeNode

func Enqueue(node *TreeNode) {
	queue = append(queue, node)
}

func Dequeue() {
	queue = queue[1:]
}

func ClearQueue() {
	queue = queue[:0]
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
		if strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":") {
			if !node.isChild(link[6:]) {
				child := TreeNode{Parent: node, Root: link[6:]}
				node.Children = append(node.Children, &child)
				Enqueue(&child)
			}
		}
	})
}

func BFS(initial string, goal string) (int64, int, []string, float64) {
	// Membersihkan Map Visited dan Queue
	ClearQueue()
	ClearVisited()

	// Memasukkan artikel awal pada queue
	Enqueue(&TreeNode{Root: GetTitle(initial)})

	// Deklarasi variabel penyimpan banyak artikel yang ditelusuri
	var artikelDiperiksa int64 = 0

	// Deklarasi variabel untuk menyimpan jalan yang dilalui
	var path []string

	// Mulai menghitung
	start := time.Now()

	// Deklarasi title goal
	var goalTitle string = GetTitle(goal)

	// Proses pencarian BFS, akan dicari selama node belum ditemukan
	for {
		// Artikel Tujuan ditemukan!
		if GetTitle(queue[0].Root) == goalTitle {
			// Menyimmpan jalan yang dilalui
			path = queue[0].GetPath(path)

			// Keluar dari while loop
			break
		}

		// Mengecek apakah artikel sudah pernah dikunjungi, jika sudah pernah dikunjungi maka tidak diproses
		if visited[GetTitle(queue[0].Root)] == 0 {
			// Menambah jumlah artikel yang diperiksa
			artikelDiperiksa++
			// Mencatat artikel telah dikunjungi
			visited[GetTitle(queue[0].Root)] = 1
			// Memasukkan semua link pada artikel tersebut ke queue
			queue[0].AddChildToQueue()
		}

		// Menghapus artikel dari queue
		Dequeue()
	}

	// Menghitung total waktu
	end := time.Since(start).Seconds()

	// Return
	return artikelDiperiksa, len(path), path, end
}
