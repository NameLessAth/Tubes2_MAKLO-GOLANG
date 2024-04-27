package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

// Fungsi untuk mengecek apakah link yang diinputkan merupakan link yang valid
func IsTitleValid(linkName string) bool {
	// Request the HTML page.
	var link string
	link = "https://en.wikipedia.org/wiki/"
	link += linkName
	res, err := http.Get(link)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return true
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

// func queryHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

//		var input InputData
//		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
//			w.WriteHeader(http.StatusBadRequest)
//			fmt.Fprintf(w, "Error: %v", err)
//			return
//		}
//		fmt.Printf("Received input: %+v", input)
//		responseData := map[string]string{"message": "Received input successfully"}
//		json.NewEncoder(w).Encode(responseData)
//	}
type RequestServer struct {
	Start       string `json:"start"`
	Destination string `json:"destination"`
}

type ResponseServer struct {
	Output string `json:"output"`
}

func main() {
	var a, d int64
	var b int
	var c []string
	a, b, c, d = BFS("Garut", "Islam")
	fmt.Println(a, b, c, d)
}

// func main() {
// 	port := 8080
// 	c := cors.New(cors.Options{
// 		AllowedOrigins: []string{"*"},
// 	})
// 	// Create a new ServeMux (router)
// 	mux := http.NewServeMux()

// 	// register a handler func for the route
// 	mux.HandleFunc("/req", func(w http.ResponseWriter, r *http.Request) {
// 		// Check if the request method is POST
// 		if r.Method != http.MethodPost {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}
// 		// Decode the request body into a NameRequest struct
// 		var request RequestServer
// 		err := json.NewDecoder(r.Body).Decode(&request)
// 		if err != nil {
// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
// 			return
// 		}
// 		// Generate the response message
// 		response := ResponseServer{
// 			Output: fmt.Sprintf("Start: %s %s", request.Start, request.Destination),
// 		}
// 		// Encode response data to JSON
// 		responseJSON, err := json.Marshal(response)
// 		if err != nil {
// 			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
// 			return
// 		}
// 		// Set the Content-Type header to application/json
// 		w.Header().Set("Content-Type", "application/json")
// 		// CORS
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		// Write the JSON response with status code 200 (OK)
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(responseJSON)
// 	})
// 	// Start the HTTP server on the specified port
// 	fmt.Printf("Server listening on port %d\n", port)

// 	http.ListenAndServe(":8080", c.Handler(mux))
// }
