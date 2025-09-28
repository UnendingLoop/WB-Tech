// Package cmd is an entrance point to the whole logic of the utility
package cmd

import (
	"log"
	"os"

	"grepClone/internal/parser"
	"grepClone/internal/processor"
)

func LaunchApp() {
	// сначала вызов обработчика аргументов ос
	if err := parser.InitSearchParam(); err != nil {
		log.Printf("Warning during parse of OS args: %q", err)
		os.Exit(1)
	}
	// определение входа - файл или стдин
	if len(parser.SP.Source) == 0 {
		err := processor.ProcessInput("")
		if err != nil {
			log.Fatalf("Warning while processing StdIn: %q", err)
		}
	} else {
		var processorErrors []error
		for _, fileName := range parser.SP.Source {
			err := processor.ProcessInput(fileName)
			if err != nil {
				log.Printf("Warning while processing file %q: %q", fileName, err)
				processorErrors = append(processorErrors, err)
			}
		}

		if len(processorErrors) != 0 {
			log.Fatalf("grepClone finished with %d errors", len(processorErrors))
		}
	}
}
