package prettyprint

import (
	"errors"
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
