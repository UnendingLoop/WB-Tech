// Package model provides a struct for storing parsed flags
package model

// Options - struct for storing all possible flags used for launching the 'sortClone'
type Options struct {
	Column        int    // done
	Numeric       bool   // done
	Reverse       bool   // done
	Unique        bool   // done
	Delimeter     string // done
	Monthly       bool   // done
	IgnSpaces     bool   // done
	CheckIfSorted bool
	HumanSort     bool   // done
	WriteToFile   string // done
}

var OptsContainer = Options{}
