package msbuild

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type solutionFile struct {
	Filename string
	Projects map[string]projectDefinition
}

type projectDefinition struct {
	Name         string
	Path         string
	ProjectGUID  string
	Dependancies []string
}

var newSolution solutionFile

func loadSolution(filename string) (solutionFile, int) {
	var solution solutionFile
	solution.Filename = filename
	solution.Projects = make(map[string]projectDefinition)

	fmt.Printf("Loading solution: " + filename + "... ")

	inFile, _ := os.Open(filename)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	foundSignature := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Microsoft Visual Studio Solution File") {
			foundSignature = true
			break
		}
	}

	if !foundSignature {
		return solution, -1
	}

	// scan for projects
	for scanner.Scan() {
		line := scanner.Text()

		// Project definition format:
		// Project("{Project Type GUID}") = "projectname", "project.vcxproj", "{Project GUID}"

		// Project type GUID values:
		// https://www.codeproject.com/Reference/720512/List-of-Visual-Studio-Project-Type-GUIDs

		found, _ := regexp.MatchString("Project\\(\\\"\\{[\\w]{8}\\-[\\w]{4}\\-[\\w]{4}\\-[\\w]{4}\\-[\\w]{12}\\}\\\"\\)", line)

		if found {
			if !scanner.Scan() {
				break
			}

			endline := scanner.Text()
			if endline != "EndProject" {
				continue
			}

			// found a valid project definition
			project, err := parseSolutionProject(line)
			if err == 0 {
				solution.Projects[project.Path] = project
			}
		}
	}

	return solution, 0
}

func parseSolutionProject(line string) (projectDefinition, int) {
	var proj projectDefinition

	projectKeypair := strings.Split(line, "=")
	if len(projectKeypair) != 2 {
		return proj, -1
	}

	projectDefFull := strings.Replace(projectKeypair[1], "\"", "", -1)
	projectDef := strings.Split(projectDefFull, ",")

	proj.Name = strings.Trim(projectDef[0], " ")
	proj.Path = strings.Trim(projectDef[1], " ")
	proj.ProjectGUID = strings.Trim(projectDef[2], " ")

	registerProject(proj)

	return proj, 0
}

func MakeNewSolution() {
	newSolution.Projects = make(map[string]projectDefinition)
}

func registerProject(proj projectDefinition) {
	newSolution.Projects[proj.Name] = proj
}

func CreateSolution() {
	fmt.Println(newSolution)
}
