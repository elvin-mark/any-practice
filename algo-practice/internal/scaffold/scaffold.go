package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	algopractice "github.com/elvin-mark/any-practice/algo-practice"
)

// UnpackExercises extracts the embedded exercises into the target directory.
func UnpackExercises(targetDir string) error {
	subFS, err := fs.Sub(algopractice.EmbeddedExercises, "exercises")
	if err != nil {
		return fmt.Errorf("accessing embedded exercises: %w", err)
	}

	destExercises := filepath.Join(targetDir, "exercises")

	return fs.WalkDir(subFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath := filepath.Join(destExercises, path)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := fs.ReadFile(subFS, path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		// Don't overwrite existing exercise files if user has already begun
		if _, err := os.Stat(targetPath); err == nil {
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
	})
}
