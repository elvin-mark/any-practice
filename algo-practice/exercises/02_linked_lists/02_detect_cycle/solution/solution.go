package solution

// ListNode represents a node in a singly linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// HasCycle returns true if the linked list has a cycle using Floyd's algorithm.
func HasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	slow := head
	fast := head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			return true
		}
	}

	return false
}
