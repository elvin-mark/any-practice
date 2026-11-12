package solution

// TreeNode represents a node in a binary tree.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// InorderTraversal returns the values of the nodes in in-order.
func InorderTraversal(root *TreeNode) []int {
	var result []int
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}
		dfs(node.Left)
		result = append(result, node.Val)
		dfs(node.Right)
	}
	dfs(root)
	if result == nil {
		return []int{}
	}
	return result
}
