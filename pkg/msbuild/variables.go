package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func replaceVariables(proj ProjectFile, value string) string {
	var attrName string

	posStart := strings.Index(value, "$(")
	if posStart != -1 {
		posEnd := strings.Index(value, ")")
		if posEnd != -1 {
			attrName = value[posStart+2 : posEnd]
		}
	}

	if attrName != "" {
		fmt.Printf("Replacing variable $(%s): \"%s\"", attrName, value)

		switch attrName {
		case "MSBuildProjectDirectory":
			newValue := filepath.Dir(proj.Filename)
			value = strings.Replace(value, "$(MSBuildProjectDirectory)", newValue, -1)

		case "MSBuildThisFileDirectory":
			newValue := filepath.Dir(proj.Filename)
			value = strings.Replace(value, "$(MSBuildThisFileDirectory)", newValue, -1)

		case "BaseDir":
			newValue := filepath.Dir(rootProjectFilename) + "\\"
			value = strings.Replace(value, "$(BaseDir)", newValue, -1)
		}
		fmt.Printf("->\"%s\"\n", value)
	}

	return value
}
