package todos

import (
	"bytes"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// GetBlameInfo runs git blame for a single line and returns author/commit/date info.
// Returns (nil, nil) for non-git repos or on any error.
func GetBlameInfo(dir, filePath string, line int) (*BlameInfo, error) {
	lineSpec := strconv.Itoa(line) + "," + strconv.Itoa(line)
	cmd := exec.Command("git", "blame", "-L", lineSpec, "--porcelain", "--", filePath)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, nil //nolint:nilerr // gracefully handle non-git repos
	}

	return parsePorcelain(out.String()), nil
}

// parsePorcelain extracts blame info from git blame --porcelain output.
func parsePorcelain(output string) *BlameInfo {
	info := &BlameInfo{}
	lines := strings.Split(output, "\n")

	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) > 0 {
			info.Commit = parts[0]
		}
	}

	for _, line := range lines[1:] {
		switch {
		case strings.HasPrefix(line, "author "):
			info.Author = strings.TrimPrefix(line, "author ")
		case strings.HasPrefix(line, "author-time "):
			info.Date = strings.TrimPrefix(line, "author-time ")
		}
	}

	return info
}

// CalculateAge returns the number of days since the blame date (unix timestamp string).
func CalculateAge(blameDate string) int {
	ts, err := strconv.ParseInt(blameDate, 10, 64)
	if err != nil {
		return 0
	}
	t := time.Unix(ts, 0)
	days := time.Since(t).Hours() / 24
	return int(math.Floor(days))
}
