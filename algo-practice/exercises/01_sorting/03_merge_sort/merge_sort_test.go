package mergesort

import (
	"reflect"
	"testing"
)

func TestMergeSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "already sorted",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "reverse order",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "unsorted with duplicates",
			input:    []int{9, 3, 7, 5, 6, 4, 8, 2, 1, 5},
			expected: []int{1, 2, 3, 4, 5, 5, 6, 7, 8, 9},
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{100},
			expected: []int{100},
		},
		{
			name:     "negative values",
			input:    []int{-12, 10, -3, 0, -50, 4},
			expected: []int{-50, -12, -3, 0, 4, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeSort(tt.input)
			if len(got) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("MergeSort(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func BenchmarkMergeSort(b *testing.B) {
	sample := []int{64, 34, 25, 12, 22, 11, 90, 45, 78, 1, 99, 3, 15, 88, 23}
	arr := make([]int, len(sample))
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		copy(arr, sample)
		_ = MergeSort(arr)
	}
}
