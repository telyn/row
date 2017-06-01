package row

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// methodByNameFunc is the equivalent of reflect.Value.FieldByNameFunc, but for methods
func methodByNameFunc(v reflect.Value, match func(string) bool) reflect.Value {
	valueType := v.Type()
	for i := 0; i < v.NumMethod(); i++ {
		m := valueType.Method(i)
		if match(m.Name) {
			return m.Func
		}

	}
	return reflect.Value{}
}

// insensitiveMatch returns a match function for case-insensitive matching to the given string
func insensitiveMatch(expected string) func(string) bool {
	return func(name string) bool {
		return strings.EqualFold(expected, name)
	}
}

// SortFields sorts the given fields on obj according to an index tag on those fields
// If no tagName is given, uses 'index'.
// see ExampleSortFields
func SortFields(obj interface{}, fieldNames []string, tagName ...string) (newFieldNames []string) {
	reflectType := reflect.TypeOf(obj)
	// remove nonexistent fields
	fields := make([]reflect.StructField, 0, len(fieldNames))
	for _, f := range fieldNames {
		if field, ok := reflectType.FieldByName(f); ok {
			fields = append(fields, field)
		}
	}

	tag := "index"
	if len(tagName) > 0 {
		tag = tagName[0]
	}
	// sort by index
	compare := mkFieldSortComparator(tag)
	sort.SliceStable(fields, func(i, j int) bool {
		return compare(fields[i], fields[j])
	})
	newFieldNames = make([]string, len(fields))
	for i, f := range fields {
		newFieldNames[i] = f.Name
	}
	return

}

func mkFieldSortComparator(tagName string) func(a, b reflect.StructField) bool {
	return func(a reflect.StructField, b reflect.StructField) bool {
		// return a < b
		idxastr, oka := a.Tag.Lookup(tagName)
		idxbstr, okb := b.Tag.Lookup(tagName)
		if oka {
			if !okb {
				// a < b when b has no tag but a does
				return true
			}
			idxa, erra := strconv.Atoi(idxastr)
			idxb, errb := strconv.Atoi(idxbstr)
			if erra == nil {
				if errb == nil {
					// compare tags when both tags are ints
					return idxa < idxb
				}
				// a < b when a's tag is int, but b's isn't
				return true
			}
			if errb == nil {
				// a > b when b's tag is int, but a's isnt
				return false
			}
			// compare field names if neither tag is an int
			return a.Name < b.Name
		}
		if okb {
			// a > b when b has a tag but a doesn't, even if tag is invalid
			return false
		}
		// compare field names if neither tag has index.
		return a.Name < b.Name

	}
}
