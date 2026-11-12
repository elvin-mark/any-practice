package solution

// MaxSubArray returns the largest sum among all non-empty contiguous subarrays.
func MaxSubArray(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	maxSoFar := nums[0]
	currentMax := nums[0]

	for i := 1; i < len(nums); i++ {
		x := nums[i]
		if x > currentMax+x {
			currentMax = x
		} else {
			currentMax = currentMax + x
		}

		if currentMax > maxSoFar {
			maxSoFar = currentMax
		}
	}

	return maxSoFar
}
