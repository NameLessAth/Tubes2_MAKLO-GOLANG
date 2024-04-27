package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var found bool = false

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
	} else if depth <= 0 { //sudah mencapai kedalaman maksimum tetapi tidak menemukan goal
		found = boolHelper(false)
	} else { //lanjutkan pencarian
		if visited[start.Root] == 0 { //menandai current node sudah dikunjungi
			visited[start.Root] = 1
		}
		start.AddChildren()
		for _, v := range start.Children { //melanjutkan DFS pada children dari node start
			if visited[v.Root] == 0 {
				DLS(wg, v, append(curPath, v), depth-1, goal, found)
			}
			if found == boolHelper(true) {
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

	for !found {
		wg.Add(1)
		go DLS(&wg, queue[0], curPath, depth, goal, &found)
		depth++
	}
	end := time.Since(start)
	fmt.Printf("Time elapsed: %s", end)
}
