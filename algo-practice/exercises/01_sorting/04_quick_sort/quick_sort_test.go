package quicksort

import (
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
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
			name:     "reversed",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "duplicates and negatives",
			input:    []int{3, -1, 4, -1, 5, 9, 2, 6, 5, 3, 5},
			expected: []int{-1, -1, 2, 3, 3, 4, 5, 5, 5, 6, 9},
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{77},
			expected: []int{77},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpy := make([]int, len(tt.input))
			copy(cpy, tt.input)
			QuickSort(cpy)
			if len(cpy) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(cpy, tt.expected) {
				t.Errorf("QuickSort(%v) = %v, want %v", tt.input, cpy, tt.expected)
			}
		})
	}
}

func BenchmarkQuickSort(b *testing.B) {
	sample := []int{64, 34, 25, 12, 22, 11, 90, 45, 78, 1, 99, 3, 15, 88, 23}
	arr := make([]int, len(sample))
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		copy(arr, sample)
		QuickSort(arr)
	}
}
