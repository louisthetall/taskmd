package nextid

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// Result holds the computed next ID and related metadata.
type Result struct {
	NextID  string `json:"next_id" yaml:"next_id"`
	MaxID   string `json:"max_id" yaml:"max_id"`
	Prefix  string `json:"prefix" yaml:"prefix"`
	Padding int    `json:"padding" yaml:"padding"`
	Total   int    `json:"total" yaml:"total"`
}

type parsedID struct {
	original string
	prefix   string
	number   int
	numStr   string
}

// Calculate determines the next available ID from a list of existing IDs.
// It finds the maximum numeric suffix and returns max+1, preserving any
// common prefix and zero-padding.
func Calculate(ids []string) Result {
	var parsed []parsedID
	for _, id := range ids {
		if p, ok := parseID(id); ok {
			parsed = append(parsed, p)
		}
	}

	if len(parsed) == 0 {
		return Result{
			NextID:  "001",
			Padding: 3,
			Total:   len(ids),
		}
	}

	maxNum := 0
	maxParsed := parsed[0]
	for _, p := range parsed {
		if p.number > maxNum {
			maxNum = p.number
			maxParsed = p
		}
	}

	prefix := detectPrefix(parsed)
	padding := max(len(maxParsed.numStr), 3)

	nextNum := maxNum + 1
	nextID := formatID(prefix, nextNum, padding)

	return Result{
		NextID:  nextID,
		MaxID:   maxParsed.original,
		Prefix:  prefix,
		Padding: padding,
		Total:   len(ids),
	}
}

// parseID extracts the trailing numeric portion and any prefix from an ID.
// Returns false if the ID contains no digits at the end.
func parseID(id string) (parsedID, bool) {
	if id == "" {
		return parsedID{}, false
	}

	// Scan backward to find where trailing digits start
	i := len(id) - 1
	for i >= 0 && id[i] >= '0' && id[i] <= '9' {
		i--
	}

	numStr := id[i+1:]
	if numStr == "" {
		return parsedID{}, false
	}

	num := 0
	for _, ch := range numStr {
		num = num*10 + int(ch-'0')
	}

	return parsedID{
		original: id,
		prefix:   id[:i+1],
		number:   num,
		numStr:   numStr,
	}, true
}

// detectPrefix returns the most common prefix if it appears in more than
// 50% of the parsed IDs. Otherwise returns "".
func detectPrefix(parsed []parsedID) string {
	if len(parsed) == 0 {
		return ""
	}

	counts := make(map[string]int)
	for _, p := range parsed {
		counts[p.prefix]++
	}

	bestPrefix := ""
	bestCount := 0
	for prefix, count := range counts {
		if count > bestCount {
			bestCount = count
			bestPrefix = prefix
		}
	}

	if bestCount*2 > len(parsed) {
		return bestPrefix
	}
	return ""
}

// GeneratePrefixed produces the next sequential ID with the given prefix.
// It filters existing IDs by prefix, finds the max numeric suffix, and
// returns prefix + zero-padded(max+1).
func GeneratePrefixed(existingIDs []string, prefix string, padding int) string {
	maxNum := 0
	for _, id := range existingIDs {
		p, ok := parseID(id)
		if !ok || !strings.EqualFold(p.prefix, prefix) {
			continue
		}
		if p.number > maxNum {
			maxNum = p.number
		}
	}
	return formatID(prefix, maxNum+1, padding)
}

// GenerateRandom produces a random base-36 alphanumeric lowercase ID of the
// given length. It retries on collision with existingIDs (max 100 attempts).
func GenerateRandom(existingIDs []string, length int) (string, error) {
	existing := make(map[string]struct{}, len(existingIDs))
	for _, id := range existingIDs {
		existing[id] = struct{}{}
	}

	const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
	charsetLen := big.NewInt(int64(len(charset)))

	for attempt := 0; attempt < 100; attempt++ {
		buf := make([]byte, length)
		for i := range buf {
			idx, err := rand.Int(rand.Reader, charsetLen)
			if err != nil {
				return "", fmt.Errorf("crypto/rand failed: %w", err)
			}
			buf[i] = charset[idx.Int64()]
		}
		id := string(buf)
		if _, taken := existing[id]; !taken {
			return id, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique ID after 100 attempts")
}

// GenerateUUID produces a random hex ID (from UUID v4 space) of the given length.
// Default length is 8 when length <= 0. It retries on collision with existingIDs (max 100 attempts).
func GenerateUUID(existingIDs []string, length int) (string, error) {
	if length <= 0 {
		length = 8
	}

	existing := make(map[string]struct{}, len(existingIDs))
	for _, id := range existingIDs {
		existing[id] = struct{}{}
	}

	const charset = "0123456789abcdef"
	charsetLen := big.NewInt(int64(len(charset)))

	for attempt := 0; attempt < 100; attempt++ {
		buf := make([]byte, length)
		for i := range buf {
			idx, err := rand.Int(rand.Reader, charsetLen)
			if err != nil {
				return "", fmt.Errorf("crypto/rand failed: %w", err)
			}
			buf[i] = charset[idx.Int64()]
		}
		id := string(buf)
		if _, taken := existing[id]; !taken {
			return id, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique UUID ID after 100 attempts")
}

// formatID assembles a prefix with a zero-padded number.
func formatID(prefix string, number int, padding int) string {
	return fmt.Sprintf("%s%0*d", prefix, padding, number)
}
