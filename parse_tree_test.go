package main

import (
	"reflect"
	"testing"
)

func TestParseTree(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected []Entry
		hasError bool
	}{
		{
			name:     "empty lines",
			lines:    []string{},
			expected: nil,
			hasError: true,
		},
		{
			name:     "only root",
			lines:    []string{"myapp/"},
			expected: []Entry{},
			hasError: false,
		},
		{
			name:     "invalid root",
			lines:    []string{""},
			expected: nil,
			hasError: true,
		},
		{
			name: "sample myapp tree with unicode box drawing",
			lines: []string{
				"myapp/",
				"├─ src/",
				"│  ├─ handlers/",
				"│  │  ├─ user.go",
				"│  │  └─ auth.go",
				"│  ├─ models/",
				"│  │  └─ user.go",
				"│  ├─ middleware/",
				"│  │  └─ logger.go",
				"│  └─ main.go",
				"├─ tests/",
				"│  ├─ user_test.go",
				"│  └─ auth_test.go",
				"├─ config/",
				"│  └─ config.yaml",
				"├─ .env",
				"├─ .gitignore",
				"├─ go.mod",
				"├─ Dockerfile",
				"└─ README.md",
			},
			expected: []Entry{
				{Path: "src", Kind: KindDir},
				{Path: "src/handlers", Kind: KindDir},
				{Path: "src/handlers/user.go", Kind: KindFile},
				{Path: "src/handlers/auth.go", Kind: KindFile},
				{Path: "src/models", Kind: KindDir},
				{Path: "src/models/user.go", Kind: KindFile},
				{Path: "src/middleware", Kind: KindDir},
				{Path: "src/middleware/logger.go", Kind: KindFile},
				{Path: "src/main.go", Kind: KindFile},
				{Path: "tests", Kind: KindDir},
				{Path: "tests/user_test.go", Kind: KindFile},
				{Path: "tests/auth_test.go", Kind: KindFile},
				{Path: "config", Kind: KindDir},
				{Path: "config/config.yaml", Kind: KindFile},
				{Path: ".env", Kind: KindFile},
				{Path: ".gitignore", Kind: KindFile},
				{Path: "go.mod", Kind: KindFile},
				{Path: "Dockerfile", Kind: KindFile},
				{Path: "README.md", Kind: KindFile},
			},
			hasError: false,
		},
		{
			name: "same tree with ASCII box drawing",
			lines: []string{
				"myapp/",
				"|-- src/",
				"|   |-- handlers/",
				"|   |   |-- user.go",
				"|   |   `-- auth.go",
				"|   |-- models/",
				"|   |   `-- user.go",
				"|   |-- middleware/",
				"|   |   `-- logger.go",
				"|   `-- main.go",
				"|-- tests/",
				"|   |-- user_test.go",
				"|   `-- auth_test.go",
				"|-- config/",
				"|   `-- config.yaml",
				"|-- .env",
				"|-- .gitignore",
				"|-- go.mod",
				"|-- Dockerfile",
				"`-- README.md",
			},
			expected: []Entry{
				{Path: "src", Kind: KindDir},
				{Path: "src/handlers", Kind: KindDir},
				{Path: "src/handlers/user.go", Kind: KindFile},
				{Path: "src/handlers/auth.go", Kind: KindFile},
				{Path: "src/models", Kind: KindDir},
				{Path: "src/models/user.go", Kind: KindFile},
				{Path: "src/middleware", Kind: KindDir},
				{Path: "src/middleware/logger.go", Kind: KindFile},
				{Path: "src/main.go", Kind: KindFile},
				{Path: "tests", Kind: KindDir},
				{Path: "tests/user_test.go", Kind: KindFile},
				{Path: "tests/auth_test.go", Kind: KindFile},
				{Path: "config", Kind: KindDir},
				{Path: "config/config.yaml", Kind: KindFile},
				{Path: ".env", Kind: KindFile},
				{Path: ".gitignore", Kind: KindFile},
				{Path: "go.mod", Kind: KindFile},
				{Path: "Dockerfile", Kind: KindFile},
				{Path: "README.md", Kind: KindFile},
			},
			hasError: false,
		},
		{
			name: "same tree with 3-space indentation",
			lines: []string{
				"myapp/",
				"   src/",
				"      handlers/",
				"         user.go",
				"         auth.go",
				"      models/",
				"         user.go",
				"      middleware/",
				"         logger.go",
				"      main.go",
				"   tests/",
				"      user_test.go",
				"      auth_test.go",
				"   config/",
				"      config.yaml",
				"   .env",
				"   .gitignore",
				"   go.mod",
				"   Dockerfile",
				"   README.md",
			},
			expected: []Entry{
				{Path: "src", Kind: KindDir},
				{Path: "src/handlers", Kind: KindDir},
				{Path: "src/handlers/user.go", Kind: KindFile},
				{Path: "src/handlers/auth.go", Kind: KindFile},
				{Path: "src/models", Kind: KindDir},
				{Path: "src/models/user.go", Kind: KindFile},
				{Path: "src/middleware", Kind: KindDir},
				{Path: "src/middleware/logger.go", Kind: KindFile},
				{Path: "src/main.go", Kind: KindFile},
				{Path: "tests", Kind: KindDir},
				{Path: "tests/user_test.go", Kind: KindFile},
				{Path: "tests/auth_test.go", Kind: KindFile},
				{Path: "config", Kind: KindDir},
				{Path: "config/config.yaml", Kind: KindFile},
				{Path: ".env", Kind: KindFile},
				{Path: ".gitignore", Kind: KindFile},
				{Path: "go.mod", Kind: KindFile},
				{Path: "Dockerfile", Kind: KindFile},
				{Path: "README.md", Kind: KindFile},
			},
			hasError: false,
		},
		{
			name: "same tree with tab indentation",
			lines: []string{
				"myapp/",
				"\tsrc/",
				"\t\thandlers/",
				"\t\t\tuser.go",
				"\t\t\tauth.go",
				"\t\tmodels/",
				"\t\t\tuser.go",
				"\t\tmiddleware/",
				"\t\t\tlogger.go",
				"\t\tmain.go",
				"\ttests/",
				"\t\tuser_test.go",
				"\t\tauth_test.go",
				"\tconfig/",
				"\t\tconfig.yaml",
				"\t.env",
				"\t.gitignore",
				"\tgo.mod",
				"\tDockerfile",
				"\tREADME.md",
			},
			expected: []Entry{
				{Path: "src", Kind: KindDir},
				{Path: "src/handlers", Kind: KindDir},
				{Path: "src/handlers/user.go", Kind: KindFile},
				{Path: "src/handlers/auth.go", Kind: KindFile},
				{Path: "src/models", Kind: KindDir},
				{Path: "src/models/user.go", Kind: KindFile},
				{Path: "src/middleware", Kind: KindDir},
				{Path: "src/middleware/logger.go", Kind: KindFile},
				{Path: "src/main.go", Kind: KindFile},
				{Path: "tests", Kind: KindDir},
				{Path: "tests/user_test.go", Kind: KindFile},
				{Path: "tests/auth_test.go", Kind: KindFile},
				{Path: "config", Kind: KindDir},
				{Path: "config/config.yaml", Kind: KindFile},
				{Path: ".env", Kind: KindFile},
				{Path: ".gitignore", Kind: KindFile},
				{Path: "go.mod", Kind: KindFile},
				{Path: "Dockerfile", Kind: KindFile},
				{Path: "README.md", Kind: KindFile},
			},
			hasError: false,
		},
		{
			name: "tree with hash in filename (not comment)",
			lines: []string{
				"project/",
				"├─ docs/",
				"│  ├─ v1#draft.md",
				"│  └─ api#spec.yaml",
				"└─ src/",
				"   └─ hash#file.go",
			},
			expected: []Entry{
				{Path: "docs", Kind: KindDir},
				{Path: "docs/v1#draft.md", Kind: KindFile},
				{Path: "docs/api#spec.yaml", Kind: KindFile},
				{Path: "src", Kind: KindDir},
				{Path: "src/hash#file.go", Kind: KindFile},
			},
			hasError: false,
		},
		{
			name: "tree with comments to be removed",
			lines: []string{
				"project/",
				"├─ src/ # source code directory",
				"│  ├─ main.go # entry point",
				"│  └─ utils.go # utility functions",
				"├─ docs/ # documentation",
				"│  └─ README.md # project readme",
				"└─ tests/ # test files",
				"   └─ main_test.go # main tests",
			},
			expected: []Entry{
				{Path: "src", Kind: KindDir},
				{Path: "src/main.go", Kind: KindFile},
				{Path: "src/utils.go", Kind: KindFile},
				{Path: "docs", Kind: KindDir},
				{Path: "docs/README.md", Kind: KindFile},
				{Path: "tests", Kind: KindDir},
				{Path: "tests/main_test.go", Kind: KindFile},
			},
			hasError: false,
		},
		{
			name: "tree with empty lines and whitespace",
			lines: []string{
				"project/",
				"",
				"├─ src/",
				"",
				"│  └─ main.go",
				"",
				"└─ README.md",
				"",
			},
			expected: []Entry{
				{Path: "src", Kind: KindDir},
				{Path: "src/main.go", Kind: KindFile},
				{Path: "README.md", Kind: KindFile},
			},
			hasError: false,
		},
		{
			name: "mixed branch token styles",
			lines: []string{
				"project/",
				"+-- src/",
				"|   +-- handlers/",
				"|   |   |-- user.go",
				"|   |   `-- auth.go",
				"|   `-- main.go",
				"`-- README.md",
			},
			expected: []Entry{
				{Path: "src", Kind: KindDir},
				{Path: "src/handlers", Kind: KindDir},
				{Path: "src/handlers/user.go", Kind: KindFile},
				{Path: "src/handlers/auth.go", Kind: KindFile},
				{Path: "src/main.go", Kind: KindFile},
				{Path: "README.md", Kind: KindFile},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTree(tt.lines)

			if tt.hasError {
				if err == nil {
					t.Errorf("ParseTree() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseTree() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseTree() mismatch:\n%s", cmpEntries(tt.expected, result))
			}
		})
	}
}

func TestCutComment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no comment",
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			name:     "comment with space hash",
			input:    "file.txt # this is a comment",
			expected: "file.txt",
		},
		{
			name:     "comment with tab hash",
			input:    "file.txt\t# this is a comment",
			expected: "file.txt\t",
		},
		{
			name:     "comment with pipe hash",
			input:    "file.txt|# this is a comment",
			expected: "file.txt|",
		},
		{
			name:     "comment with box drawing hash",
			input:    "file.txt│# this is a comment",
			expected: "file.txt│",
		},
		{
			name:     "hash at beginning",
			input:    "# comment only",
			expected: "",
		},
		{
			name:     "hash in filename (no preceding space/tab/pipe)",
			input:    "file#name.txt",
			expected: "file#name.txt",
		},
		{
			name:     "hash in filename with path",
			input:    "docs/v1#draft.md",
			expected: "docs/v1#draft.md",
		},
		{
			name:     "multiple hashes with space",
			input:    "file.txt # comment # more",
			expected: "file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cutComment(tt.input)
			if result != tt.expected {
				t.Errorf("cutComment(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConsumeIndent(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedLevel int
		expectedRest  string
	}{
		{
			name:          "no indent",
			input:         "file.txt",
			expectedLevel: 0,
			expectedRest:  "file.txt",
		},
		{
			name:          "single level with unicode box drawing and spaces",
			input:         "│  file.txt",
			expectedLevel: 1,
			expectedRest:  "file.txt",
		},
		{
			name:          "single level with pipe and spaces",
			input:         "|  file.txt",
			expectedLevel: 1,
			expectedRest:  "file.txt",
		},
		{
			name:          "single level with 3 spaces",
			input:         "   file.txt",
			expectedLevel: 1,
			expectedRest:  "file.txt",
		},
		{
			name:          "single level with tab",
			input:         "\tfile.txt",
			expectedLevel: 1,
			expectedRest:  "file.txt",
		},
		{
			name:          "multiple levels with unicode",
			input:         "│  │  file.txt",
			expectedLevel: 2,
			expectedRest:  "file.txt",
		},
		{
			name:          "multiple levels with pipes",
			input:         "|  |  file.txt",
			expectedLevel: 2,
			expectedRest:  "file.txt",
		},
		{
			name:          "multiple levels with 3-space blocks",
			input:         "      file.txt",
			expectedLevel: 2,
			expectedRest:  "file.txt",
		},
		{
			name:          "mixed indent types",
			input:         "│  \tfile.txt",
			expectedLevel: 2,
			expectedRest:  "file.txt",
		},
		{
			name:          "lone pipe with trailing spaces",
			input:         "|    file.txt",
			expectedLevel: 1,
			expectedRest:  "file.txt",
		},
		{
			name:          "lone unicode pipe with trailing spaces",
			input:         "│    file.txt",
			expectedLevel: 1,
			expectedRest:  "file.txt",
		},
		{
			name:          "stray leading spaces before branch token",
			input:         "  ├─ file.txt",
			expectedLevel: 0,
			expectedRest:  "├─ file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, rest := consumeIndent(tt.input)
			if level != tt.expectedLevel {
				t.Errorf("consumeIndent(%q) level = %d, want %d", tt.input, level, tt.expectedLevel)
			}
			if rest != tt.expectedRest {
				t.Errorf("consumeIndent(%q) rest = %q, want %q", tt.input, rest, tt.expectedRest)
			}
		})
	}
}

// Test helper functions for consumeIndent
func TestMatchUnicodeBoxPattern(t *testing.T) {
	tests := []struct {
		name     string
		runes    []rune
		pos      int
		expected bool
	}{
		{
			name:     "matches unicode box pattern",
			runes:    []rune("│  file.txt"),
			pos:      0,
			expected: true,
		},
		{
			name:     "no match - insufficient length",
			runes:    []rune("│ "),
			pos:      0,
			expected: false,
		},
		{
			name:     "no match - wrong pattern",
			runes:    []rune("│x file.txt"),
			pos:      0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchUnicodeBoxPattern(tt.runes, tt.pos)
			if result != tt.expected {
				t.Errorf("matchUnicodeBoxPattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMatchPipePattern(t *testing.T) {
	tests := []struct {
		name     string
		runes    []rune
		pos      int
		expected bool
	}{
		{
			name:     "matches pipe pattern",
			runes:    []rune("|  file.txt"),
			pos:      0,
			expected: true,
		},
		{
			name:     "no match - insufficient length",
			runes:    []rune("| "),
			pos:      0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchPipePattern(tt.runes, tt.pos)
			if result != tt.expected {
				t.Errorf("matchPipePattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMatchThreeSpacePattern(t *testing.T) {
	tests := []struct {
		name     string
		runes    []rune
		pos      int
		expected bool
	}{
		{
			name:     "matches three space pattern",
			runes:    []rune("   file.txt"),
			pos:      0,
			expected: true,
		},
		{
			name:     "no match - insufficient length",
			runes:    []rune("  "),
			pos:      0,
			expected: false,
		},
		{
			name:     "no match - mixed characters",
			runes:    []rune("  x"),
			pos:      0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchThreeSpacePattern(tt.runes, tt.pos)
			if result != tt.expected {
				t.Errorf("matchThreeSpacePattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConsumeSinglePipeWithSpaces(t *testing.T) {
	tests := []struct {
		name         string
		runes        []rune
		pos          int
		expectedPos  int
		expectedFound bool
	}{
		{
			name:         "pipe with spaces",
			runes:        []rune("|   file.txt"),
			pos:          0,
			expectedPos:  4,
			expectedFound: true,
		},
		{
			name:         "unicode pipe with spaces",
			runes:        []rune("│   file.txt"),
			pos:          0,
			expectedPos:  4,
			expectedFound: true,
		},
		{
			name:         "pipe without spaces",
			runes:        []rune("|file.txt"),
			pos:          0,
			expectedPos:  0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newPos, found := consumeSinglePipeWithSpaces(tt.runes, tt.pos)
			if newPos != tt.expectedPos {
				t.Errorf("consumeSinglePipeWithSpaces() pos = %d, want %d", newPos, tt.expectedPos)
			}
			if found != tt.expectedFound {
				t.Errorf("consumeSinglePipeWithSpaces() found = %v, want %v", found, tt.expectedFound)
			}
		})
	}
}

func TestTrimBranch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no branch token",
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			name:     "with unicode branch ├─",
			input:    "├─ file.txt",
			expected: "file.txt",
		},
		{
			name:     "with unicode branch └─",
			input:    "└─ file.txt",
			expected: "file.txt",
		},
		{
			name:     "with ASCII branch |--",
			input:    "|-- file.txt",
			expected: "file.txt",
		},
		{
			name:     "with ASCII branch `--",
			input:    "`-- file.txt",
			expected: "file.txt",
		},
		{
			name:     "with ASCII branch +--",
			input:    "+-- file.txt",
			expected: "file.txt",
		},
		{
			name:     "multiple branch tokens",
			input:    "├─└─ file.txt",
			expected: "file.txt",
		},
		{
			name:     "with leading spaces",
			input:    "  ├─ file.txt",
			expected: "file.txt",
		},
		{
			name:     "with box drawing residue",
			input:    "├─────── file.txt",
			expected: "file.txt",
		},
		{
			name:     "with pipes and unicode box drawing",
			input:    "│ │ file.txt",
			expected: "file.txt",
		},
		{
			name:     "dot files should not be trimmed",
			input:    ".gitignore",
			expected: ".gitignore",
		},
		{
			name:     "branch token with dot file",
			input:    "├─ .env",
			expected: ".env",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimBranch(tt.input)
			if result != tt.expected {
				t.Errorf("trimBranch(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Test helper functions for trimBranch
func TestTrimLeadingWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		expected []rune
	}{
		{
			name:     "no leading whitespace",
			input:    []rune("file.txt"),
			expected: []rune("file.txt"),
		},
		{
			name:     "spaces and tabs",
			input:    []rune("  \tfile.txt"),
			expected: []rune("file.txt"),
		},
		{
			name:     "only whitespace",
			input:    []rune("   "),
			expected: []rune(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimLeadingWhitespace(tt.input)
			if string(result) != string(tt.expected) {
				t.Errorf("trimLeadingWhitespace(%q) = %q, want %q", string(tt.input), string(result), string(tt.expected))
			}
		})
	}
}

func TestRemoveBranchPrefixes(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		expected []rune
	}{
		{
			name:     "unicode branch",
			input:    []rune("├─ file.txt"),
			expected: []rune("file.txt"),
		},
		{
			name:     "ASCII branch",
			input:    []rune("|-- file.txt"),
			expected: []rune("file.txt"),
		},
		{
			name:     "multiple branches",
			input:    []rune("├─└─ file.txt"),
			expected: []rune("file.txt"),
		},
		{
			name:     "no branch",
			input:    []rune("file.txt"),
			expected: []rune("file.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeBranchPrefixes(tt.input)
			if string(result) != string(tt.expected) {
				t.Errorf("removeBranchPrefixes(%q) = %q, want %q", string(tt.input), string(result), string(tt.expected))
			}
		})
	}
}

func TestRemoveResidualBoxDrawing(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		expected []rune
	}{
		{
			name:     "box drawing characters",
			input:    []rune("─│| file.txt"),
			expected: []rune("file.txt"),
		},
		{
			name:     "mixed characters with spaces",
			input:    []rune("─ │ | file.txt"),
			expected: []rune("file.txt"),
		},
		{
			name:     "no box drawing",
			input:    []rune("file.txt"),
			expected: []rune("file.txt"),
		},
		{
			name:     "preserve dots",
			input:    []rune(".gitignore"),
			expected: []rune(".gitignore"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeResidualBoxDrawing(tt.input)
			if string(result) != string(tt.expected) {
				t.Errorf("removeResidualBoxDrawing(%q) = %q, want %q", string(tt.input), string(result), string(tt.expected))
			}
		})
	}
}

// cmpEntries provides a detailed comparison between expected and actual entries for debugging
func cmpEntries(want, got []Entry) string {
	result := ""
	maxLen := len(want)
	if len(got) > maxLen {
		maxLen = len(got)
	}

	result += "Expected vs Got:\n"
	for i := 0; i < maxLen; i++ {
		wantStr := "                              "
		gotStr := "                              "

		if i < len(want) {
			kind := "Dir "
			if want[i].Kind == KindFile {
				kind = "File"
			}
			wantStr = want[i].Path + " (" + kind + ")"
		}

		if i < len(got) {
			kind := "Dir "
			if got[i].Kind == KindFile {
				kind = "File"
			}
			gotStr = got[i].Path + " (" + kind + ")"
		}

		marker := "  "
		if i < len(want) && i < len(got) && !reflect.DeepEqual(want[i], got[i]) {
			marker = "❌"
		} else if i >= len(want) {
			marker = "➕"
		} else if i >= len(got) {
			marker = "➖"
		}

		result += marker + " " + wantStr + " | " + gotStr + "\n"
	}

	return result
}
