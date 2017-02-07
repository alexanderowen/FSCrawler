package main

import "sync"

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
		list:        make([]*Node, 10),
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

func main() {
	work := NewQueue(1)
	n := &Node{
		val: "hello",
	}
	work.Push(n)
	work.Pop()
}
