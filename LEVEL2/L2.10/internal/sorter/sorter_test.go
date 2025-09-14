package sorter

import (
	"reflect"
	"testing"

	"sortClone/internal/model"
)

func TestSortBasic(t *testing.T) {
	lines := []string{"banana", "apple", "cherry"}
	expected := []string{"apple", "banana", "cherry"}
	opts := model.Options{}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort() = %v, want %v", got, expected)
	}
}

func TestSortNumeric(t *testing.T) {
	lines := []string{"10", "2", "1", "20"}
	expected := []string{"1", "2", "10", "20"}
	opts := model.Options{Numeric: true}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort(Numeric) = %v, want %v", got, expected)
	}
}

func TestSortNumericDelimeterColumn(t *testing.T) {
	lines := []string{"10;4", "2;3", "1;1", "20;2"}
	expected := []string{"1;1", "20;2", "2;3", "10;4"}
	opts := model.Options{Numeric: true, Delimeter: ";", Column: 2}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort(NumericDelimeterColumn) = %v, want %v", got, expected)
	}
}

func TestSortReverse(t *testing.T) {
	lines := []string{"a", "b", "c"}
	expected := []string{"c", "b", "a"}
	opts := model.Options{Reverse: true}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort(Reverse) = %v, want %v", got, expected)
	}
}

func TestSortMonthly(t *testing.T) {
	lines := []string{"Mar", "Jan", "Dec"}
	expected := []string{"Jan", "Mar", "Dec"}
	opts := model.Options{Monthly: true}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort(Monthly) = %v, want %v", got, expected)
	}
}

func TestSortHuman(t *testing.T) {
	lines := []string{"1K", "2M", "512", "3G"}
	expected := []string{"512", "1K", "2M", "3G"}
	opts := model.Options{HumanSort: true}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort(Human) = %v, want %v", got, expected)
	}
}

func TestSortColumn(t *testing.T) {
	lines := []string{
		"apple 10",
		"banana 2",
		"cherry 5",
	}
	expected := []string{
		"banana 2",
		"cherry 5",
		"apple 10",
	}
	opts := model.Options{Column: 2, Numeric: true}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort(Column+Numeric) = %v, want %v", got, expected)
	}
}

func TestSortIgnoreSpaces(t *testing.T) {
	lines := []string{"apple  ", "banana", "cherry "}
	expected := []string{"apple  ", "banana", "cherry "}
	opts := model.Options{IgnSpaces: true}
	model.OptsContainer = opts
	got := Sort(lines)
	// просто проверяем, что порядок сортировки корректный
	if got[0] != "apple  " || got[1] != "banana" || got[2] != "cherry " {
		t.Errorf("Sort(IgnSpaces) = %v, want %v", got, expected)
	}
}

func TestSortReverseNumericColumn(t *testing.T) {
	lines := []string{
		"a 1",
		"b 10",
		"c 2",
	}
	expected := []string{
		"b 10",
		"c 2",
		"a 1",
	}
	opts := model.Options{Column: 2, Numeric: true, Reverse: true, Delimeter: " "}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort(Column+Numeric+Reverse) = %v, want %v", got, expected)
	}
}

func TestSortUnique(t *testing.T) {
	lines := []string{
		"a",
		"b",
		"a",
		"c",
		"b",
	}
	expected := []string{"a", "b", "c"}
	opts := model.Options{Unique: true}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort(Unique) = %v, want %v", got, expected)
	}
}

func TestSortAllFlags(t *testing.T) {
	lines := []string{
		"Mar 1K",
		"Jan 2M",
		"Feb 512",
		"Mar 2K",
	}
	expected := []string{
		"Feb 512",
		"Mar 1K",
		"Mar 2K",
		"Jan 2M",
	}
	opts := model.Options{
		Column:    2,
		Numeric:   false,
		Reverse:   false,
		Unique:    false,
		Monthly:   true,
		HumanSort: true,
		IgnSpaces: false,
		Delimeter: " ",
	}
	model.OptsContainer = opts
	got := Sort(lines)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Sort(AllFlags) = %v, want %v", got, expected)
	}
}
