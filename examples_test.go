package row_test

import (
	"fmt"
	"github.com/BytemarkHosting/row"
)

func ExampleSortedFieldsFrom() {
	// this example demonstrates how to get a list of a struct's fields ordered by field indexes
	// the quotes are necessary, and the indices must be integers
	type Person struct {
		ID                 int    `index:"4"`
		Name               string `index:"5"`
		DriftCompatibility int    `index:"0"`
		Zebras             int    `index:"1"`
	}
	p := Person{ID: 345, Name: "Smelulok", DriftCompatibility: -100}
	fmt.Printf("%+v", row.SortedFieldsFrom(p))
	// Output: [DriftCompatibility Zebras ID Name]
}

func ExampleSortFieldsWithTagName() {
	// this example demonstrates using a custom index name and that SortedFieldsFrom is a shorthand for a longer call to SortFields.
	type Person struct {
		ID                 int    `mycoolindex:"4"`
		Name               string `mycoolindex:"5"`
		DriftCompatibility int    `mycoolindex:"0"`
		Zebras             int    `mycoolindex:"1"`
	}
	p := Person{ID: 345, Name: "Smelulok", DriftCompatibility: -100}
	fmt.Printf("%+v\n", row.SortFields(p, row.FieldsFrom(p), "mycoolindex"))
	fmt.Printf("%+v", row.SortedFieldsFrom(p, "mycoolindex"))
	// Output: [DriftCompatibility Zebras ID Name]
	// [DriftCompatibility Zebras ID Name]
}

func ExampleSortFieldsMixed() {
	// This example demonstrates that un-indexed fields always get sorted to the end, in alphabetical order
	type ComplexStruct struct {
		Shoes    int `index:"3"`
		Ribaldry int
		Hats     int
		Creases  int `index:"4"`
	}
	cs := ComplexStruct{1, 1, 1, 1}
	fmt.Printf("%+v", row.SortFields(cs, row.FieldsFrom(cs)))
	// Output: [Shoes Creases Hats Ribaldry]
}
