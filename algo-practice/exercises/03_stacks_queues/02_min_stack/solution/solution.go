package solution

// MinStack supports push, pop, top, and retrieving the minimum element in O(1) time.
type MinStack struct {
	stack    []int
	minStack []int
}

func Constructor() MinStack {
	return MinStack{
		stack:    make([]int, 0),
		minStack: make([]int, 0),
	}
}

func (s *MinStack) Push(val int) {
	s.stack = append(s.stack, val)
	if len(s.minStack) == 0 || val <= s.minStack[len(s.minStack)-1] {
		s.minStack = append(s.minStack, val)
	} else {
		s.minStack = append(s.minStack, s.minStack[len(s.minStack)-1])
	}
}

func (s *MinStack) Pop() {
	if len(s.stack) > 0 {
		s.stack = s.stack[:len(s.stack)-1]
		s.minStack = s.minStack[:len(s.minStack)-1]
	}
}

func (s *MinStack) Top() int {
	if len(s.stack) == 0 {
		return 0
	}
	return s.stack[len(s.stack)-1]
}

func (s *MinStack) GetMin() int {
	if len(s.minStack) == 0 {
		return 0
	}
	return s.minStack[len(s.minStack)-1]
}
