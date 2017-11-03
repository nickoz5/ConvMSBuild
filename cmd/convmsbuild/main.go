package main

import (
	"fmt"
	"os"
)

type ProjectFile struct {
	Filename    string
	ProjectData Project
}

var rootProjectFilename string

func main() {
	if len(os.Args) > 1 {
		rootProjectFilename = os.Args[1]
	}
	fmt.Println("Root project filename: ", rootProjectFilename)

	proj := loadProject(rootProjectFilename)

	for _, item := range proj.ProjectData.ItemGroups {
		for _, buildProj := range item.BuildProjects {
			projectFilename := buildProj.Include
			//subProj :=
			loadProject(projectFilename)
		}
	}

}
