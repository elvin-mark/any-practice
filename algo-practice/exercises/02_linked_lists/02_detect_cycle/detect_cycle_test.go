package detectcycle

import "testing"

func createListWithCycle(nums []int, pos int) *ListNode {
	if len(nums) == 0 {
		return nil
	}
	head := &ListNode{Val: nums[0]}
	curr := head
	var cycleEntry *ListNode
	if pos == 0 {
		cycleEntry = head
	}

	for i := 1; i < len(nums); i++ {
		curr.Next = &ListNode{Val: nums[i]}
		curr = curr.Next
		if i == pos {
			cycleEntry = curr
		}
	}

	if pos >= 0 {
		curr.Next = cycleEntry
	}

	return head
}

func TestHasCycle(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		pos      int // -1 means no cycle
		expected bool
	}{
		{
			name:     "cycle at index 1",
			nums:     []int{3, 2, 0, -4},
			pos:      1,
			expected: true,
		},
		{
			name:     "cycle at index 0 (self/loop)",
			nums:     []int{1, 2},
			pos:      0,
			expected: true,
		},
		{
			name:     "no cycle single element",
			nums:     []int{1},
			pos:      -1,
			expected: false,
		},
		{
			name:     "no cycle multiple elements",
			nums:     []int{1, 2, 3, 4, 5},
			pos:      -1,
			expected: false,
		},
		{
			name:     "empty list",
			nums:     []int{},
			pos:      -1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := createListWithCycle(tt.nums, tt.pos)
			got := HasCycle(head)
			if got != tt.expected {
				t.Errorf("HasCycle() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func BenchmarkHasCycle(b *testing.B) {
	nums := make([]int, 1000)
	head := createListWithCycle(nums, 500)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = HasCycle(head)
	}
}
