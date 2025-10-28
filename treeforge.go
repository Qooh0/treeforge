package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	version = "1.0.0"
)

// Helper functions for main
func processInput(inputFile string, verbose bool) ([]string, error) {
	if inputFile != "" {
		return readFromFile(inputFile, verbose)
	}
	return readFromStdin(verbose)
}

func determineRootName(rootName, firstLine string) string {
	if rootName != "" {
		return rootName
	}
	
	root := strings.TrimSpace(firstLine)
	root = strings.TrimSuffix(root, "/")
	if root == "" {
		return "output"
	}
	return root
}

func printDryRun(basePath string, entries []Entry) {
	fmt.Println("=== Dry-run mode (use --apply to create files) ===")
	fmt.Printf("Base: %s\n\n", basePath)
	for _, entry := range entries {
		fullPath := filepath.Join(basePath, entry.Path)
		if entry.Kind == KindDir {
			fmt.Printf("  [DIR]  %s\n", fullPath)
		} else {
			fmt.Printf("  [FILE] %s\n", fullPath)
		}
	}
	fmt.Printf("\nTotal: %d directories, %d files\n", countDirs(entries), countFiles(entries))
}

func createEntry(entry Entry, basePath string, force, verbose bool) (string, error) {
	fullPath := filepath.Join(basePath, entry.Path)

	if entry.Kind == KindDir {
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return "", fmt.Errorf("creating directory %s: %w", fullPath, err)
		}
		if verbose {
			fmt.Printf("  [DIR]  %s\n", fullPath)
		}
		return "created", nil
	} else {
		// Check if file exists
		if _, err := os.Stat(fullPath); err == nil && !force {
			if verbose {
				fmt.Printf("  [SKIP] %s (already exists)\n", fullPath)
			}
			return "skipped", nil
		}

		// Create parent directory if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("creating directory for file %s: %w", fullPath, err)
		}

		// Create empty file
		file, err := os.Create(fullPath)
		if err != nil {
			return "", fmt.Errorf("creating file %s: %w", fullPath, err)
		}
		file.Close()

		if verbose {
			fmt.Printf("  [FILE] %s\n", fullPath)
		}
		return "created", nil
	}
}

func applyEntries(basePath string, entries []Entry, force, verbose bool) error {
	// Create base directory
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("creating base directory: %w", err)
	}

	created := 0
	skipped := 0

	for _, entry := range entries {
		result, err := createEntry(entry, basePath, force, verbose)
		if err != nil {
			return err
		}
		
		if result == "created" {
			created++
		} else if result == "skipped" {
			skipped++
		}
	}

	fmt.Printf("\n✓ Done! Created: %d, Skipped: %d\n", created, skipped)
	return nil
}

func main() {
	var (
		inputFile = flag.String("i", "", "Input tree structure file (default: stdin)")
		parent    = flag.String("parent", ".", "Parent directory to create structure in")
		rootName  = flag.String("root-name", "", "Override root directory name (from first line if empty)")
		apply     = flag.Bool("apply", false, "Actually create files/directories (default: dry-run)")
		force     = flag.Bool("force", false, "Overwrite existing files (directories are not deleted)")
		verbose   = flag.Bool("v", false, "Verbose output")
		showVer   = flag.Bool("version", false, "Show version")
	)
	flag.Parse()

	if *showVer {
		fmt.Printf("treeforge v%s\n", version)
		return
	}

	// Read input
	lines, err := processInput(*inputFile, *verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	if len(lines) == 0 {
		fmt.Fprintf(os.Stderr, "Error: empty input\n")
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Read %d lines\n", len(lines))
	}

	// Parse tree structure
	entries, err := ParseTree(lines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing tree: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Parsed %d entries\n", len(entries))
	}

	// Determine root name and create base path
	root := determineRootName(*rootName, lines[0])
	basePath := filepath.Join(*parent, root)

	// Dry-run or apply
	if !*apply {
		printDryRun(basePath, entries)
		return
	}

	// Apply mode: create files and directories
	if *verbose {
		fmt.Printf("Creating structure in: %s\n", basePath)
	}

	if err := applyEntries(basePath, entries, *force, *verbose); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func readFromFile(path string, verbose bool) ([]string, error) {
	if verbose {
		fmt.Printf("Reading from file: %s\n", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func readFromStdin(verbose bool) ([]string, error) {
	if verbose {
		fmt.Println("Reading from stdin...")
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// stdin is from terminal (not piped)
		fmt.Fprintln(os.Stderr, "Paste your tree structure (press Ctrl+D when done):")
	}

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func countDirs(entries []Entry) int {
	count := 0
	for _, e := range entries {
		if e.Kind == KindDir {
			count++
		}
	}
	return count
}

func countFiles(entries []Entry) int {
	count := 0
	for _, e := range entries {
		if e.Kind == KindFile {
			count++
		}
	}
	return count
}
