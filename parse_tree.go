package main

import (
	"errors"
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
		return nil, errors.New("empty tree")
	}

	// Check root line validity
	root := strings.TrimSpace(lines[0])
	root = strings.TrimSuffix(root, "/")
	if root == "" {
		return nil, errors.New("invalid root line")
	}

	// Skip root line (line 0)
	var entries []Entry = make([]Entry, 0) // Initialize as empty slice, not nil
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
	// Priority: " #" (space + hash)
	if idx := strings.Index(s, " #"); idx >= 0 {
		return s[:idx]
	}

	// Otherwise, check if # is preceded by space, tab, │, or |
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '#' {
			if i == 0 {
				return ""
			}
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

// Helper functions for consumeIndent
func matchUnicodeBoxPattern(runes []rune, pos int) bool {
	return pos+2 < len(runes) && runes[pos] == '│' && runes[pos+1] == ' ' && runes[pos+2] == ' '
}

func matchPipePattern(runes []rune, pos int) bool {
	return pos+2 < len(runes) && runes[pos] == '|' && runes[pos+1] == ' ' && runes[pos+2] == ' '
}

func matchThreeSpacePattern(runes []rune, pos int) bool {
	return pos+2 < len(runes) && runes[pos] == ' ' && runes[pos+1] == ' ' && runes[pos+2] == ' '
}

func consumeSinglePipeWithSpaces(runes []rune, pos int) (int, bool) {
	if pos >= len(runes) {
		return pos, false
	}
	
	if runes[pos] == '│' || runes[pos] == '|' {
		start := pos
		pos++
		// Consume trailing spaces
		for pos < len(runes) && runes[pos] == ' ' {
			pos++
		}
		if pos > start+1 {
			// Had spaces after │ or |, count as one level
			return pos, true
		}
		// No spaces, revert
		return start, false
	}
	
	return pos, false
}

func consumeIndent(s string) (level int, rest string) {
	runes := []rune(s)
	pos := 0

	for pos < len(runes) {
		// Try '│  ' (U+2502 + 2 spaces)
		if matchUnicodeBoxPattern(runes, pos) {
			level++
			pos += 3
			continue
		}

		// Try '|  ' (pipe + 2 spaces)
		if matchPipePattern(runes, pos) {
			level++
			pos += 3
			continue
		}

		// Try '   ' (3 spaces)
		if matchThreeSpacePattern(runes, pos) {
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
		if newPos, found := consumeSinglePipeWithSpaces(runes, pos); found {
			level++
			pos = newPos
			continue
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

// Helper functions for trimBranch
func trimLeadingWhitespace(runes []rune) []rune {
	for len(runes) > 0 && (runes[0] == ' ' || runes[0] == '\t') {
		runes = runes[1:]
	}
	return runes
}

func removeBranchPrefixes(runes []rune) []rune {
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
					// After removing a branch token, trim spaces/tabs again
					runes = trimLeadingWhitespace(runes)
					break
				}
			}
		}
		if !trimmed {
			break
		}
	}
	return runes
}

func removeResidualBoxDrawing(runes []rune) []rune {
	for len(runes) > 0 {
		r := runes[0]
		if r == '─' || r == '—' || r == '│' || r == '|' || r == '-' || r == ' ' || r == '\t' {
			runes = runes[1:]
		} else {
			break
		}
	}
	return runes
}

func trimBranch(s string) string {
	runes := []rune(s)

	// First, trim any leading spaces or tabs
	runes = trimLeadingWhitespace(runes)

	// Remove branch prefixes repeatedly
	runes = removeBranchPrefixes(runes)

	// Remove residual box drawing characters, dashes, pipes, and spaces between them
	runes = removeResidualBoxDrawing(runes)

	// Final trim of any remaining leading spaces
	result := string(runes)
	return strings.TrimLeft(result, " \t")
}
