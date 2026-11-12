package solution

// TwoSum returns 1-based indices [index1, index2] of two numbers that add up to target.
func TwoSum(numbers []int, target int) []int {
	left := 0
	right := len(numbers) - 1

	for left < right {
		sum := numbers[left] + numbers[right]
		if sum == target {
			return []int{left + 1, right + 1}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}

	return nil
}
