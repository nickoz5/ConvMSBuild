package msbuild

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

)

type Project struct {
	XMLName    xml.Name    `xml:"Project"`
	Imports    []Import    `xml:"Import"`
	ItemGroups []ItemGroup `xml:"ItemGroup"`
}
type Import struct {
	XMLName xml.Name `xml:"Import"`
	Project string   `xml:"Project,attr"`
}
type ItemGroup struct {
	XMLName       xml.Name       `xml:"ItemGroup"`
	BuildProjects []BuildProject `xml:"BuildProject"`
}
type BuildProject struct {
	XMLName xml.Name `xml:"BuildProject"`
	Include string   `xml:"Include,attr"`
}

func loadProject(filename string) ProjectFile {
	var proj ProjectFile
	proj.Filename = filename

	xmlFile, err := os.Open(proj.Filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return proj
	}
	defer xmlFile.Close()

	msbuild.Setvariables("MSBuildProjectDirectory", filepath.Dir(projectFile))
	msbuild.Setvariables("MSBuildThisFileDirectory", filepath.Dir(projectFile))
	msbuild.Setvariables("BaseDir", filepath.Dir(projectFile))

	byteValue, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(byteValue, &proj.ProjectData)

	fmt.Println("Loaded file: ", proj.Filename)
	for _, value := range proj.ProjectData.Imports {
		value.Project = replaceVariables(proj, value.Project)

		// TODO import properties..
	}

	for i1, _ := range proj.ProjectData.ItemGroups {
		for i2, _ := range proj.ProjectData.ItemGroups[i1].BuildProjects {
			item := &proj.ProjectData.ItemGroups[i1].BuildProjects[i2]
			item.Include = replaceVariables(proj, item.Include)
		}
	}

	return proj
}
