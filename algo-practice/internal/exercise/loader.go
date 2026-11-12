package exercise

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
)

// Repository manages a collection of exercises.
type Repository struct {
	Exercises []*Exercise
	byID      map[string]*Exercise
	tracks    []string
}

// LoadRepository discovers and loads all exercises rooted at exercisesDir.
func LoadRepository(exercisesDir string) (*Repository, error) {
	var exercises []*Exercise
	byID := make(map[string]*Exercise)

	err := filepath.WalkDir(exercisesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "exercise.json" {
			dir := filepath.Dir(path)
			ex, err := LoadFromDir(dir)
			if err != nil {
				return fmt.Errorf("loading exercise at %s: %w", dir, err)
			}
			exercises = append(exercises, ex)
			byID[ex.ID] = ex
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking exercises dir: %w", err)
	}

	// Sort exercises consistently: TrackOrder -> Order -> ID
	sort.Slice(exercises, func(i, j int) bool {
		if exercises[i].TrackOrder != exercises[j].TrackOrder {
			return exercises[i].TrackOrder < exercises[j].TrackOrder
		}
		if exercises[i].Order != exercises[j].Order {
			return exercises[i].Order < exercises[j].Order
		}
		return exercises[i].ID < exercises[j].ID
	})

	trackSeen := make(map[string]bool)
	var tracks []string
	for _, ex := range exercises {
		if !trackSeen[ex.Track] {
			trackSeen[ex.Track] = true
			tracks = append(tracks, ex.Track)
		}
	}

	return &Repository{
		Exercises: exercises,
		byID:      byID,
		tracks:    tracks,
	}, nil
}

// FindByID returns an exercise by its unique ID.
func (r *Repository) FindByID(id string) (*Exercise, bool) {
	ex, ok := r.byID[id]
	return ex, ok
}

// Tracks returns the ordered list of distinct track names.
func (r *Repository) Tracks() []string {
	return r.tracks
}

// ByTrack returns all exercises belonging to a track.
func (r *Repository) ByTrack(track string) []*Exercise {
	var result []*Exercise
	for _, ex := range r.Exercises {
		if ex.Track == track {
			result = append(result, ex)
		}
	}
	return result
}

// Next returns the next exercise following the given ID, or nil if at the end.
func (r *Repository) Next(currentID string) *Exercise {
	for i, ex := range r.Exercises {
		if ex.ID == currentID {
			if i+1 < len(r.Exercises) {
				return r.Exercises[i+1]
			}
			return nil
		}
	}
	if len(r.Exercises) > 0 {
		return r.Exercises[0]
	}
	return nil
}

// Prev returns the previous exercise before the given ID, or nil if at start.
func (r *Repository) Prev(currentID string) *Exercise {
	for i, ex := range r.Exercises {
		if ex.ID == currentID {
			if i > 0 {
				return r.Exercises[i-1]
			}
			return nil
		}
	}
	return nil
}

// IndexOf returns the 0-based index of an exercise, or -1 if not found.
func (r *Repository) IndexOf(id string) int {
	for i, ex := range r.Exercises {
		if ex.ID == id {
			return i
		}
	}
	return -1
}
