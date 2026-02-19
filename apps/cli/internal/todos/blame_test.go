package todos

import (
	"fmt"
	"testing"
	"time"
)

func TestParsePorcelain_BasicOutput(t *testing.T) {
	porcelain := `abc1234567890123456789012345678901234567 3 3 1
author John Doe
author-mail <john@example.com>
author-time 1700000000
author-tz +0000
committer John Doe
committer-mail <john@example.com>
committer-time 1700000000
committer-tz +0000
summary Fix the thing
filename main.go
	// TODO: fix this
`
	info := parsePorcelain(porcelain)
	if info.Author != "John Doe" {
		t.Errorf("expected author 'John Doe', got %q", info.Author)
	}
	if info.Commit != "abc1234567890123456789012345678901234567" {
		t.Errorf("expected commit hash, got %q", info.Commit)
	}
	if info.Date != "1700000000" {
		t.Errorf("expected date '1700000000', got %q", info.Date)
	}
}

func TestParsePorcelain_EmptyOutput(t *testing.T) {
	info := parsePorcelain("")
	if info.Author != "" || info.Commit != "" || info.Date != "" {
		t.Errorf("expected empty info for empty output, got %+v", info)
	}
}

func TestCalculateAge_RecentDate(t *testing.T) {
	ts := time.Now().Add(-10 * 24 * time.Hour).Unix()
	tsStr := fmt.Sprintf("%d", ts)
	age := CalculateAge(tsStr)
	if age < 9 || age > 11 {
		t.Errorf("expected age ~10 days, got %d", age)
	}
}

func TestCalculateAge_ValidTimestamp(t *testing.T) {
	age := CalculateAge("1700000000")
	if age <= 0 {
		t.Errorf("expected positive age for past timestamp, got %d", age)
	}
}

func TestCalculateAge_InvalidTimestamp(t *testing.T) {
	age := CalculateAge("not-a-number")
	if age != 0 {
		t.Errorf("expected 0 for invalid timestamp, got %d", age)
	}
}

func TestCalculateAge_EmptyString(t *testing.T) {
	age := CalculateAge("")
	if age != 0 {
		t.Errorf("expected 0 for empty string, got %d", age)
	}
}

func TestGetBlameInfo_NonGitDir(t *testing.T) {
	dir := t.TempDir()
	info, err := GetBlameInfo(dir, "nonexistent.go", 1)
	if err != nil {
		t.Errorf("expected nil error for non-git dir, got %v", err)
	}
	if info != nil {
		t.Errorf("expected nil info for non-git dir, got %+v", info)
	}
}
