package parser

import (
	"reflect"
	"testing"
)

func TestParsePipes(t *testing.T) {
	p := ParsePipes("echo hello | tr a-z A-Z")
	if len(p.Commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(p.Commands))
	}
	if !reflect.DeepEqual(p.Commands[0].Args, []string{"echo", "hello"}) {
		t.Errorf("unexpected args: %v", p.Commands[0].Args)
	}
	if !reflect.DeepEqual(p.Commands[1].Args, []string{"tr", "a-z", "A-Z"}) {
		t.Errorf("unexpected args: %v", p.Commands[1].Args)
	}
}

func TestParseConditional(t *testing.T) {
	c := ParseConditional("echo hi && ls || echo fail")
	if len(c.Pipelines) != 3 {
		t.Fatalf("expected 3 pipelines, got %d", len(c.Pipelines))
	}
	if len(c.Operators) != 2 || c.Operators[0] != "&&" || c.Operators[1] != "||" {
		t.Errorf("unexpected operators: %v", c.Operators)
	}
}

func TestParseCommands_Redirects(t *testing.T) {
	cmd := parseCommands([]string{"cat", "<", "in.txt", ">", "out.txt"})
	if cmd.InputFile != "in.txt" {
		t.Errorf("expected input file in.txt, got %q", cmd.InputFile)
	}
	if cmd.OutputFile != "out.txt" {
		t.Errorf("expected output file out.txt, got %q", cmd.OutputFile)
	}
	if !reflect.DeepEqual(cmd.Args, []string{"cat"}) {
		t.Errorf("unexpected args: %v", cmd.Args)
	}
}
