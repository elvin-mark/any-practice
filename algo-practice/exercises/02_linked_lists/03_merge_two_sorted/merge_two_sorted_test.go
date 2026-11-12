package mergetwosorted

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

func TestMergeTwoLists(t *testing.T) {
	tests := []struct {
		name     string
		l1       []int
		l2       []int
		expected []int
	}{
		{
			name:     "equal lengths",
			l1:       []int{1, 2, 4},
			l2:       []int{1, 3, 4},
			expected: []int{1, 1, 2, 3, 4, 4},
		},
		{
			name:     "both empty",
			l1:       []int{},
			l2:       []int{},
			expected: []int{},
		},
		{
			name:     "one empty",
			l1:       []int{},
			l2:       []int{0},
			expected: []int{0},
		},
		{
			name:     "disjoint ranges",
			l1:       []int{1, 2, 3},
			l2:       []int{4, 5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list1 := sliceToList(tt.l1)
			list2 := sliceToList(tt.l2)
			merged := MergeTwoLists(list1, list2)
			got := listToSlice(merged)
			if len(got) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("MergeTwoLists(%v, %v) = %v, want %v", tt.l1, tt.l2, got, tt.expected)
			}
		})
	}
}

func BenchmarkMergeTwoLists(b *testing.B) {
	l1Nums := []int{1, 3, 5, 7, 9}
	l2Nums := []int{2, 4, 6, 8, 10}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		l1 := sliceToList(l1Nums)
		l2 := sliceToList(l2Nums)
		_ = MergeTwoLists(l1, l2)
	}
}
