package prettyprint

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// ToTableRow converts the given object into a row for a olekukonko/tablewriter table, using reflection. fields should be an array of (exported) fields on obj which are strings, ints, bools, fmt.Stringers, or slices thereof, or are methods taking no arguments which return a single value of the same, or a single value of the same and an error.
func ToTableRow(obj interface{}, fields []string) (row []string, err error) {
	row = make([]string, len(fields))

	value := reflect.ValueOf(obj)
	for i, field := range fields {
		// check field starts with capital to ensure we only use exported types
		v := value.FieldByName(field)
		if !v.IsValid() {
			v = value.MethodByName(field)
		}
		str, err := valueToString(v)
		if err != nil {
			return nil, err
		}
		row[i] = str

	}
	return

}

// valueToString will convert v to a string by ANY MEANS NECESSARY (either it already is a string, or
func valueToString(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int()), nil
	case reflect.Float32, reflect.Float64:
		// format float at max precision losing without trailing zeroes
		return strconv.FormatFloat(v.Float(), 'f', 2, 64), nil
	case reflect.Ptr:
		return valueToString(reflect.Indirect(v))
	case reflect.Invalid:
		// oh shit ma dudes
		return "", errors.New("v wasn't a valid value!")
	default:
		if isStringer(v) {
			ret := v.MethodByName("String").Call([]reflect.Value{})
			return ret[0].Interface().(string), nil
		}
		return "", errors.New("v wasn't a type we were ready for")
	}
}

func isStringer(v reflect.Value) bool {
	stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	return v.Type().Implements(stringerType)
}

func methodToString(m reflect.Value) (string, error) {
	methodType := m.Type()
	if methodType.NumIn() != 0 {
		return "", errors.New("Wrong number of parameters in methodToString")
	}
	// make sure this method is either func() T or func() T, err
	nOuts := methodType.NumOut()
	if nOuts == 2 {
		retType := methodType.Out(1)
		errorType := reflect.TypeOf((*error)(nil)).Elem()

		if !retType.Implements(errorType) {
			return "", errors.New("2nd value returned from method is not an error")
		}
	}
	if nOuts != 1 && nOuts != 2 {
		return "", errors.New("Method returns wrong number of values")
	}

	ret := m.Call([]reflect.Value{})
	// check error first
	if nOuts == 2 {
		err := ret[1].Interface()
		if err != nil {
			errErr := err.(error)
			return "", errErr
		}
	}
	// if no error, turn the output into a string in the normal fashion
	return valueToString(ret[0])
}
