package maxsubarray

import "testing"

func TestMaxSubArray(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "mixed positive and negative",
			nums:     []int{-2, 1, -3, 4, -1, 2, 1, -5, 4},
			expected: 6, // [4, -1, 2, 1]
		},
		{
			name:     "single element positive",
			nums:     []int{1},
			expected: 1,
		},
		{
			name:     "all positive elements",
			nums:     []int{5, 4, -1, 7, 8},
			expected: 23,
		},
		{
			name:     "all negative elements",
			nums:     []int{-8, -3, -6, -2, -5, -4},
			expected: -2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxSubArray(tt.nums)
			if got != tt.expected {
				t.Errorf("MaxSubArray(%v) = %d, want %d", tt.nums, got, tt.expected)
			}
		})
	}
}

func BenchmarkMaxSubArray(b *testing.B) {
	nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4, 10, -3, 8, -12, 15, -4}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = MaxSubArray(nums)
	}
}
