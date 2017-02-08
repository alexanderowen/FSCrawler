package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"sync"
)

type Node struct {
	val string
}

type Queue struct {
	list        []*Node
	cond        *sync.Cond
	numRoutines int
	numWaiting  int
}

func NewQueue(num int) *Queue {
	return &Queue{
		list:        make([]*Node, 0),
		cond:        &sync.Cond{L: &sync.Mutex{}},
		numRoutines: num,
		numWaiting:  0,
	}
}

func (q *Queue) Push(n *Node) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.list = append(q.list, n)
}

func (q *Queue) Pop() *Node {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	defer q.cond.Broadcast()
	q.numWaiting++
	for q.Len() == 0 {
		//fmt.Printf("Waiting for pop. Waiting: %d and Total: %d\n", q.numWaiting, q.numRoutines)
		if q.numWaiting == q.numRoutines { // barrier condition, no more work
			return nil
		}
		q.cond.Wait()
	}
	q.numWaiting--
	n := (q.list)[0]
	q.list = (q.list)[1:]
	return n
}

func (q *Queue) Len() int {
	return len(q.list)
}

func crawler(workQ *Queue, reg *regexp.Regexp, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		n := workQ.Pop()
		if n == nil {
			return //all done
		}

		//fmt.Printf("Testing: '%s'\n", n.val)
		files, err := ioutil.ReadDir(n.val)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Found an error: '%s'\n", err)
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
	var wg sync.WaitGroup

	var numRoutines int
	if env := os.Getenv("CRAWLER_THREADS"); env != "" {
		numRoutines, _ = strconv.Atoi(env)
	} else {
		numRoutines = 2
	}

	work := NewQueue(numRoutines)
	n := &Node{
		val: ".",
	}
	work.Push(n)
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go crawler(work, reg, &wg)
	}

	wg.Wait()
}
