// Package parser parses input line of commands dividing them into pipes and commands in model structure
package parser

import (
	"strings"

	"miniShell/internal/model"
)

func ParsePipes(input string) model.Pipeline {
	pipes := strings.Split(input, "|")

	var pipeline model.Pipeline

	for _, v := range pipes {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		args := strings.Fields(v)
		cmd := parseCommands(args)
		pipeline.Commands = append(pipeline.Commands, cmd)

	}
	return pipeline
}

func ParseConditional(input string) model.Conditional {
	var cond model.Conditional
	tokens := strings.Fields(input)

	var currentCmd []string

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		if tok == "&&" || tok == "||" {
			pipeline := ParsePipes(strings.Join(currentCmd, " "))
			cond.Pipelines = append(cond.Pipelines, pipeline)
			cond.Operators = append(cond.Operators, tok)
			currentCmd = nil
		} else {
			currentCmd = append(currentCmd, tok)
		}
	}

	if len(currentCmd) > 0 {
		pipeline := ParsePipes(strings.Join(currentCmd, " "))
		cond.Pipelines = append(cond.Pipelines, pipeline)
	}

	return cond
}

func parseCommands(tokens []string) model.Command {
	var cmd model.Command
	i := 0
	for i < len(tokens) {
		tok := tokens[i]
		switch tok {
		case ">":
			if i+1 < len(tokens) {
				cmd.OutputFile = tokens[i+1]
				cmd.Append = false
				i += 2
			} else {
				i++
			}
		case ">>":
			if i+1 < len(tokens) {
				cmd.OutputFile = tokens[i+1]
				cmd.Append = true
				i += 2
			} else {
				i++
			}
		case "<":
			if i+1 < len(tokens) {
				cmd.InputFile = tokens[i+1]
				i += 2
			} else {
				i++
			}
		default:
			cmd.Args = append(cmd.Args, tok)
			i++
		}
	}
	return cmd
}
