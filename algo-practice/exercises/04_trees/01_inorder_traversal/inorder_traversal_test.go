package inordertraversal

import (
	"reflect"
	"testing"
)

func TestInorderTraversal(t *testing.T) {
	// Tree 1:
	//     1
	//      \
	//       2
	//      /
	//     3
	t1 := &TreeNode{
		Val: 1,
		Right: &TreeNode{
			Val:  2,
			Left: &TreeNode{Val: 3},
		},
	}

	// Tree 2: Full balanced
	//       4
	//     /   \
	//    2     6
	//   / \   / \
	//  1   3 5   7
	t2 := &TreeNode{
		Val: 4,
		Left: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 1},
			Right: &TreeNode{Val: 3},
		},
		Right: &TreeNode{
			Val:   6,
			Left:  &TreeNode{Val: 5},
			Right: &TreeNode{Val: 7},
		},
	}

	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{
			name:     "example 1",
			root:     t1,
			expected: []int{1, 3, 2},
		},
		{
			name:     "balanced tree produces sorted output for BST",
			root:     t2,
			expected: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:     "nil root",
			root:     nil,
			expected: []int{},
		},
		{
			name:     "single node",
			root:     &TreeNode{Val: 42},
			expected: []int{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InorderTraversal(tt.root)
			if len(got) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("InorderTraversal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func BenchmarkInorderTraversal(b *testing.B) {
	root := &TreeNode{
		Val: 4,
		Left: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 1},
			Right: &TreeNode{Val: 3},
		},
		Right: &TreeNode{
			Val:   6,
			Left:  &TreeNode{Val: 5},
			Right: &TreeNode{Val: 7},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = InorderTraversal(root)
	}
}
