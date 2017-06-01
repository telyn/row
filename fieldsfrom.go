package row

import (
	"reflect"
)

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
