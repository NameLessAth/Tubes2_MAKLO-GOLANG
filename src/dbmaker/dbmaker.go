package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func isIn(list []string, target string) bool {
	for _, val := range list {
		if val == target {
			return true
		}
	}
	return false
}

var data = make(map[string][]string)

func addToData(title string) {
	var child []string
	var link = "https://en.wikipedia.org/wiki/"
	link += title
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
		if strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":") && !isIn(child, link[6:]) {
			if strings.Contains(link, "#") {
				hashtagIndex := strings.Index(link, "#")
				link = link[:hashtagIndex]
			}
			child = append(child, link[6:])
		}
	})

	data[title] = child
}

func main() {
	file, err := os.Open("all") // Change to your file path
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a Scanner to read the file
	scanner := bufio.NewScanner(file)
	// Remove first line from all file
	scanner.Scan()

	var n int = 2
	for scanner.Scan() {
		fmt.Println("Processing", scanner.Text())
		addToData(scanner.Text())
		n--
		if n == 0 {
			break
		}
	}

	// Marshal the map to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error marshalling to JSON:", err)
		return
	}

	// Create a file
	file2, err := os.Create("output.json")
	if err != nil {
		fmt.Println("error creating file:", err)
		return
	}
	defer file2.Close()

	// Write JSON data to the file
	_, err = file2.Write(jsonData)
	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}

	fmt.Println("JSON data written to file successfully")
}
