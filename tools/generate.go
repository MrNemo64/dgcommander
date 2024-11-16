package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	examplesDir := "examples"
	makefile := "Makefile"
	launchFile := ".vscode/launch.json"

	// Gather example directories
	dirs, err := os.ReadDir(examplesDir)
	if err != nil {
		fmt.Printf("Error reading examples directory: %v\n", err)
		return
	}

	var examples []string
	for _, dir := range dirs {
		if dir.IsDir() {
			examples = append(examples, dir.Name())
		}
	}

	// Generate Makefile content
	var makefileContent bytes.Buffer
	makefileContent.WriteString("# Auto-generated Makefile. Do not edit manually.\n\n")
	makefileContent.WriteString("generate:\n")
	makefileContent.WriteString("\tgo run tools/generate.go\n\n")
	for _, example := range examples {
		makefileContent.WriteString(fmt.Sprintf("example.%s:\n", example))
		makefileContent.WriteString(fmt.Sprintf("\tgo run examples/%s/main.go\n\n", example))
	}

	// Write Makefile
	err = os.WriteFile(makefile, makefileContent.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error writing Makefile: %v\n", err)
		return
	}

	// Generate launch.json content
	launchConfig := `{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
`
	var launchConfigurations []string
	for _, example := range examples {
		launchConfigurations = append(launchConfigurations, fmt.Sprintf(`    {
      "name": "%s",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "examples/%s/main.go",
      "cwd": "${cwd}"
    }`, title(example), example))
	}
	launchConfig += strings.Join(launchConfigurations, ",\n")
	launchConfig += "\n  ]\n}"

	// Write launch.json
	err = os.WriteFile(launchFile, []byte(launchConfig), 0644)
	if err != nil {
		fmt.Printf("Error writing launch.json: %v\n", err)
		return
	}

	fmt.Println("Makefile and launch.json have been regenerated.")
}

func title(s string) string {
	return s
}
