package draw

import (
	"strings"
)

// Constants for box drawing characters and formatting
const (
	topLeftCorner     = '┌'
	topRightCorner    = '┐'
	bottomLeftCorner  = '└'
	bottomRightCorner = '┘'
	horizontalLine    = "─"
	verticalLine      = '│'
	space             = ' '
	spaceStr          = " "
	colonSpace        = ": "
	newline           = '\n'
)

// KeyValue represents a key-value pair for the box display
type KeyValue struct {
	Key   string
	Value string
}

// Returns a slice of KeyValue pairs as a formatted box in the console
// with right-aligned keys (padded on the left) and aligned borders
func Box(data []KeyValue) string {
	if len(data) == 0 {
		return ""
	}

	// Find the length of the longest key and value
	maxKeyLen := 0
	maxValueLen := 0
	for _, item := range data {
		if len(item.Key) > maxKeyLen {
			maxKeyLen = len(item.Key)
		}
		if len(item.Value) > maxValueLen {
			maxValueLen = len(item.Value)
		}
	}

	// Calculate the total width of the box content
	// Format: │ <padded-key>: <value><padding> │
	contentWidth := maxKeyLen + 4 + maxValueLen

	// Pre-calculate the horizontal border
	horizontalBorder := strings.Repeat(horizontalLine, contentWidth)

	// Calculate byte size precisely based on UTF-8 encoding
	// Each box drawing character (┌,┐,└,┘,│,─) takes 3 bytes in UTF-8
	// Normal ASCII characters (spaces, colons, letters) take 1 byte each
	// Newlines take 1 byte each

	// Box corners: 4 corners × 3 bytes each = 12 bytes
	// Horizontal borders: 2 lines × contentWidth × 3 bytes each = 6 × contentWidth bytes
	// Vertical borders: data lines × 2 sides × 3 bytes each = 6 × len(data) bytes
	// Content: sum of all key lengths + value lengths + padding + formatting
	// Newlines: (len(data) + 2) lines × 1 byte each = len(data) + 2 bytes

	// Conservative estimate for content and precise calculation for box elements
	totalSize := 14 + (6 * contentWidth) + (6 * len(data)) + (contentWidth * len(data)) + len(data)

	var builder strings.Builder
	builder.Grow(totalSize)

	// Draw the top border
	builder.WriteRune(topLeftCorner)
	builder.WriteString(horizontalBorder)
	builder.WriteRune(topRightCorner)
	builder.WriteByte(newline)

	// Build each key-value pair line
	for _, item := range data {
		// Right-align the keys by adding padding to the left
		keyPadding := strings.Repeat(spaceStr, maxKeyLen-len(item.Key))

		// Calculate padding needed after the value
		valuePadding := strings.Repeat(spaceStr, maxValueLen-len(item.Value))

		builder.WriteRune(verticalLine)
		builder.WriteByte(space)
		builder.WriteString(keyPadding)
		builder.WriteString(item.Key)
		builder.WriteString(colonSpace)
		builder.WriteString(item.Value)
		builder.WriteString(valuePadding)
		builder.WriteByte(space)
		builder.WriteRune(verticalLine)
		builder.WriteByte(newline)
	}

	// Draw the bottom border
	builder.WriteRune(bottomLeftCorner)
	builder.WriteString(horizontalBorder)
	builder.WriteRune(bottomRightCorner)
	builder.WriteByte(newline)

	return builder.String()
}
