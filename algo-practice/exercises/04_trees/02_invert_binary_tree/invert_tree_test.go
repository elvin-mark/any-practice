package inverttree

import "testing"

func isSameTree(p, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}
	return p.Val == q.Val && isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)
}

func TestInvertTree(t *testing.T) {
	// Original:
	//      4
	//    /   \
	//   2     7
	//  / \   / \
	// 1   3 6   9
	tree := &TreeNode{
		Val: 4,
		Left: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 1},
			Right: &TreeNode{Val: 3},
		},
		Right: &TreeNode{
			Val:   7,
			Left:  &TreeNode{Val: 6},
			Right: &TreeNode{Val: 9},
		},
	}

	// Inverted:
	//      4
	//    /   \
	//   7     2
	//  / \   / \
	// 9   6 3   1
	expected := &TreeNode{
		Val: 4,
		Left: &TreeNode{
			Val:   7,
			Left:  &TreeNode{Val: 9},
			Right: &TreeNode{Val: 6},
		},
		Right: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 3},
			Right: &TreeNode{Val: 1},
		},
	}

	inverted := InvertTree(tree)
	if !isSameTree(inverted, expected) {
		t.Fatalf("InvertTree did not match expected mirrored tree")
	}

	if InvertTree(nil) != nil {
		t.Fatalf("InvertTree(nil) should return nil")
	}
}

func BenchmarkInvertTree(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		root := &TreeNode{
			Val: 4,
			Left: &TreeNode{
				Val:   2,
				Left:  &TreeNode{Val: 1},
				Right: &TreeNode{Val: 3},
			},
			Right: &TreeNode{
				Val:   7,
				Left:  &TreeNode{Val: 6},
				Right: &TreeNode{Val: 9},
			},
		}
		_ = InvertTree(root)
	}
}
