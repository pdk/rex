package rex

import (
	"log"
	"strconv"
)

// Record holds a single record of values
type Record struct {
	values  map[string]interface{}
	dropped bool
}

// NewRecord constructs a new Record from a JSON parse result.
func NewRecord(value interface{}) Record {

	switch obj := value.(type) {
	case nil:
		return Record{
			dropped: true,
		}
	case map[string]interface{}:
		return Record{
			values: obj,
		}
	}

	return Record{
		values: map[string]interface{}{
			"value": value,
		},
	}
}

// Has returns true IFF the record has the key, and the value is non-null/non-empty.
func (r Record) Has(key string) bool {
	if r.dropped {
		log.Fatalf("Int called on dropped record")
	}

	val, ok := r.values[key]
	if !ok {
		return false
	}

	switch val := val.(type) {
	case nil:
		return false
	case []interface{}:
		return len(val) > 0
	case map[string]interface{}:
		return len(val) > 0
	}

	return true
}

// Drop returns a dropped Record if true.
func (r Record) Drop(f bool) Record {
	if f {
		return Record{dropped: true}
	}

	return r
}

// Keep will return the record if true, otherwise a dropped record.
func (r Record) Keep(f bool) Record {
	if f {
		return r
	}

	return Record{dropped: true}
}

// Set adds/overwrites a value for a key to the record.
func (r Record) Set(key string, val interface{}) Record {
	r.values[key] = val
	return r
}

// Value returns any value in the Record.
func (r Record) Value(key string) interface{} {
	val, ok := r.values[key]
	if !ok {
		return nil
	}

	return val
}

// Bool returns true or false.  Returns true for a string value of "true". For
// numbers, 0 is false, otherwise true. Returns false for missing/unconvertable
// values.
func (r Record) Bool(key string) bool {
	if r.dropped {
		log.Fatalf("Int called on dropped record")
	}

	val, ok := r.values[key]
	if !ok {
		return false
	}

	switch val := val.(type) {
	case nil:
		return false
	case string:
		return val == "true"
	case int64:
		return val != 0
	case float64:
		return val != 0.0
	case bool:
		return val
	}

	return false
}

// Int returns an int64 of the key. Returns 0 for missing/unconvertable values.
func (r Record) Int(key string) int64 {
	if r.dropped {
		log.Fatalf("Int called on dropped record")
	}

	val, ok := r.values[key]
	if !ok {
		return 0
	}

	switch val := val.(type) {
	case nil:
		return 0
	case string:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0
		}
		return i
	case int64:
		return val
	}

	return 0
}

type stringer interface {
	String() string
}

// String returns a string value. Will return "" for missing/null/unconvertable values.
func (r Record) String(key string) string {
	if r.dropped {
		log.Fatalf("String called on dropped record")
	}

	val, ok := r.values[key]
	if !ok {
		return ""
	}

	switch val := val.(type) {
	case nil:
		return ""
	case string:
		return val
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		if val {
			return "true"
		}
		return "false"
	}

	stringerVal, ok := val.(stringer)
	if ok {
		return stringerVal.String()
	}

	return ""
}
