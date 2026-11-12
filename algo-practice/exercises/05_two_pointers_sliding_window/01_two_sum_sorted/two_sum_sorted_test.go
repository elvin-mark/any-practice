package twosumsorted

import (
	"reflect"
	"testing"
)

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		target   int
		expected []int
	}{
		{
			name:     "standard example",
			numbers:  []int{2, 7, 11, 15},
			target:   9,
			expected: []int{1, 2},
		},
		{
			name:     "three elements",
			numbers:  []int{2, 3, 4},
			target:   6,
			expected: []int{1, 3},
		},
		{
			name:     "negative numbers",
			numbers:  []int{-1, 0},
			target:   -1,
			expected: []int{1, 2},
		},
		{
			name:     "duplicate elements",
			numbers:  []int{0, 0, 3, 4},
			target:   0,
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwoSum(tt.numbers, tt.target)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("TwoSum(%v, %d) = %v, want %v", tt.numbers, tt.target, got, tt.expected)
			}
		})
	}
}

func BenchmarkTwoSum(b *testing.B) {
	n := 10000
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = i * 2
	}
	target := (n - 1) * 2

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = TwoSum(nums, target)
	}
}
