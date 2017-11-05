package main

import (
	"fmt"
	"os"

	"github.com/nickoz5/convmsbuild/msbuild"
)

var rootProjectFilename string

func main() {
	if len(os.Args) > 1 {
		rootProjectFilename = os.Args[1]
	}
	fmt.Println("Root project filename: ", rootProjectFilename)

	MakeNewSolution()

	_ = msbuild.LoadProject(rootProjectFilename)

	CreateSolution()
}
