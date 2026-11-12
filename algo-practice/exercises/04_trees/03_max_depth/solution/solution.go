package solution

// TreeNode represents a node in a binary tree.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// MaxDepth returns the maximum depth of a binary tree.
func MaxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	left := MaxDepth(root.Left)
	right := MaxDepth(root.Right)
	if left > right {
		return 1 + left
	}
	return 1 + right
}
