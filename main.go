package main

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"io/ioutil"
	"encoding/xml"
)

type ProjectFile struct {
	Filename string
	ProjectData Project
}
type Project struct {
    XMLName    xml.Name     `xml:"Project"`
	Imports    []Import		`xml:"Import"`
    ItemGroups []ItemGroup  `xml:"ItemGroup"`
}
type Import struct {
	XMLName xml.Name `xml:"Import"`
	Project string   `xml:"Project,attr"`
}
type ItemGroup struct {
	XMLName          xml.Name       `xml:"ItemGroup"`
	BuildProjects    []BuildProject `xml:"BuildProject"`
}
type BuildProject struct {
	XMLName xml.Name `xml:"BuildProject"`
	Include string   `xml:"Include,attr"`
}

var rootProjectFilename string

func main() {
	rootProjectFilename = os.Args[1]
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

func loadProject(filename string) ProjectFile {
	var proj ProjectFile
	proj.Filename = filename

	xmlFile, err := os.Open(proj.Filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return proj
	}
	defer xmlFile.Close()

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

func replaceVariables(proj ProjectFile, value string) string {
	var attrName string

	
	posStart := strings.Index(value, "$(")
	if posStart != -1 {
		posEnd := strings.Index(value, ")")
		if posEnd != -1 {
			attrName = value[posStart+2:posEnd]
		}
	}
	
	if attrName != "" {
		fmt.Printf("Replacing variable $(%s): \"%s\"", attrName, value)

		switch attrName {
		case "MSBuildProjectDirectory":
			newValue := filepath.Dir(proj.Filename)
			value = strings.Replace(value, "$(MSBuildProjectDirectory)", newValue, -1)
			
		case "MSBuildThisFileDirectory":
			newValue := filepath.Dir(proj.Filename)
			value = strings.Replace(value, "$(MSBuildThisFileDirectory)", newValue, -1)
			
		case "BaseDir":
			newValue := filepath.Dir(rootProjectFilename) + "\\"
			value = strings.Replace(value, "$(BaseDir)", newValue, -1)
		}
		fmt.Printf("->\"%s\"\n",  value)
	}
		
	return value
}
