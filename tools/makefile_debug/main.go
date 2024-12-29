package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	examplesDir := "../../examples"
	makefile := "../../Makefile"
	launchFile := "../../.vscode/launch.json"

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
	makefileContent.WriteString("generate.makefile-debug:\n")
	makefileContent.WriteString("\tcd tools/makefile_debug && go run .\n\n")
	makefileContent.WriteString("# If we kill a debug process from VSCode the commands are not deleted. This is to clean up for when this happens.\n")
	makefileContent.WriteString("delete.commands:\n")
	makefileContent.WriteString("\tcd tools/delete_commands && go run .\n\n")
	for _, example := range examples {
		makefileContent.WriteString(fmt.Sprintf("example.%s:\n", example))
		makefileContent.WriteString(fmt.Sprintf("\tcd examples/%s && go run .\n\n", example))
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
      "program": "examples/%s",
      "cwd": "examples/%s"
    }`, title(example), example, example))
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
