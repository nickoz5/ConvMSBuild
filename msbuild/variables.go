package msbuild

import (
	"fmt"
	"path/filepath"
	"strings"
)

var savedVars map[string]string = make(map[string]string)

func setVar(key string, value string) {
	if savedVars[key] == "" {
		savedVars[key] = value
		fmt.Printf("Environment: $(%s) = \"%s\"\n", key, value)
	}
}

func SubstituteVar(projectFile string, value string) string {
	for {
		posStart := strings.Index(value, "$(")
		if posStart == -1 {
			break
		}

		posEnd := strings.Index(value[posStart:], ")")
		if posEnd == -1 {
			break
		}

		attrName := value[posStart+2 : posStart+posEnd]

		if attrName != "" {
			// check system variable first
			attrValue, isSystem := getSystemVar(projectFile, attrName)
			if !isSystem {
				attrValue = savedVars[attrName]
			}

			// replace the variable with the saved value
			value = strings.Replace(value, "$("+attrName+")", attrValue, -1)

			if attrValue == "" && !isSystem {
				fmt.Println("Variable not found: $(" + attrName + ")")
			}
		}
	}

	return value
}

func getSystemVar(projectFile string, attrName string) (string, bool) {
	var attrValue string

	found := true

	switch attrName {
	case "MSBuildProjectDirectory":
		attrValue = filepath.Dir(projectFile) + "\\"
	case "MSBuildThisFileDirectory":
		attrValue = filepath.Dir(projectFile) + "\\"
	case "MSBuildToolsPath":
		attrValue = ""
	case "MSBuildToolsVersion":
		attrValue = ""

	default:
		found = false
	}

	return attrValue, found
}
