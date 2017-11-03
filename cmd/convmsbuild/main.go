package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nickoz5/convmsbuild/pkg/msbuild"
	//"github.com/nickoz5/convmsbuild/pkg/msbuild"
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

	msbuild.Setvariables("MSBuildProjectDirectory", filepath.Dir(projectFile))
	msbuild.Setvariables("", filepath.Dir(projectFile))

	proj := loadProject(rootProjectFilename)

	for _, item := range proj.ProjectData.ItemGroups {
		for _, buildProj := range item.BuildProjects {
			projectFilename := buildProj.Include
			//subProj :=
			loadProject(projectFilename)
		}
	}

}
