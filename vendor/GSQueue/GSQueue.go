package GSQueue

import (
	"sync"
)

/*************** Defining a Gouroutine-safe Queue class ************/
/**************** that doubles as a barrier construct **************/

type Node struct {
	val string
}

func NewNode(v string) *Node {
	return &Node{
		val: v,
	}
}

func (n *Node) GetValue() string {
	return n.val
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
