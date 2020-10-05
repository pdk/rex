package rex

import (
	"strconv"
	"time"
)

// Helpers to make best effort conversion to most correct type.

// GuessType takes any value, does some guesswork, and returns maybe the same
// value, maybe a "corrected" type value.
func GuessType(val interface{}) interface{} {

	switch val := val.(type) {
	case string:
		return GuessStringType(val)
	// all ints -> int64
	case int:
		return int64(val)
	case int8:
		return int64(val)
	case int16:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	// all uints -> uint64
	case uint:
		return uint64(val)
	case uint8:
		return uint64(val)
	case uint16:
		return uint64(val)
	case uint32:
		return uint64(val)
	case uint64:
		return val
	// all floats -> float64
	case float32:
		return float64(val)
	case float64:
		return val
	case map[string]interface{}:
		newMap := map[string]interface{}{}
		for k, v := range val {
			newMap[k] = GuessType(v)
		}
		return newMap
	}

	return val
}

// GuessStringType tries various ways to parse a string.
func GuessStringType(s string) interface{} {

	switch s {
	case "":
		return ""
	case "null", "nil":
		return nil
	case "true", "True":
		return true
	case "false", "False":
		return false
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return i
	}

	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}

	for _, f := range timeFormats {
		t, err := time.Parse(f, s)
		if err == nil {
			return t
		}
	}

	return s
}

var (
	timeFormats = []string{
		time.ANSIC,    // "Mon Jan _2 15:04:05 2006"
		time.UnixDate, // "Mon Jan _2 15:04:05 MST 2006"
		time.RubyDate, // "Mon Jan 02 15:04:05 -0700 2006"
		time.RFC822,   // "02 Jan 06 15:04 MST"
		time.RFC822Z,  // "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
		time.RFC850,   // "Monday, 02-Jan-06 15:04:05 MST"
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
		time.RFC3339,  // "2006-01-02T15:04:05Z07:00"
	}
)
