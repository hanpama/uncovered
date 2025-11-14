package coverage

import (
	"sort"
)

// UncoveredLine represents a single uncovered line in a source file.
// It contains the line number and column where the uncovered code begins.
type UncoveredLine struct {
	Line int // Line number (1-indexed)
	Col  int // Column number (1-indexed)
}

// FileUncovered represents all uncovered lines in a single source file.
// Lines are sorted by line number in ascending order.
type FileUncovered struct {
	FileName string          // Relative path to the source file
	Lines    []UncoveredLine // Sorted list of uncovered lines
}

// GetUncoveredLines extracts all uncovered lines from a coverage profile.
//
// It processes all blocks in the profile and identifies those with zero
// execution count. For each uncovered block, all lines within that block
// are marked as uncovered.
//
// The returned slice contains one FileUncovered for each file that has
// at least one uncovered line. Both the files and the lines within each
// file are sorted for consistent output.
func GetUncoveredLines(profile *Profile) []*FileUncovered {
	fileMap := make(map[string]map[int]int) // filename -> line -> column

	for _, block := range profile.Blocks {
		if block.IsCovered() {
			continue
		}

		if _, exists := fileMap[block.FileName]; !exists {
			fileMap[block.FileName] = make(map[int]int)
		}

		// Add all lines in the uncovered block
		for line := block.StartLine; line <= block.EndLine; line++ {
			col := 0
			if line == block.StartLine {
				col = block.StartCol
			}
			// Keep the earliest column for each line
			if existingCol, exists := fileMap[block.FileName][line]; !exists || col < existingCol {
				fileMap[block.FileName][line] = col
			}
		}
	}

	// Convert map to sorted slice
	result := make([]*FileUncovered, 0, len(fileMap))
	for fileName, lineMap := range fileMap {
		lines := make([]UncoveredLine, 0, len(lineMap))
		for line, col := range lineMap {
			lines = append(lines, UncoveredLine{Line: line, Col: col})
		}

		// Sort lines by line number
		sort.Slice(lines, func(i, j int) bool {
			return lines[i].Line < lines[j].Line
		})

		result = append(result, &FileUncovered{
			FileName: fileName,
			Lines:    lines,
		})
	}

	// Sort files by filename
	sort.Slice(result, func(i, j int) bool {
		return result[i].FileName < result[j].FileName
	})

	return result
}
