package maxdepth

import "testing"

func TestMaxDepth(t *testing.T) {
	//       3
	//      / \
	//     9   20
	//        /  \
	//       15   7
	t1 := &TreeNode{
		Val:  3,
		Left: &TreeNode{Val: 9},
		Right: &TreeNode{
			Val:   20,
			Left:  &TreeNode{Val: 15},
			Right: &TreeNode{Val: 7},
		},
	}

	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name:     "depth 3",
			root:     t1,
			expected: 3,
		},
		{
			name:     "single node",
			root:     &TreeNode{Val: 1},
			expected: 1,
		},
		{
			name:     "empty tree",
			root:     nil,
			expected: 0,
		},
		{
			name: "linear chain depth 4",
			root: &TreeNode{
				Val: 1,
				Right: &TreeNode{
					Val: 2,
					Right: &TreeNode{
						Val:   3,
						Right: &TreeNode{Val: 4},
					},
				},
			},
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxDepth(tt.root)
			if got != tt.expected {
				t.Errorf("MaxDepth() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func BenchmarkMaxDepth(b *testing.B) {
	root := &TreeNode{
		Val:  3,
		Left: &TreeNode{Val: 9},
		Right: &TreeNode{
			Val:   20,
			Left:  &TreeNode{Val: 15},
			Right: &TreeNode{Val: 7},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = MaxDepth(root)
	}
}
