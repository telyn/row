olekukonko/tablewriter helper
=============================

[![Build Status](https://travis-ci.org/telyn/row.svg?branch=develop)](https://travis-ci.org/telyn/row) [![Coverage Status](https://coveralls.io/repos/github/telyn/row/badge.svg?branch=master)](https://coveralls.io/github/telyn/row?branch=master)

row is a tiny library to help make rows for olekukonko/tablewriter tables (but it can also be used any time you want to select a bunch of fields of misc types from a struct and get a slice of strings)

how to use
----------

```
$ go get github.com/telyn/row
```

```go
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/telyn/row"
)

type Country {
	Name string
	Population int
	Cities []string
}
// NumCities could also have the signature 'func (c Country) NumCities() (int, error)
func (c Country) NumCities() int {
	return len(c.Cities)
}

// Output:
// TODO include output
func Main() {
	fields := []string{"Name", "Population", "NumCities"}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(fields)
	
	country := Country{"Argentine Republic", 4341700, []string{"Buenos Aires","Córdoba", "Rosario", "Mendoza", "La Plata"}}

	table.Append(row.From(country, fields))
}
```
