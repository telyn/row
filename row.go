package row

import (
	"reflect"
)

// From converts the given object into a row for a olekukonko/tablewriter table, using reflection. fields should be an array of (exported) fields on obj which are strings, ints, bools, fmt.Stringers, or slices thereof, or are methods taking no arguments which return a single value of the same, or a single value of the same and an error.
func From(obj interface{}, fields []string) (row []string, err error) {
	row = make([]string, len(fields))

	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	for i, field := range fields {
		// TODO(telyn): check field starts with capital to ensure we only use exported types
		v := value.FieldByName(field)
		str := ""
		if v.Kind() == reflect.Invalid {
			v := value.MethodByName(field)
			if v.Kind() == reflect.Invalid {
				str = "no field called " + field
			} else {
				str, err = methodToString(v)
			}
		} else {
			str, err = valueToString(v)
		}

		if err != nil {
			return nil, err
		}
		row[i] = str

	}
	return

}

// FieldsFrom lists all the fields available for the given object, and can be fed to From too.
func FieldsFrom(obj interface{}) (fields []string) {
	value := reflect.ValueOf(obj)
	t := value.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		return fieldsFromType(t.Elem())
	case reflect.Struct:
		return fieldsFromType(t)
	}
	return []string{}
}

func SortedFieldsFrom(obj interface{}, tagName ...string) (fields []string) {
	return SortFields(obj, FieldsFrom(obj), tagName...)
}
