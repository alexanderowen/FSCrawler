package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sync"
	"time"
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
	q.numWaiting++
	for q.Len() == 0 {
		if q.numWaiting == q.numRoutines { // barrier condition, no more work
			q.cond.Broadcast()
			return nil
		}
		q.cond.Wait()
	}
	n := (q.list)[0]
	q.list = (q.list)[1:]

	q.cond.Broadcast()
	return n
}

func (q *Queue) Len() int {
	return len(q.list)
}

func crawler(workQ *Queue, reg *regexp.Regexp) {
	for {
		n := workQ.Pop()
		if n == nil {
			return //all done
		}

		fmt.Printf("Testing: '%s'\n", n.val)
		files, err := ioutil.ReadDir(n.val)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Found an error: '%s'\n", err)
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				workQ.Push(&Node{val: n.val + "/" + file.Name()})
			}
			if reg.MatchString(file.Name()) {
				fmt.Printf("Found a match: %s\n", file.Name())
			}
		}
	}
}

func convertToRegexp(pat string) string {
	var reg string
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
	return reg
}

func main() {
	pattern := os.Args[1]
	var validID = regexp.MustCompile(convertToRegexp(pattern))

	work := NewQueue(1)
	n := &Node{
		val: ".",
	}
	work.Push(n)
	go crawler(work, validID)
	/*
		work.cond.L.Lock()
		for work.numWaiting != work.numRoutines {
			work.cond.Wait()
		}
		work.cond.L.Unlock()
	*/
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main finished.")

}
