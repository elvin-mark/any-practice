package solution

// QuickSort sorts an integer slice in ascending order in-place.
func QuickSort(arr []int) {
	if len(arr) <= 1 {
		return
	}
	sortHelper(arr, 0, len(arr)-1)
}

func sortHelper(arr []int, low, high int) {
	if low < high {
		p := partition(arr, low, high)
		sortHelper(arr, low, p-1)
		sortHelper(arr, p+1, high)
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return i
}
