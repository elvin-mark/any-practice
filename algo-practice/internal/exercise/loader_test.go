package exercise

import (
	"testing"
)

func TestLoadRepository(t *testing.T) {
	repo, err := LoadRepository("../../exercises")
	if err != nil {
		t.Fatalf("failed to load repository: %v", err)
	}

	if len(repo.Exercises) == 0 {
		t.Fatalf("expected exercises to be loaded, got 0")
	}

	if len(repo.Tracks()) == 0 {
		t.Fatalf("expected tracks to be found, got 0")
	}

	first := repo.Exercises[0]
	if first.ID != "sorting_01_bubble_sort" {
		t.Errorf("expected first exercise to be sorting_01_bubble_sort, got %s", first.ID)
	}

	next := repo.Next(first.ID)
	if next == nil || next.ID != "sorting_02_binary_search" {
		t.Errorf("expected next exercise to be sorting_02_binary_search, got %v", next)
	}

	prev := repo.Prev(next.ID)
	if prev == nil || prev.ID != first.ID {
		t.Errorf("expected prev exercise to be %s, got %v", first.ID, prev)
	}
}
