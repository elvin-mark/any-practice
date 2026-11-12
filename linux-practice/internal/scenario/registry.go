package scenario

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/elvin-mark/any-practice/linux-practice/internal/config"
)

// Registry manages the catalog of embedded and external scenarios.
type Registry struct {
	scenarios map[string]*Scenario
	byTrack   map[string][]*Scenario
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{
		scenarios: make(map[string]*Scenario),
		byTrack:   make(map[string][]*Scenario),
	}
}

// Register adds a scenario to the registry.
func (r *Registry) Register(sc *Scenario) error {
	if err := sc.Validate(); err != nil {
		return fmt.Errorf("invalid scenario %q: %w", sc.ID, err)
	}

	r.scenarios[sc.ID] = sc
	trackKey := strings.ToLower(sc.Track)
	r.byTrack[trackKey] = append(r.byTrack[trackKey], sc)
	return nil
}

// Get returns a scenario by its unique identifier.
func (r *Registry) Get(id string) (*Scenario, bool) {
	sc, ok := r.scenarios[id]
	return sc, ok
}

// All returns all registered scenarios sorted alphabetically by track and ID.
func (r *Registry) All() []*Scenario {
	var list []*Scenario
	for _, sc := range r.scenarios {
		list = append(list, sc)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Track == list[j].Track {
			return list[i].ID < list[j].ID
		}
		return list[i].Track < list[j].Track
	})
	return list
}

// Tracks returns the sorted list of unique tracks.
func (r *Registry) Tracks() []string {
	var tracks []string
	for t := range r.byTrack {
		tracks = append(tracks, t)
	}
	sort.Strings(tracks)
	return tracks
}

// ByTrack returns all scenarios belonging to a specific track.
func (r *Registry) ByTrack(track string) []*Scenario {
	return r.byTrack[strings.ToLower(track)]
}

// LoadEmbedded discovers and loads scenarios from an embedded file system.
func (r *Registry) LoadEmbedded(embeddedFS fs.FS, rootDir string) error {
	return fs.WalkDir(embeddedFS, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			manifest := filepath.Join(path, "scenario.yaml")
			if _, err := fs.Stat(embeddedFS, manifest); err == nil {
				sc, err := LoadFromFS(embeddedFS, path)
				if err != nil {
					return fmt.Errorf("loading embedded scenario at %s: %w", path, err)
				}
				if err := r.Register(sc); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// LoadExternal discovers and loads scenarios from a local filesystem directory.
func (r *Registry) LoadExternal(rootDir string) error {
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			manifest := filepath.Join(path, "scenario.yaml")
			if _, err := os.Stat(manifest); err == nil {
				sc, err := LoadFromDisk(path)
				if err != nil {
					return fmt.Errorf("loading external scenario at %s: %w", path, err)
				}
				// External scenario overwrites or registers
				if err := r.Register(sc); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// LoadAll loads both embedded scenarios and user custom scenarios from ~/.config/linux-practice/scenarios.
func (r *Registry) LoadAll(embeddedFS fs.FS) error {
	if err := r.LoadEmbedded(embeddedFS, "scenarios"); err != nil {
		return fmt.Errorf("loading embedded scenarios: %w", err)
	}

	customDir, err := config.ScenariosDir()
	if err == nil {
		_ = r.LoadExternal(customDir)
	}

	return nil
}
