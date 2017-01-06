package row_test

import (
	"github.com/olekukonko/tablewriter"
	"github.com/telyn/row"
	"os"
)

// Country represents a country
type Country struct {
	Name       string
	Population int     `row.thousands:","`
	HDI        float32 `row.precision:"2"`
	Cities     []string
}

// NumCities could also have the signature 'func (c Country) NumCities() (int, error)
func (c Country) NumCities() int {
	return len(c.Cities)
}

// Output:
// TODO include output
func ExampleCountryTable() error {
	fields := []string{"Name", "Population", "NumCities"}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(fields)

	country := Country{"Argentine Republic", 4341700, 0.836, []string{"Buenos Aires", "Córdoba", "Rosario", "Mendoza", "La Plata"}}

	values, err := row.From(country, fields)
	if err != nil {
		return err
	}
	table.Append(values)
	return nil
}
