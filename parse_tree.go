package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Kind int

const (
	KindDir Kind = iota
	KindFile
)

type Entry struct {
	Path string
	Kind Kind
}

func ParseTree(lines []string) ([]Entry, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	// Skip root line (line 0)
	var entries []Entry
	levelParent := map[int]string{0: ""}

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Remove comment
		line = cutComment(line)

		// Get indentation level
		level, rest := consumeIndent(line)

		// Remove branch decorations
		rest = trimBranch(rest)

		// Extract name
		name := strings.TrimSpace(rest)
		if name == "" {
			continue
		}

		// Check if directory
		isDir := strings.HasSuffix(name, "/")
		name = strings.TrimSuffix(name, "/")

		// Build relative path
		parent := levelParent[level]
		var relPath string
		if parent == "" {
			relPath = name
		} else {
			relPath = filepath.Join(parent, name)
		}

		if isDir {
			entries = append(entries, Entry{Path: relPath, Kind: KindDir})
			levelParent[level+1] = relPath
			// Clear deeper levels (sibling branches)
			for k := range levelParent {
				if k > level+1 {
					delete(levelParent, k)
				}
			}
		} else {
			entries = append(entries, Entry{Path: relPath, Kind: KindFile})
		}
	}

	return entries, nil
}

func cutComment(s string) string {
	// Priority: " # " (space + hash + space)
	if idx := strings.Index(s, " # "); idx >= 0 {
		return s[:idx]
	}

	// Otherwise, check if # is preceded by space, tab, │, or |
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '#' {
			if i > 0 {
				prev := runes[i-1]
				if prev == ' ' || prev == '\t' || prev == '│' || prev == '|' {
					return string(runes[:i])
				}
			}
		}
	}

	return s
}

func consumeIndent(s string) (level int, rest string) {
	runes := []rune(s)
	pos := 0

	for pos < len(runes) {
		// Try '│  ' (U+2502 + 2 spaces)
		if pos+2 < len(runes) && runes[pos] == '│' && runes[pos+1] == ' ' && runes[pos+2] == ' ' {
			level++
			pos += 3
			continue
		}

		// Try '|  ' (pipe + 2 spaces)
		if pos+2 < len(runes) && runes[pos] == '|' && runes[pos+1] == ' ' && runes[pos+2] == ' ' {
			level++
			pos += 3
			continue
		}

		// Try '   ' (3 spaces)
		if pos+2 < len(runes) && runes[pos] == ' ' && runes[pos+1] == ' ' && runes[pos+2] == ' ' {
			level++
			pos += 3
			continue
		}

		// Try '\t' (single tab)
		if runes[pos] == '\t' {
			level++
			pos++
			continue
		}

		// Try single '│' or '|' followed by spaces
		if runes[pos] == '│' || runes[pos] == '|' {
			start := pos
			pos++
			// Consume trailing spaces
			for pos < len(runes) && runes[pos] == ' ' {
				pos++
			}
			if pos > start+1 {
				// Had spaces after │ or |, count as one level
				level++
				continue
			}
			// No spaces, revert
			pos = start
			break
		}

		// Check if it's a leading space that doesn't form a pattern
		if runes[pos] == ' ' {
			// Skip single spaces that don't form part of an indent token
			pos++
			continue
		}

		// No more indent tokens
		break
	}

	return level, string(runes[pos:])
}

func trimBranch(s string) string {
	runes := []rune(s)

	// Remove branch prefixes
	branches := []string{"├─", "└─", "|--", "`--", "+--"}
	for {
		trimmed := false
		for _, branch := range branches {
			branchRunes := []rune(branch)
			if len(runes) >= len(branchRunes) {
				match := true
				for i, br := range branchRunes {
					if runes[i] != br {
						match = false
						break
					}
				}
				if match {
					runes = runes[len(branchRunes):]
					trimmed = true
					break
				}
			}
		}
		if !trimmed {
			break
		}
	}

	// Remove residual box drawing and dashes
	for len(runes) > 0 {
		r := runes[0]
		if r == '─' || r == '—' || r == '│' || r == '|' || r == '-' {
			runes = runes[1:]
		} else {
			break
		}
	}

	// Trim left spaces
	result := string(runes)
	return strings.TrimLeft(result, " ")
}