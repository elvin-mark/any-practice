package algopractice

import "embed"

// EmbeddedExercises embeds the entire exercises directory into the compiled binary.
//
//go:embed all:exercises
var EmbeddedExercises embed.FS
