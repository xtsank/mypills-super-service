//go:build ignore

// This command counts covered source lines from a Go cover profile.
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var block = regexp.MustCompile(`^(.+):(\d+)\.\d+,(\d+)\.(\d+) \d+ (\d+)$`)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run scripts/line-coverage.go coverage.out")
		os.Exit(2)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	lines := make(map[string]map[int]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := scanner.Text()
		if strings.HasPrefix(entry, "mode: ") {
			continue
		}
		match := block.FindStringSubmatch(entry)
		if match == nil {
			fatal(fmt.Errorf("invalid cover profile entry: %q", entry))
		}
		start, _ := strconv.Atoi(match[2])
		end, _ := strconv.Atoi(match[3])
		endColumn, _ := strconv.Atoi(match[4])
		count, _ := strconv.Atoi(match[5])
		if endColumn == 1 && end > start {
			end--
		}
		if lines[match[1]] == nil {
			lines[match[1]] = make(map[int]bool)
		}
		for line := start; line <= end; line++ {
			lines[match[1]][line] = lines[match[1]][line] || count > 0
		}
	}
	if err := scanner.Err(); err != nil {
		fatal(err)
	}

	total, covered := 0, 0
	for _, fileLines := range lines {
		for _, hit := range fileLines {
			total++
			if hit {
				covered++
			}
		}
	}
	if total == 0 {
		fatal(fmt.Errorf("cover profile has no instrumented lines"))
	}
	fmt.Printf("instrumented line coverage: %d/%d (%.1f%%)\n", covered, total, float64(covered)*100/float64(total))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
