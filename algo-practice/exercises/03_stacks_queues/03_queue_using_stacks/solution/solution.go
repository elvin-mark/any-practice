package solution

// MyQueue implements a FIFO queue using two stacks.
type MyQueue struct {
	inStack  []int
	outStack []int
}

func Constructor() MyQueue {
	return MyQueue{
		inStack:  make([]int, 0),
		outStack: make([]int, 0),
	}
}

func (q *MyQueue) Push(x int) {
	q.inStack = append(q.inStack, x)
}

func (q *MyQueue) Pop() int {
	q.shift()
	val := q.outStack[len(q.outStack)-1]
	q.outStack = q.outStack[:len(q.outStack)-1]
	return val
}

func (q *MyQueue) Peek() int {
	q.shift()
	return q.outStack[len(q.outStack)-1]
}

func (q *MyQueue) Empty() bool {
	return len(q.inStack) == 0 && len(q.outStack) == 0
}

func (q *MyQueue) shift() {
	if len(q.outStack) == 0 {
		for len(q.inStack) > 0 {
			n := len(q.inStack) - 1
			q.outStack = append(q.outStack, q.inStack[n])
			q.inStack = q.inStack[:n]
		}
	}
}
