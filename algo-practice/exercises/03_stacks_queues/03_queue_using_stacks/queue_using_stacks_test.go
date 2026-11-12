package queueusingstacks

import "testing"

func TestMyQueue(t *testing.T) {
	q := Constructor()
	if !q.Empty() {
		t.Fatalf("expected queue to be empty initially")
	}

	q.Push(1)
	q.Push(2)

	if got := q.Peek(); got != 1 {
		t.Fatalf("Peek() = %d, want 1", got)
	}

	if got := q.Pop(); got != 1 {
		t.Fatalf("Pop() = %d, want 1", got)
	}

	if q.Empty() {
		t.Fatalf("expected queue to not be empty after 1 pop of 2 elements")
	}

	q.Push(3)
	if got := q.Pop(); got != 2 {
		t.Fatalf("Pop() = %d, want 2", got)
	}

	if got := q.Pop(); got != 3 {
		t.Fatalf("Pop() = %d, want 3", got)
	}

	if !q.Empty() {
		t.Fatalf("expected queue to be empty after popping all elements")
	}
}

func BenchmarkMyQueue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		q := Constructor()
		for j := 0; j < 100; j++ {
			q.Push(j)
		}
		for j := 0; j < 100; j++ {
			_ = q.Pop()
		}
	}
}
