package main

import (
	GSQ "GSQueue"
	GSH "GSStringHeap"
	"container/heap"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"sync"
)

// Worker function for goroutines
func crawler(workQ *GSQ.Queue, done *GSH.GSStringHeap, reg *regexp.Regexp, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		n := workQ.Pop()
		if n == nil { // All goroutines waiting, time to die
			return
		}

		files, err := ioutil.ReadDir(n.GetValue())
		if err != nil { // error occurred
			fmt.Fprintf(os.Stderr, "Error attempting to read dir. Message: '%s'\n", err)
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				workQ.Push(GSQ.NewNode(n.GetValue() + "/" + file.Name()))
			} else if reg.MatchString(file.Name()) {
				//fmt.Printf("%s/%s\n", n.GetValue(), file.Name())
				str := fmt.Sprintf("%s/%s", n.GetValue(), file.Name())
				heap.Push(done, str)
			}
		}
	}
}

// Converts regular expression from Bash syntax to Go.regexp syntax
func convertToRegexp(pat string) string {
	reg := "^"
	for _, char := range pat {
		switch char {
		case '*':
			reg = reg + ".*"
		case '.':
			reg = reg + "\\."
		case '?':
			reg = reg + "."
		default:
			reg = reg + string(char)
		}
	}
	return reg + "$"
}

func main() {
	pattern := os.Args[1]
	reg := regexp.MustCompile(convertToRegexp(pattern))

	var dir string
	if len(os.Args) == 3 {
		dir = os.Args[2]
	} else {
		dir = "."
	}

	var numRoutines int
	if env := os.Getenv("CRAWLER_THREADS"); env != "" { //get env var
		numRoutines, _ = strconv.Atoi(env)
	} else {
		numRoutines = 2
	}

	done := GSH.NewGSStringHeap()
	heap.Init(done) // can now use the heap interface
	work := GSQ.NewQueue(numRoutines)
	n := GSQ.NewNode(dir)
	work.Push(n)
	var wg sync.WaitGroup // Use WaitGroup so main thread knows when execution is complete
	for i := 0; i < numRoutines; i++ {
		wg.Add(1) // For each goroutine created, there is another one to wait on
		go crawler(work, done, reg, &wg)
	}

	wg.Wait() // main; don't terminate until all goroutines finished
	for done.Len() > 0 {
		fmt.Printf("%s\n", heap.Pop(done))
	}
}
