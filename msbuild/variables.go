package msbuild

import (
	"fmt"
	"strings"
)

var savedVars map[string]string = make(map[string]string)

func SetVar(key string, value string) {
	savedVars[key] = value
}

func SubstituteVar(projectFile string, value string) string {
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
		value = strings.Replace(value, "$("+attrName+")", newValue, -1)

		fmt.Printf("->\"%s\"\n", value)
	}

	return value
}
