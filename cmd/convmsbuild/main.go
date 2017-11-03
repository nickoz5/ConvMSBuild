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

	proj := msbuild.LoadProject(rootProjectFilename)

	for _, item := range proj.ProjectData.ItemGroups {
		for _, buildProj := range item.BuildProjects {
			projectFilename := buildProj.Include
			msbuild.LoadProject(projectFilename)
		}
	}

}
