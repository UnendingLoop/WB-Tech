// Package cmd is an entrance point to the whole logic of the utility
package cmd

import (
	"log"
	"os"

	"cutClone/internal/parser"
	"cutClone/internal/reader"
)

func LaunchApp() {
	// сначала вызов обработчика аргументов ос
	if err := parser.InitSearchParam(); err != nil {
		log.Printf("Warning while parsing OS args: %q", err)
		os.Exit(1)
	}
	// определение входа - файл или стдин
	if len(parser.CC.Source) == 0 {
		err := reader.ProcessInput("")
		if err != nil {
			log.Fatalf("Warning while processing StdIn: %q", err)
		}
	} else {
		var readerErrors []error
		for _, fileName := range parser.CC.Source {
			err := reader.ProcessInput(fileName)
			if err != nil {
				log.Printf("Warning while processing file %q: %q", fileName, err)
				readerErrors = append(readerErrors, err)
			}
		}

		if len(readerErrors) != 0 {
			log.Fatalf("cutClone finished with %d errors", len(readerErrors))
		}
	}
}
