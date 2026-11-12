package solution

// ListNode represents a node in a singly linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ReverseList reverses a singly linked list in-place.
func ReverseList(head *ListNode) *ListNode {
	var prev *ListNode
	curr := head
	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}
	return prev
}
