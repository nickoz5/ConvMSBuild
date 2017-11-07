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

	addTargets(sln, proj)

	// determine number of target projects..

	msbuild.CreateSolutionFile(sln, baseDir)
}

func addTargets(sln msbuild.SolutionFile, proj msbuild.ProjectFile) {

	sln.Projects = make(map[string]msbuild.ProjectDefinition)

	for targetidx, _ := range proj.ProjectData.Targets {
		target := &proj.ProjectData.Targets[targetidx]

		// only interested in "Build" targets..
		for buildidx, _ := range target.Builds {
			buildtarget := &target.Builds[buildidx]

			for _, target := range buildtarget.Projects {

				//				sln.Projects[target. Path] = target
			}
		}
	}
}
