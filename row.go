package row

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
			v = value.MethodByName(field)
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

func fieldsFromType(t reflect.Type) (fields []string) {
	fields = make([]string, 0)
	// grab all the field-y methods first, because we indirect later and lose access to the pointery ones
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		//fmt.Printf("testing %s..", m.Name)
		err := methodIsField(m.Type, true)
		if err == nil {
			//	fmt.Println("yes")
			fields = append(fields, m.Name)
		} else {
			//	fmt.Println(err)
		}
	}
	// now indirect if this is a pointery type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return []string{}
	}
	// loop over all struct fields
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields = append(fields, f.Name)
	}
	return
}

// valueToString will convert v to a string by ANY MEANS NECESSARY (either it already is a string, or
func valueToString(v reflect.Value) (string, error) {
	if v.Kind() == reflect.Invalid {
		// oh shit ma dudes
		return "", errors.New("v wasn't a valid value!")
	}
	if isStringer(v) {
		ret := v.MethodByName("String").Call([]reflect.Value{})
		return ret[0].Interface().(string), nil
	}
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
	case reflect.Array, reflect.Slice:
		output := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			elem, err := valueToString(v.Index(i))
			if err != nil {
				return "", err
			}
			output[i] = elem
		}
		return strings.Join(output, "\n"), nil
	case reflect.Ptr:
		if v.IsNil() {
			return "nil", nil
		}
		return valueToString(v.Elem())
	default:
		return "", fmt.Errorf("v (%v) (%T) wasn't a type we were ready for. Its kind is %s", v.Interface(), v.Interface(), v.Kind())
	}
}

func isStringer(v reflect.Value) bool {
	stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	return v.Type().Implements(stringerType)
}

// if error==nil, method is field. if isReceiver is true, expects there to be a receiver (so 1 arg), if not, 0 args
func methodIsField(m reflect.Type, isReceiver bool) error {
	args := 0
	if isReceiver {
		args = 1
	}
	if m.NumIn() != args {
		return errors.New("Wrong number of parameters in methodToString")
	}
	// make sure this method is either func() T or func() T, err
	nOuts := m.NumOut()
	if nOuts == 2 {
		retType := m.Out(1)
		errorType := reflect.TypeOf((*error)(nil)).Elem()

		if !retType.Implements(errorType) {
			return errors.New("2nd value returned from method is not an error")
		}
	}
	if nOuts != 1 && nOuts != 2 {
		return errors.New("Method returns wrong number of values")
	}
	return nil
}

func methodToString(m reflect.Value) (string, error) {
	err := methodIsField(m.Type(), false)
	if err != nil {
		return "", err
	}

	ret := m.Call([]reflect.Value{})
	// check error first
	if m.Type().NumOut() == 2 {
		err := ret[1].Interface()
		if err != nil {
			errErr := err.(error)
			return "", errErr
		}
	}
	// if no error, turn the output into a string in the normal fashion
	return valueToString(ret[0])
}
