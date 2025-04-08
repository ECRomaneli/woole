package parser

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// regexs
var (
	// DurationRegex matches a duration string with optional whitespace
	// and valid qualifiers (d, h, m, s, ms, ns).
	durationRegex = regexp.MustCompile(`^((\d+(\.\d{1,2})?)([dhms]|min|ms|ns)\s?)+$`)
	// BytesRegex matches a size string with optional whitespace
	// and valid qualifiers (b, kb, mb, gb, tb).
	bytesRegex = regexp.MustCompile(`^((\d+(\.\d{1,2})?)([kmgt]?b)\s?)+$`)
)

// ParseDuration parses a string into a time.Duration.
// Valid qualifiers: d, h, m, min, s, ms, ns.
func ParseDuration(input string) (time.Duration, error) {
	input = strings.ToLower(input)

	if input == "0" {
		return 0, nil
	}

	if input == "max" {
		return time.Duration(math.MaxInt64), nil
	}

	// Validate and split input
	components, err := validateAndSplit(input, durationRegex)
	if err != nil {
		return 0, err
	}

	var total time.Duration
	for _, component := range components {
		// Extract number and qualifier
		value, qualifier, err := extractNumberAndQualifier(component, durationRegex)
		if err != nil {
			return 0, err
		}

		num, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("invalid number in component: %s", value)
		}

		// Add the corresponding duration based on the qualifier
		switch strings.ToLower(qualifier) {
		case "d":
			total += time.Duration(num) * 24 * time.Hour
		case "h":
			total += time.Duration(num) * time.Hour
		case "m", "min":
			total += time.Duration(num) * time.Minute
		case "s":
			total += time.Duration(num) * time.Second
		case "ms":
			total += time.Duration(num) * time.Millisecond
		case "ns":
			total += time.Duration(num) * time.Nanosecond
		default:
			return 0, fmt.Errorf("invalid qualifier in component: %s", qualifier)
		}
	}

	return total, nil
}

// ParseBytes parses a string into a size in bytes (int64).
// Valid qualifiers: b, kb, mb, gb, tb.
func ParseBytes(input string) (int64, error) {
	input = strings.ToLower(input)

	if input == "0" {
		return 0, nil
	}

	if input == "maxint32" {
		return math.MaxInt32, nil
	}

	if input == "maxint64" || input == "max" {
		return math.MaxInt64, nil
	}

	// Validate and split input
	components, err := validateAndSplit(input, bytesRegex)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, component := range components {
		// Extract number and qualifier
		value, qualifier, err := extractNumberAndQualifier(component, bytesRegex)
		if err != nil {
			return 0, err
		}

		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid number in component: %s", value)
		}

		// Add the corresponding size based on the qualifier
		switch strings.ToLower(qualifier) {
		case "b":
			total += num
		case "kb":
			total += num * 1024
		case "mb":
			total += num * 1024 * 1024
		case "gb":
			total += num * 1024 * 1024 * 1024
		case "tb":
			total += num * 1024 * 1024 * 1024 * 1024
		default:
			return 0, fmt.Errorf("invalid qualifier in component: %s", qualifier)
		}

		if math.IsInf(total, 1) {
			return 0, fmt.Errorf("size exceeded float64 limit: %s", input)
		}
	}

	return int64(math.Round(total)), nil
}

// validateAndSplit validates the input string using the provided regex and splits it into components.
func validateAndSplit(input string, re *regexp.Regexp) ([]string, error) {
	if !re.MatchString(input) {
		return nil, fmt.Errorf("invalid format: %s", input)
	}
	return strings.Fields(input), nil
}

// extractNumberAndQualifier extracts the numeric value and qualifier from a component.
func extractNumberAndQualifier(component string, re *regexp.Regexp) (string, string, error) {
	matches := re.FindStringSubmatch(component)
	if len(matches) < 5 {
		return "", "", fmt.Errorf("invalid component: %s", component)
	}

	// Extract the numeric value and qualifier
	numStr := matches[2]
	qualifier := matches[4]

	return numStr, qualifier, nil
}
