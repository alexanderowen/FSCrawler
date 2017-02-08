package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"sync"
)

/*************** Defining a Gouroutine-safe Queue class ************/
/**************** that doubles as a barrier construct **************/

type Node struct {
	val string
}

// No mutex reference; sync.Cond will create
// one in it's initialization
type Queue struct {
	list        []*Node
	cond        *sync.Cond
	numRoutines int
	numWaiting  int
}

// Creates a new Queue
// Overlaying structure is a Go slice
func NewQueue(num int) *Queue {
	return &Queue{
		list:        make([]*Node, 0),
		cond:        &sync.Cond{L: &sync.Mutex{}},
		numRoutines: num,
		numWaiting:  0,
	}
}

// Queue.push method
func (q *Queue) Push(n *Node) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock() // defer; will execute this line just before func returns
	q.list = append(q.list, n)
}

// Queue.pop method; blocks the calling function until there is work, or when all
// goroutines are waiting, at which point it sends 'nil'
func (q *Queue) Pop() *Node {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	defer q.cond.Broadcast()
	q.numWaiting++
	for q.Len() == 0 {
		if q.numWaiting == q.numRoutines { // barrier condition met, no more work
			return nil
		}
		q.cond.Wait()
	}
	q.numWaiting--
	n := (q.list)[0]
	q.list = (q.list)[1:]
	return n
}

// Queue.Len method; returns length of queue
func (q *Queue) Len() int {
	return len(q.list)
}

/************ End Queue definition *************/

// Worker function for goroutines
func crawler(workQ *Queue, reg *regexp.Regexp, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		n := workQ.Pop()
		if n == nil { // All goroutines waiting, time to die
			return
		}

		files, err := ioutil.ReadDir(n.val)
		if err != nil { // error occurred
			fmt.Fprintf(os.Stderr, "Error attempting to read dir. Message: '%s'\n", err)
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				workQ.Push(&Node{val: n.val + "/" + file.Name()})
			} else if reg.MatchString(file.Name()) {
				fmt.Printf("%s/%s\n", n.val, file.Name())
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
	var wg sync.WaitGroup // Use WaitGroup so main thread knows when execution is complete

	var numRoutines int
	if env := os.Getenv("CRAWLER_THREADS"); env != "" { //get env var
		numRoutines, _ = strconv.Atoi(env)
	} else {
		numRoutines = 2
	}

	work := NewQueue(numRoutines)
	n := &Node{
		val: dir,
	}
	work.Push(n)
	for i := 0; i < numRoutines; i++ {
		wg.Add(1) // For each goroutine created, there is another one to wait on
		go crawler(work, reg, &wg)
	}

	wg.Wait() // main; don't terminate until all goroutines finished
}
