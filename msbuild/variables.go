package msbuild

import (
	"fmt"
	"path/filepath"
	"strings"
)

func setVar(ctx ProjectFile, key string, value string) {
	if ctx.ProjectData.PropValues[key] == "" {
		ctx.ProjectData.PropValues[key] = value
	}
}

func SubstituteVar(ctx ProjectFile, value string) string {
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
			attrValue, isSystem := getSystemVar(ctx, attrName)
			if !isSystem {
				attrValue = ctx.ProjectData.PropValues[attrName]
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

func getSystemVar(ctx ProjectFile, attrName string) (string, bool) {
	var attrValue string

	found := true

	switch attrName {
	case "MSBuildProjectDirectory":
		attrValue = filepath.Dir(ctx.Filename) + "\\"
	case "MSBuildThisFileDirectory":
		attrValue = filepath.Dir(ctx.Filename) + "\\"
	case "MSBuildToolsPath":
		attrValue = ""
	case "MSBuildToolsVersion":
		attrValue = ""

	default:
		found = false
	}

	return attrValue, found
}
