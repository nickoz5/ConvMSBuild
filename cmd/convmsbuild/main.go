package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/nickoz5/convmsbuild/msbuild"
)

var rootProjectFilename string

func main() {
	rootProjectFilenamePtr := flag.String("proj", "", "specifies the project to load and convert")
	outputFilenamePtr := flag.String("out", "output.sln", "specifies the output solution filename")
	flag.Parse()

	rootProjectFilename := *rootProjectFilenamePtr
	outputFilename := *outputFilenamePtr

	if rootProjectFilename == "" {
		flag.Usage()
		return
	}

	fmt.Println("Root project filename: ", rootProjectFilename)

	baseDir := filepath.Dir(rootProjectFilename) + "\\"

	proj := msbuild.LoadProject(rootProjectFilename)

	//
	var sln msbuild.SolutionFile
	sln.Filename = baseDir + outputFilename

	sln.Projects = make(map[string]msbuild.ProjectDefinition)
	addTargets(&sln, proj)

	fmt.Println("Found project definition total: ", len(sln.Projects))

	// determine number of target projects..

	msbuild.CreateSolutionFile(sln, baseDir)
}

func addTargets(sln *msbuild.SolutionFile, proj msbuild.ProjectFile) {

	for _, item := range proj.ProjectData.ItemGroups {
		// only interested in "Build" targets..
		for _, buildproj := range item.BuildProjects {
			// recurse on all referenced projects
			addTargets(sln, buildproj.Project)
		}
	}

	fmt.Println("Scanning msbuild project: ", proj.Filename)

	for _, t := range proj.ProjectData.Targets {

		// only interested in "Build" targets..
		for _, b := range t.Builds {

			for _, sproj := range b.Projects {

				// these are our solution files..
				for k, v := range sproj.Projects {

					if _, ok := sln.Projects[k]; ok == false {
						// this is a VS project file..
						fmt.Println("Found project definition: ", v)

						sln.Projects[k] = v
					}
				}
			}
		}
	}
}
