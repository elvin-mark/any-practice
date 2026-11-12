package reverselist

import (
	"reflect"
	"testing"
)

func sliceToList(nums []int) *ListNode {
	dummy := &ListNode{}
	curr := dummy
	for _, n := range nums {
		curr.Next = &ListNode{Val: n}
		curr = curr.Next
	}
	return dummy.Next
}

func listToSlice(head *ListNode) []int {
	var result []int
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}

func TestReverseList(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "multiple elements",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{5, 4, 3, 2, 1},
		},
		{
			name:     "two elements",
			input:    []int{1, 2},
			expected: []int{2, 1},
		},
		{
			name:     "single element",
			input:    []int{42},
			expected: []int{42},
		},
		{
			name:     "empty list",
			input:    []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := sliceToList(tt.input)
			reversed := ReverseList(head)
			got := listToSlice(reversed)
			if len(got) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ReverseList(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func BenchmarkReverseList(b *testing.B) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	head := sliceToList(nums)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		head = ReverseList(head)
	}
}
