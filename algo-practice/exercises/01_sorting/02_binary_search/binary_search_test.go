package binarysearch

import "testing"

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{
			name:     "element in middle",
			arr:      []int{1, 3, 5, 7, 9, 11},
			target:   7,
			expected: 3,
		},
		{
			name:     "first element",
			arr:      []int{2, 4, 6, 8, 10},
			target:   2,
			expected: 0,
		},
		{
			name:     "last element",
			arr:      []int{2, 4, 6, 8, 10},
			target:   10,
			expected: 4,
		},
		{
			name:     "element not present (smaller than min)",
			arr:      []int{10, 20, 30, 40},
			target:   5,
			expected: -1,
		},
		{
			name:     "element not present (greater than max)",
			arr:      []int{10, 20, 30, 40},
			target:   50,
			expected: -1,
		},
		{
			name:     "element not present (in range)",
			arr:      []int{1, 3, 5, 7, 9},
			target:   4,
			expected: -1,
		},
		{
			name:     "empty slice",
			arr:      []int{},
			target:   5,
			expected: -1,
		},
		{
			name:     "single element found",
			arr:      []int{42},
			target:   42,
			expected: 0,
		},
		{
			name:     "single element not found",
			arr:      []int{42},
			target:   99,
			expected: -1,
		},
		{
			name:     "negative numbers",
			arr:      []int{-50, -20, -5, 0, 15, 30},
			target:   -5,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BinarySearch(tt.arr, tt.target)
			if got != tt.expected {
				t.Errorf("BinarySearch(%v, %d) = %d, want %d", tt.arr, tt.target, got, tt.expected)
			}
		})
	}
}

func BenchmarkBinarySearch(b *testing.B) {
	n := 10000
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i * 2
	}
	target := 8888

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = BinarySearch(arr, target)
	}
}
