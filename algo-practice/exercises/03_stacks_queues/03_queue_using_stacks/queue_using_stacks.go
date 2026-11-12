package queueusingstacks

// MyQueue implements a FIFO queue using two stacks.
type MyQueue struct {
	// TODO: Define fields
}

func Constructor() MyQueue {
	// TODO: Initialize MyQueue
	return MyQueue{}
}

func (q *MyQueue) Push(x int) {
	// TODO: Push element x to the back of queue
}

func (q *MyQueue) Pop() int {
	// TODO: Remove and return element from front of queue
	return 0
}

func (q *MyQueue) Peek() int {
	// TODO: Get the front element
	return 0
}

func (q *MyQueue) Empty() bool {
	// TODO: Return whether the queue is empty
	return false
}
