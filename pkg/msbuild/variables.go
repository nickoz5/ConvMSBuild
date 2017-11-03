package msbuild

import (
	"fmt"
	"strings"
)

var savedVars map[string]string

func SetVariables(vars map[string]string) {
	savedVars = vars
}
func ReplaceVariables(projectFile string, value string) string {
	var attrName string

	posStart := strings.Index(value, "$(")
	if posStart != -1 {
		posEnd := strings.Index(value, ")")
		if posEnd != -1 {
			attrName = value[posStart+2 : posEnd]
		}
	}

	if attrName != "" {
		fmt.Printf("Substituting var $(%s): \"%s\"", attrName, value)

		newValue := savedVars[attrName]
		value = strings.Replace(value, attrName, newValue, -1)

		fmt.Printf("->\"%s\"\n", value)
	}

	return value
}
