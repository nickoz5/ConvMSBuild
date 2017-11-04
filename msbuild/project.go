package msbuild

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
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
}

func LoadProject(filename string) ProjectFile {
	var proj ProjectFile
	proj.Filename = filename

	xmlFile, err := os.Open(proj.Filename)
	if err != nil {
		fmt.Println("Error opening file: ", proj.Filename)
		return proj
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(byteValue, &proj.ProjectData)

	fmt.Println("Loaded file: ", proj.Filename)

	// parse all properties
	for _, group := range proj.ProjectData.PropertyGroups {
		for _, prop := range group.Nodes {
			value := SubstituteVar(proj.Filename, prop.Content)
			setVar(prop.XMLName.Local, value)
		}
	}

	for _, item := range proj.ProjectData.Imports {
		projectFilename := SubstituteVar(proj.Filename, item.Project)
		LoadProject(projectFilename)
	}

	for _, item := range proj.ProjectData.ItemGroups {
		for _, buildProj := range item.BuildProjects {
			projectFilename := SubstituteVar(proj.Filename, buildProj.Include)
			LoadProject(projectFilename)
		}
	}

	return proj
}
