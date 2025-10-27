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
	var lines []string
	var err error

	if *inputFile != "" {
		lines, err = readFromFile(*inputFile, *verbose)
	} else {
		lines, err = readFromStdin(*verbose)
	}

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

	// Determine root name
	root := *rootName
	if root == "" {
		root = strings.TrimSpace(lines[0])
		root = strings.TrimSuffix(root, "/")
		if root == "" {
			root = "output"
		}
	}

	// Create base path
	basePath := filepath.Join(*parent, root)

	// Dry-run or apply
	if !*apply {
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
		return
	}

	// Apply mode: create files and directories
	if *verbose {
		fmt.Printf("Creating structure in: %s\n", basePath)
	}

	// Create base directory
	if err := os.MkdirAll(basePath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating base directory: %v\n", err)
		os.Exit(1)
	}

	created := 0
	skipped := 0

	for _, entry := range entries {
		fullPath := filepath.Join(basePath, entry.Path)

		if entry.Kind == KindDir {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", fullPath, err)
				os.Exit(1)
			}
			if *verbose {
				fmt.Printf("  [DIR]  %s\n", fullPath)
			}
			created++
		} else {
			// Check if file exists
			if _, err := os.Stat(fullPath); err == nil && !*force {
				if *verbose {
					fmt.Printf("  [SKIP] %s (already exists)\n", fullPath)
				}
				skipped++
				continue
			}

			// Create parent directory if needed
			dir := filepath.Dir(fullPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory for file %s: %v\n", fullPath, err)
				os.Exit(1)
			}

			// Create empty file
			file, err := os.Create(fullPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", fullPath, err)
				os.Exit(1)
			}
			file.Close()

			if *verbose {
				fmt.Printf("  [FILE] %s\n", fullPath)
			}
			created++
		}
	}

	fmt.Printf("\n✓ Done! Created: %d, Skipped: %d\n", created, skipped)
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
