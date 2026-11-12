package minstack

import "testing"

func TestMinStack(t *testing.T) {
	ms := Constructor()

	ms.Push(-2)
	ms.Push(0)
	ms.Push(-3)

	if got := ms.GetMin(); got != -3 {
		t.Fatalf("GetMin() = %d, want -3", got)
	}

	ms.Pop()

	if got := ms.Top(); got != 0 {
		t.Fatalf("Top() = %d, want 0", got)
	}

	if got := ms.GetMin(); got != -2 {
		t.Fatalf("GetMin() = %d, want -2", got)
	}

	ms.Push(-5)
	if got := ms.GetMin(); got != -5 {
		t.Fatalf("GetMin() = %d, want -5", got)
	}
}

func BenchmarkMinStack(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ms := Constructor()
		for j := 0; j < 100; j++ {
			ms.Push(j)
			_ = ms.GetMin()
		}
		for j := 0; j < 100; j++ {
			ms.Pop()
		}
	}
}
