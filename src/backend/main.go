package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/cors"
)

// ADT Tree untuk tiap artikel
type TreeNode struct {
	Parent   *TreeNode
	Root     string
	Children []*TreeNode
}

// Fungsi untuk menentukan apakah sebuah artikel dapat diakses dari root
func (node *TreeNode) isChild(str string) bool {
	for _, child := range node.Children {
		if child.Root == str {
			return true
		}
	}
	return false
}

// Handle Redirect
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

// Fungsi untuk mendapatkan artikel yang dilalui
func (node *TreeNode) GetPath(path []string) []string {
	if node.Parent != nil {
		path = node.Parent.GetPath(path)
		path = append(path, GetTitle(node.Root))
		return path
	} else {
		path = append(path, GetTitle(node.Root))
		return path
	}
}

var visited = make(map[string]int)

func ClearVisited() {
	for key := range visited {
		delete(visited, key)
	}
}

// WikipediaPage represents a Wikipedia page
type WikipediaPage struct {
	Query struct {
		Pages map[string]struct {
			Missing bool `json:"missing"`
		} `json:"pages"`
	} `json:"query"`
}

// Function to check if a URL corresponds to a valid article on English Wikipedia
func IsTitleValid(url string) bool {
	// Extract page title from URL
	parts := strings.Split(url, "/")
	pageTitle := parts[len(parts)-1]

	// Make a request to the Wikipedia API
	resp, err := http.Get(fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&format=json&prop=info&titles=%s&redirects=1", pageTitle))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var wikiPage WikipediaPage
	if err := json.NewDecoder(resp.Body).Decode(&wikiPage); err != nil {
		return false
	}

	// Check if the page exists
	if len(wikiPage.Query.Pages) == 0 || wikiPage.Query.Pages[pageTitle].Missing {
		return false
	}
	return true
}

type RequestServer struct {
	Start       string `json:"start"`
	Destination string `json:"destination"`
	Algo        string `json:"algo"`
}

type ResponseServer struct {
	Success string `json:"success"`
	Output  string `json:"output"`
}

func main() {
	port := 8080
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	// Create a new ServeMux (router)
	mux := http.NewServeMux()

	// register a handler func for the route
	mux.HandleFunc("/req", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Decode the request body into a NameRequest struct
		var request RequestServer
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Generate the response message
		var response ResponseServer
		if IsTitleValid(request.Start) && IsTitleValid((request.Destination)) {
			// Call the BFS/IDS Function
			// make a temp variable
			var b int
			var c []string
			var a, d int64
			if request.Algo == "BFS" {
				a, b, c, d = BFS(request.Start, request.Destination)
			} else {
				a, b, c, d = IDS(request.Start, request.Destination)
			}
			response.Output = fmt.Sprintf("Jumlah Artikel yang diperiksa : %d<br/>Jumlah Artikel yang dilalui : %d<br/>Rute penjelajahan : %s", a, b, c[0])
			for i := 1; i < len(c); i++ {
				response.Output += fmt.Sprintf("->%s", c[i])
			}
			response.Output += fmt.Sprintf("<br />Waktu pencarian : %d ms", d)
			response.Success = "Success"
		} else {
			response.Success = "Fail"
		}
		// Marshal the response
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
		// Set the Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Write the JSON response with status code 200 (OK)
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	})
	// Start the HTTP server on the specified port
	fmt.Printf("Server listening on port %d\n", port)

	http.ListenAndServe(":8080", c.Handler(mux))
}
