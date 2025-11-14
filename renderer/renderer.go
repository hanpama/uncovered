// Package renderer provides formatted output of uncovered code lines with context.
//
// This package takes the list of uncovered lines identified by the coverage
// package and renders them in a human-readable format with:
//   - File headers showing statistics (e.g., "29 of 59 lines uncovered")
//   - Context lines (3 before and after) around each uncovered section
//   - Color-coded output with red highlighting for uncovered lines
//   - Line numbers for easy navigation to the source
//
// # Output Format
//
// Each file is rendered with a header followed by groups of uncovered lines:
//
//	example/calculator.go (29 of 59 lines uncovered)
//	================================================
//	     21 }
//	     22
//	     23 // Divide returns the quotient of two numbers
//	>    24 func (c *Calculator) Divide(a, b int) (int, error) {
//	>    25   if b == 0 {
//	     ...
//
// Lines marked with ">" are uncovered (shown in red in terminal).
// Consecutive uncovered sections are intelligently grouped to avoid
// showing overlapping context.
package renderer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hanpama/uncovered/coverage"
)

const (
	contextLines = 3 // Number of context lines to show before and after
)

// Renderer renders uncovered lines with context.
//
// It reads the source files to display actual code content and
// uses ANSI color codes for terminal output.
type Renderer struct{}

// New creates a new renderer
func New() *Renderer {
	return &Renderer{}
}

// Render renders uncovered lines with context to the writer.
//
// For each file with uncovered lines, it:
//  1. Reads the source file from disk
//  2. Prints a header with file statistics
//  3. Groups nearby uncovered lines together
//  4. Displays each group with surrounding context
//
// The output uses ANSI color codes, so it's best suited for terminal output.
// Returns an error if any source file cannot be read.
func (r *Renderer) Render(w io.Writer, uncovered []*coverage.FileUncovered) error {
	for i, fileUncovered := range uncovered {
		if i > 0 {
			fmt.Fprintln(w)
		}

		if err := r.renderFile(w, fileUncovered); err != nil {
			return fmt.Errorf("rendering %s: %w", fileUncovered.FileName, err)
		}
	}
	return nil
}

// renderFile renders uncovered lines for a single file
func (r *Renderer) renderFile(w io.Writer, fileUncovered *coverage.FileUncovered) error {
	// Read file content
	file, err := os.Open(fileUncovered.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read all lines
	lines, err := r.readLines(file)
	if err != nil {
		return err
	}

	// Print file header with statistics
	r.printFileHeader(w, fileUncovered.FileName, len(fileUncovered.Lines), len(lines))

	// Group consecutive lines to avoid duplicate context
	groups := r.groupLines(fileUncovered.Lines)

	// Render each group
	for i, group := range groups {
		if i > 0 {
			fmt.Fprintln(w)
		}
		r.renderGroup(w, lines, group)
	}

	return nil
}

// lineGroup represents a group of consecutive uncovered lines
type lineGroup struct {
	start int // first uncovered line number
	end   int // last uncovered line number
	lines []coverage.UncoveredLine
}

// groupLines groups consecutive uncovered lines
func (r *Renderer) groupLines(lines []coverage.UncoveredLine) []lineGroup {
	if len(lines) == 0 {
		return nil
	}

	groups := make([]lineGroup, 0)
	currentGroup := lineGroup{
		start: lines[0].Line,
		end:   lines[0].Line,
		lines: []coverage.UncoveredLine{lines[0]},
	}

	for i := 1; i < len(lines); i++ {
		// If lines are close enough (within 2*contextLines), merge them into the same group
		if lines[i].Line <= currentGroup.end+2*contextLines+1 {
			currentGroup.end = lines[i].Line
			currentGroup.lines = append(currentGroup.lines, lines[i])
		} else {
			groups = append(groups, currentGroup)
			currentGroup = lineGroup{
				start: lines[i].Line,
				end:   lines[i].Line,
				lines: []coverage.UncoveredLine{lines[i]},
			}
		}
	}
	groups = append(groups, currentGroup)

	return groups
}

// renderGroup renders a group of uncovered lines with context
func (r *Renderer) renderGroup(w io.Writer, allLines []string, group lineGroup) {
	// Create a set of uncovered line numbers for quick lookup
	uncoveredSet := make(map[int]bool)
	for _, line := range group.lines {
		uncoveredSet[line.Line] = true
	}

	// Calculate context range
	startLine := max(1, group.start-contextLines)
	endLine := min(len(allLines), group.end+contextLines)

	// Render lines with context
	for lineNum := startLine; lineNum <= endLine; lineNum++ {
		isUncovered := uncoveredSet[lineNum]
		r.renderLine(w, lineNum, allLines[lineNum-1], isUncovered)
	}
}

// renderLine renders a single line
func (r *Renderer) renderLine(w io.Writer, lineNum int, content string, isUncovered bool) {
	prefix := "  "
	lineNumStr := fmt.Sprintf("%5d", lineNum)

	if isUncovered {
		prefix = "\033[31m>\033[0m " // Red arrow
		lineNumStr = fmt.Sprintf("\033[31m%5d\033[0m", lineNum) // Red line number
	}

	fmt.Fprintf(w, "%s %s %s\n", prefix, lineNumStr, content)
}

// printFileHeader prints the file header with statistics
func (r *Renderer) printFileHeader(w io.Writer, fileName string, uncoveredCount, totalLines int) {
	header := fmt.Sprintf("%s (%d of %d lines uncovered)", fileName, uncoveredCount, totalLines)
	fmt.Fprintf(w, "\033[1m%s\033[0m\n", header) // Bold
	fmt.Fprintln(w, strings.Repeat("=", min(80, len(header))))
}

// readLines reads all lines from a file
func (r *Renderer) readLines(file *os.File) ([]string, error) {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
