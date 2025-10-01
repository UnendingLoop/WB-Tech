package printer_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"cutClone/internal/parser"
	"cutClone/internal/printer"
)

func TestSplitAndPrintLine(t *testing.T) {
	cases := []struct {
		name      string
		line      string
		fields    string
		delim     string
		expOutput string
	}{
		{
			name:      "single column",
			line:      "a,b,c,d",
			fields:    "2",
			delim:     ",",
			expOutput: "b",
		},
		{
			name:      "range columns",
			line:      "1,2,3,4,5",
			fields:    "2-4",
			delim:     ",",
			expOutput: "2,3,4",
		},
		{
			name:      "columns from start",
			line:      "x,y,z",
			fields:    "-2",
			delim:     ",",
			expOutput: "x,y",
		},
		{
			name:      "columns reversed",
			line:      "x,y,z",
			fields:    "3-1",
			delim:     ",",
			expOutput: "z,y,x",
		},
		{
			name:      "columns reversed",
			line:      "x,y,z",
			fields:    "4-5",
			delim:     ",",
			expOutput: "",
		},
		{
			name:      "columns reversed",
			line:      "x,y,z",
			fields:    "6-5",
			delim:     ",",
			expOutput: "",
		},
		{
			name:      "columns to end",
			line:      "p,q,r,s",
			fields:    "3-",
			delim:     ",",
			expOutput: "r,s",
		},
		{
			name:      "multiple columns combined",
			line:      "1,2,3,4",
			fields:    "1,3-4,-2",
			delim:     ",",
			expOutput: "1,3,4,1,2",
		},
		{
			name:      "out of range column",
			line:      "a,b",
			fields:    "1-5",
			delim:     ",",
			expOutput: "a,b",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser.CC.Fields = tc.fields
			parser.CC.Delim = tc.delim

			// перехват stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			printer.SplitAndPrintLine(tc.line)

			w.Close()
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			os.Stdout = oldStdout

			actual := buf.String()
			actual = strings.TrimSpace(actual)

			if actual != tc.expOutput {
				t.Errorf("unexpected output:\nGot:  %q\nWant: %q", actual, tc.expOutput)
			}
		})
	}
}
