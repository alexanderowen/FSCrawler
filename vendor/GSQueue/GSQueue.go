/* Author: Alexander Owen
 *
 * Package defining a Goroutine-safe Queue that acts as a Barrier
 * Implemented specifically for File System Crawler
 */

package GSQueue

import (
	"sync"
)

/*************** Defining a Gouroutine-safe Queue class ************/
/**************** that doubles as a barrier construct **************/

// Node in our Queue class
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

// Queue as a Barrier, requires knowledge of goroutines using it
// Note: No mutex field; sync.Cond will create one
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

// Returns length of queue
func (q *Queue) Len() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return len(q.list)
}
