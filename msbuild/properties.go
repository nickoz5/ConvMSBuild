package msbuild


import (
	"fmt"
	"strings"
)

type PropertyGroup struct (
	XMLName xml.Name `xml:"PropertyGroup"`
)

type PropertyValue struct (
	XMLName xml.Name `xml:"PropertyGroup"`
)