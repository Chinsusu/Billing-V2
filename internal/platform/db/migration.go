package db

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var migrationFileNamePattern = regexp.MustCompile(`^([0-9]{4})_([a-z0-9_]+)\.sql$`)

type Migration struct {
	Version  string
	Name     string
	Path     string
	SQL      string
	Checksum string
}

func LoadMigrations(fsys fs.FS) ([]Migration, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}

	migrations := make([]Migration, 0, len(entries))
	seenVersions := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		matches := migrationFileNamePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			return nil, fmt.Errorf("migration %q must match 0001_descriptive_name.sql", entry.Name())
		}

		version := matches[1]
		name := matches[2]
		if previous, exists := seenVersions[version]; exists {
			return nil, fmt.Errorf("duplicate migration version %s in %s and %s", version, previous, entry.Name())
		}
		seenVersions[version] = entry.Name()

		body, err := fs.ReadFile(fsys, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read migration %q: %w", entry.Name(), err)
		}
		sqlText := strings.TrimSpace(string(body))
		if sqlText == "" {
			return nil, fmt.Errorf("migration %q is empty", entry.Name())
		}

		migrations = append(migrations, Migration{
			Version:  version,
			Name:     name,
			Path:     filepath.ToSlash(entry.Name()),
			SQL:      sqlText,
			Checksum: checksum(sqlText),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	return migrations, nil
}

func checksum(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
