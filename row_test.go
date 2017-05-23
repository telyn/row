package row

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type testStringer struct {
}

func (ts testStringer) String() string {
	return "test stringer"
}

func TestValueToString(t *testing.T) {
	someint := 6
	somefloat := 342.23
	var someNilPtr *int
	var someNilStringer *testStringer
	tests := []struct {
		Value  interface{}
		String string
	}{
		{"test", "test"},
		{7, "7"},
		{&someint, "6"},
		{false, "false"},
		{53.345, "53.34"},
		{&somefloat, "342.23"},
		{testStringer{}, "test stringer"},
		{&testStringer{}, "test stringer"},
		{someNilPtr, "nil"},
		{someNilStringer, "nil"},
	}

	for i, test := range tests {
		str, err := valueToString(reflect.ValueOf(test.Value))
		if err != nil {
			t.Fatalf("testValueToString %d error: %v", i, err)
		}
		if str != test.String {
			t.Errorf("testValueToString %d failed: expecting '%s', got '%s'", i, test.String, str)
		}
	}

	str, err := valueToString(reflect.ValueOf(complex(3.4, 5.6)))
	if err == nil {
		t.Fatalf("testValueToString complex didn't return error.")
	}
	if str != "" {
		t.Fatalf("testValueToString complex didn't return empty string: %v", str)
	}

	str, err = valueToString(reflect.Value{})
	if err == nil {
		t.Fatalf("testValueToString ZeroValue didn't return error.")
	}
	if str != "" {
		t.Fatalf("testValueToString ZeroValue didn't return empty string: %v", str)
	}

}

type structWithMethods struct{}

// repeat for int, float, string, stringer...
func (s structWithMethods) Int() int {
	return 9
}
func (s structWithMethods) IntNoError() (int, error) {
	return 143, nil
}
func (s structWithMethods) IntError() (int, error) {
	return 634, errors.New("int with error")
}

func (s structWithMethods) Float() float32 {
	return 9.1
}
func (s structWithMethods) FloatNoError() (float32, error) {
	return 143.3, nil
}
func (s structWithMethods) FloatError() (float32, error) {
	return 634.8, errors.New("float with error")
}

func (s structWithMethods) String() string {
	return "string"
}
func (s structWithMethods) StringNoError() (string, error) {
	return "string with no error", nil
}
func (s structWithMethods) StringError() (string, error) {
	return "string with error", errors.New("string with error")
}

func (s structWithMethods) Stringer() testStringer {
	return testStringer{}
}
func (s structWithMethods) StringerNoError() (testStringer, error) {
	return testStringer{}, nil
}
func (s structWithMethods) StringerError() (testStringer, error) {
	return testStringer{}, errors.New("float with error")
}

func (s structWithMethods) StringerPtr() *testStringer {
	return &testStringer{}
}
func (s structWithMethods) StringerPtrNoError() (*testStringer, error) {
	return &testStringer{}, nil
}
func (s structWithMethods) StringerPtrError() (*testStringer, error) {
	return &testStringer{}, errors.New("float with error")
}
func (s structWithMethods) BadFuncTypeTooManyRets() (string, error, error) {
	return "result", errors.New("error 1"), errors.New("error 2")
}
func (s structWithMethods) BadFuncTakesArguments(v string) string {
	return v
}
func (s structWithMethods) BadFunc2ndRetNotError() (string, string) {
	return "result one", "result two"
}

func TestMethodToString(t *testing.T) {
	s := structWithMethods{}
	tests := []struct {
		Method      string
		String      string
		ShouldError bool
	}{
		{"Int", "9", false},
		{"IntNoError", "143", false},
		{"IntError", "", true},
		{"Float", "9.10", false},
		{"FloatNoError", "143.30", false},
		{"FloatError", "", true},
		{"String", "string", false},
		{"StringNoError", "string with no error", false},
		{"StringError", "", true},
		{"Stringer", "test stringer", false},
		{"StringerNoError", "test stringer", false},
		{"StringerError", "", true},
		{"StringerPtr", "test stringer", false},
		{"StringerPtrNoError", "test stringer", false},
		{"StringerPtrError", "", true},
		{"BadFuncTypeTooManyRets", "", true},
		{"BadFuncTakesArguments", "", true},
		{"BadFunc2ndRetNotError", "", true},
	}
	sVal := reflect.ValueOf(s)
	for i, test := range tests {
		m := sVal.MethodByName(test.Method)
		str, err := methodToString(m)
		if test.ShouldError && err == nil {
			t.Errorf("testMethodToString %d (%s) expected to error but didn't", i, test.Method, err)
		} else if !test.ShouldError && err != nil {
			t.Errorf("testMethodToString %d (%s) expect nil error but: %v", i, test.Method, err)
		}

		if str != test.String {
			t.Errorf("testMethodToString %d (%s) failed: expecting '%s', got '%s'", i, test.Method, test.String, str)
		}
	}
}

type Craftiness int

const (
	Pedestrian Craftiness = iota
	Crafty
	Racoonesque
)

func (c Craftiness) String() string {
	switch c {
	case Pedestrian:
		return "Pedestrian"
	case Crafty:
		return "Crafty"
	case Racoonesque:
		return "Racoon"
	}
	return fmt.Sprintf("%d craftiness", c)
}

type TestStruct struct {
	Colour     string
	Weight     int
	Gamma      float32
	Void       *string
	Craftiness Craftiness
	Complexity complex64
}

func TestFieldsFrom(t *testing.T) {
	tests := []struct {
		In       interface{}
		Expected []string
	}{{
		In:       TestStruct{},
		Expected: []string{"Colour", "Weight", "Gamma", "Void", "Craftiness", "Complexity"},
	}, {
		In:       []TestStruct{},
		Expected: []string{"Colour", "Weight", "Gamma", "Void", "Craftiness", "Complexity"},
	}, {
		In:       &TestStruct{},
		Expected: []string{"Colour", "Weight", "Gamma", "Void", "Craftiness", "Complexity"},
	}, {
		In:       []*TestStruct{},
		Expected: []string{"Colour", "Weight", "Gamma", "Void", "Craftiness", "Complexity"},
	}, {
		In:       "",
		Expected: []string{},
	}, {
		In:       83,
		Expected: []string{},
	}}

	for i, test := range tests {
		actual := FieldsFrom(test.In)
		if !reflect.DeepEqual(test.Expected, actual) {
			t.Errorf("TestFieldsFrom %d FAIL: \r\n expected: %#v\r\n actual:  %#v", i, test.Expected, actual)
		}
	}
}

func TestRowFrom(t *testing.T) {

	testObj := TestStruct{
		Colour:     "Green",
		Weight:     45,
		Gamma:      0.3,
		Void:       nil,
		Craftiness: Crafty,
		Complexity: complex(3.4, 4.5),
	}

	tests := []struct {
		In  []string
		Out []string
	}{
		{
			In:  []string{"Gamma", "Craftiness", "PowerConsumption", "Void"},
			Out: []string{"0.30", "Crafty", "no field called PowerConsumption", "nil"},
		},
	}

	for i, test := range tests {
		out, err := From(testObj, test.In)
		if err != nil {
			t.Errorf("TestRowFrom %d (%v) ERROR: %s", i, test.In, err.Error())
		}
		if !reflect.DeepEqual(out, test.Out) {
			t.Errorf("TestRowFrom %d (%v) FAIL:\r\n%#v\r\n%#v", i, test.In, out, test.Out)
		}
	}

	out, err := From(testObj, []string{"Complexity"})
	if err == nil {
		t.Errorf("TestRowFrom Complexity NO ERROR")
	}
	if !reflect.DeepEqual(out, []string(nil)) {
		t.Errorf("TestRowFrom Complexity FAIL:\r\n%#v", out)
	}

}
