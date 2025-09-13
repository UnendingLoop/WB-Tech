package processor

import "testing"

func TestUnZipString(t *testing.T) {
	testList := []struct {
		input    string
		expected string
		wantErr  error
	}{
		{"a4bc2d5e", "aaaabccddddde", nil},
		{"abcd", "abcd", nil},
		{"", "", nil},
		{"qwe\\45", "qwe44444", nil},
		{"a1b2c3\\23", "abbccc222", nil},
		{"a123b", "aaaab", nil},
		{"@123b", "@@@@b", nil},
		{"2qwe\\45", "", errFirstNotLetter},
		{"qwe\\\\45", "", errTooManySlashes},
		{"qwe\\45\\", "", errLastRuneIsSlash},
		{"qwe\\g", "", errLetterAfterSlash},
		{"45", "", errFoundOnlyNumbers},
		{"a0b", "", errZeroAsMultiplier},
	}

	for _, test := range testList {
		t.Run(test.input, func(t *testing.T) {
			result, err := UnZipString(test.input)
			if result != test.expected {
				t.Errorf("UnzipString input '%v':\nActual result: '%v',\nExpected result: '%v'\n", test.input, result, test.expected)
			}
			if err != test.wantErr {
				t.Errorf("UnzipString input '%v':\nGot error: '%v'\nExpected error: '%v'\n", test.input, err, test.wantErr)
			}
		})
	}
}
