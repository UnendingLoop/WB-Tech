// Package model holds datastructs for running minishell
package model

type Command struct {
	Args       []string
	InputFile  string
	OutputFile string
	Append     bool
}

type Pipeline struct {
	Commands []Command
}

type Conditional struct {
	Pipelines []Pipeline
	Operators []string // например ["&&", "||"]
}
