package msbuild

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type ProjectFile struct {
	Filename    string
	ProjectData Project
}

type Project struct {
	XMLName        xml.Name        `xml:"Project"`
	Imports        []Import        `xml:"Import"`
	ItemGroups     []ItemGroup     `xml:"ItemGroup"`
	PropertyGroups []PropertyGroup `xml:"PropertyGroup"`
	Targets        []Target        `xml:"Target"`
	PropValues     map[string]string
}
type Import struct {
	XMLName xml.Name
	Project string `xml:"Project,attr"`
}
type ItemGroup struct {
	XMLName       xml.Name
	BuildProjects []BuildProject `xml:"BuildProject"`
}
type BuildProject struct {
	XMLName xml.Name
	Include string `xml:"Include,attr"`
	Project ProjectFile
}
type PropertyGroup struct {
	XMLName xml.Name
	Nodes   []Property `xml:",any"`
}
type Property struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}
type Target struct {
	XMLName xml.Name
	Name    string          `xml:"Name,attr"`
	Builds  []MSBuildTarget `xml:"MSBuild"`
}
type MSBuildTarget struct {
	XMLName      xml.Name
	ProjectNames string `xml:"Projects,attr"`
	Projects     []SolutionFile
}

func LoadProject(filename string) ProjectFile {
	var proj ProjectFile
	proj.Filename = filename
	proj.ProjectData.PropValues = make(map[string]string)

	xmlFile, err := os.Open(proj.Filename)
	if err != nil {
		fmt.Println("Error opening file: ", proj.Filename)
		return proj
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(byteValue, &proj.ProjectData)

	// parse all properties
	for _, group := range proj.ProjectData.PropertyGroups {
		for _, prop := range group.Nodes {
			value := SubstituteVar(proj, prop.Content)
			setVar(proj, prop.XMLName.Local, value)
		}
	}

	for _, item := range proj.ProjectData.Imports {
		projectFilename := SubstituteVar(proj, item.Project)
		imp := LoadProject(projectFilename)

		// import data from imp into proj
		for k, v := range imp.ProjectData.PropValues {
			proj.ProjectData.PropValues[k] = v
		}
	}

	for itemidx, _ := range proj.ProjectData.ItemGroups {
		item := &proj.ProjectData.ItemGroups[itemidx]

		for buildidx, _ := range item.BuildProjects {
			buildproj := &item.BuildProjects[buildidx]

			projectFilename := SubstituteVar(proj, buildproj.Include)
			buildproj.Project = LoadProject(projectFilename)
		}
	}

	for targetidx := range proj.ProjectData.Targets {
		target := &proj.ProjectData.Targets[targetidx]

		// only interested in "Build" targets..
		for buildidx := range target.Builds {
			buildtarget := &target.Builds[buildidx]

			targetnames := strings.Split(buildtarget.ProjectNames, ";")

			count := len(targetnames)
			if count > 0 {
				buildtarget.Projects = make([]SolutionFile, count)

				for targetidx := range targetnames {
					projectFilename := SubstituteVar(proj, targetnames[targetidx])
					sln, _ := LoadSolution(projectFilename)

					buildtarget.Projects[targetidx] = sln

					fmt.Println(buildtarget.Projects)
				}
			}
		}
	}

	return proj
}
