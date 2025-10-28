package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessInput(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "myapp/\n├─ src/\n│  └─ main.go\n└─ README.md"
	
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		inputFile   string
		verbose     bool
		expectError bool
		expectedLen int
	}{
		{
			name:        "read from file",
			inputFile:   testFile,
			verbose:     false,
			expectError: false,
			expectedLen: 4,
		},
		{
			name:        "nonexistent file",
			inputFile:   "/nonexistent/file.txt",
			verbose:     false,
			expectError: true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, err := processInput(tt.inputFile, tt.verbose)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("processInput() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("processInput() unexpected error: %v", err)
				return
			}
			
			if len(lines) != tt.expectedLen {
				t.Errorf("processInput() got %d lines, want %d", len(lines), tt.expectedLen)
			}
		})
	}
}

func TestDetermineRootName(t *testing.T) {
	tests := []struct {
		name     string
		rootName string
		firstLine string
		expected string
	}{
		{
			name:     "explicit root name",
			rootName: "custom-root",
			firstLine: "myapp/",
			expected: "custom-root",
		},
		{
			name:     "extract from first line",
			rootName: "",
			firstLine: "myapp/",
			expected: "myapp",
		},
		{
			name:     "extract from first line without slash",
			rootName: "",
			firstLine: "myapp",
			expected: "myapp",
		},
		{
			name:     "empty first line falls back to output",
			rootName: "",
			firstLine: "",
			expected: "output",
		},
		{
			name:     "whitespace only first line falls back to output",
			rootName: "",
			firstLine: "   ",
			expected: "output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineRootName(tt.rootName, tt.firstLine)
			if result != tt.expected {
				t.Errorf("determineRootName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestPrintDryRun(t *testing.T) {
	entries := []Entry{
		{Path: "src", Kind: KindDir},
		{Path: "src/main.go", Kind: KindFile},
		{Path: "README.md", Kind: KindFile},
	}
	
	basePath := "/tmp/test"
	
	// This function would normally print to stdout
	// We're testing that it doesn't panic and handles the entries correctly
	printDryRun(basePath, entries)
	
	// No assertion needed - just checking it doesn't panic
}

func TestCreateEntry(t *testing.T) {
	tmpDir := t.TempDir()
	
	tests := []struct {
		name        string
		entry       Entry
		basePath    string
		force       bool
		verbose     bool
		expectError bool
	}{
		{
			name:        "create directory",
			entry:       Entry{Path: "testdir", Kind: KindDir},
			basePath:    tmpDir,
			force:       false,
			verbose:     false,
			expectError: false,
		},
		{
			name:        "create file",
			entry:       Entry{Path: "testfile.txt", Kind: KindFile},
			basePath:    tmpDir,
			force:       false,
			verbose:     false,
			expectError: false,
		},
		{
			name:        "create nested file",
			entry:       Entry{Path: "nested/deep/file.txt", Kind: KindFile},
			basePath:    tmpDir,
			force:       false,
			verbose:     false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createEntry(tt.entry, tt.basePath, tt.force, tt.verbose)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("createEntry() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("createEntry() unexpected error: %v", err)
				return
			}
			
			// Check that the result indicates success
			if result != "created" && result != "skipped" {
				t.Errorf("createEntry() got unexpected result: %s", result)
			}
			
			// Verify the file/directory was actually created
			fullPath := filepath.Join(tt.basePath, tt.entry.Path)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("createEntry() did not create %s", fullPath)
			}
		})
	}
}

func TestCreateEntryExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.txt")
	
	// Create an existing file
	if err := os.WriteFile(existingFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}
	
	entry := Entry{Path: "existing.txt", Kind: KindFile}
	
	// Test without force - should skip
	result, err := createEntry(entry, tmpDir, false, false)
	if err != nil {
		t.Errorf("createEntry() unexpected error: %v", err)
	}
	if result != "skipped" {
		t.Errorf("createEntry() expected 'skipped', got %q", result)
	}
	
	// Test with force - should overwrite
	result, err = createEntry(entry, tmpDir, true, false)
	if err != nil {
		t.Errorf("createEntry() unexpected error: %v", err)
	}
	if result != "created" {
		t.Errorf("createEntry() expected 'created', got %q", result)
	}
}