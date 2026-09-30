package migrations

import "embed"

// Files keeps migrations available when the CLI runs outside the project.
//
//go:embed sqlite/*.sql
var Files embed.FS
