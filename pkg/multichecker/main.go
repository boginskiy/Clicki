/*
Package 'multichecker' need for static analize code.

	Package include:
		standart tools from golang.org/x/tools/go/analysis/passes/:
			printf.Analyzer,
			shadow.Analyzer,
			structtag.Analyzer,
			shift.Analyzer

		staticchecker with settings from 'staticcheck.json'
		custom analizer with error of calling os.Exit in main func

	Compilation of package:
		$go build -o <name utility> <path to main.go>

	Use:
		checking one file:
			$./multichecker test.go
		checking one folder:
			$./multichecker ./internal/...
		checking all:
			$./multichecker ./...
*/
package main

import (
	custom "github.com/boginskiy/Clicki/pkg/multichecker/customAnalyzer"
	mult "github.com/boginskiy/Clicki/pkg/multichecker/multichecker"
	static "github.com/boginskiy/Clicki/pkg/multichecker/staticchecker"
)

var config = "staticcheck.json"

func main() {
	// Init Staticcheck with loading config.
	statchecker := static.NewStaticcheck()
	statchecker.LoadConfig(config)

	// Init Multicheck with standart checks.
	multchecker := mult.NewMulticheck(mult.StandartStaticChecker...)

	// Add statchecker.
	multchecker.SetChecks(statchecker.GetChecks()...)

	// Add custom checks.
	multchecker.SetChecks(custom.ErrUsingOsExit)

	multchecker.Start()
}
