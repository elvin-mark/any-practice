package embedded

import "embed"

// Dockerfile contains the base systemd Ubuntu 24.04 image recipe.
//
//go:embed Dockerfile
var Dockerfile []byte

// ScenariosFS contains all embedded core practice scenarios.
//
//go:embed scenarios/*
var ScenariosFS embed.FS
