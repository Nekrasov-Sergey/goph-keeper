// Package buildinfo содержит информацию о сборке приложения.
package buildinfo

import (
	"fmt"
)

var (
	// buildVersion — версия сборки.
	buildVersion = "N/A"
	// buildDate — дата сборки.
	buildDate = "N/A"
	// buildCommit — коммит сборки.
	buildCommit = "N/A"
)

// Print выводит информацию о сборке: версию, дату и коммит.
func Print() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
}
