package msbuild

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SolutionFile struct {
	Filename string
	Projects map[string]ProjectDefinition
}

type ProjectDefinition struct {
	TypeGUID     string
	Name         string
	Path         string
	ProjectGUID  string
	Dependancies []string
}

var newSolution SolutionFile

func LoadSolution(filename string) (SolutionFile, int) {
	var solution SolutionFile
	solution.Filename = filename
	solution.Projects = make(map[string]ProjectDefinition)

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
			project, err := parseSolutionProject(filename, line)
			if err == 0 {
				solution.Projects[project.Path] = project
			}
		}
	}

	return solution, 0
}

func parseSolutionProject(solutionfilename string, line string) (ProjectDefinition, int) {
	var proj ProjectDefinition

	projectKeypair := strings.Split(line, "=")
	if len(projectKeypair) != 2 {
		return proj, -1
	}

	proj.TypeGUID = line[9:47]

	projectDefFull := strings.Replace(projectKeypair[1], "\"", "", -1)
	projectDef := strings.Split(projectDefFull, ",")

	proj.Name = strings.Trim(projectDef[0], " ")
	proj.Path = strings.Trim(projectDef[1], " ")
	proj.ProjectGUID = strings.Trim(projectDef[2], " ")

	proj.Path = filepath.Dir(solutionfilename) + "\\" + proj.Path

	return proj, 0
}

func registerProject(proj ProjectDefinition) {
	newSolution.Projects[proj.Name] = proj
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateSolutionFile(sln SolutionFile, baseDir string) {

	f, err := os.Create(sln.Filename)
	checkError(err)

	defer f.Close()

	w := bufio.NewWriter(f)

	//	var wlen int = -1

	// dump the header
	_, err = w.WriteString("\n" +
		"Microsoft Visual Studio Solution File, Format Version 12.00\n" +
		"# Visual Studio 2012\n")
	checkError(err)

	for _, proj := range sln.Projects {
		filename := strings.Replace(proj.Path, baseDir, "", 1)

		_, err = w.WriteString("Project(\"" + proj.TypeGUID + "\") = \"" + proj.Name + "\", \"" + filename + "\", \"" + proj.ProjectGUID + "\"\n")
		checkError(err)

		_, err = w.WriteString("\tProjectSection(ProjectDependencies) = postProject")
		_, err = w.WriteString("\t\t{B1AD7565-468E-4675-B684-FDC9BD1A35EB} = {B1AD7565-468E-4675-B684-FDC9BD1A35EB}")
		_, err = w.WriteString("\tEndProjectSection")

		_, err = w.WriteString("EndProject\n")
		checkError(err)
	}

	_, err = w.WriteString("Global")

	_, err = w.WriteString("\tGlobalSection(SolutionConfigurationPlatforms) = preSolution")
	_, err = w.WriteString("\t\tDebug|Any CPU = Debug|Any CPU")
	_, err = w.WriteString("\t\tDebug|x64 = Debug|x64")
	_, err = w.WriteString("\t\tDebug|x86 = Debug|x86")
	_, err = w.WriteString("\t\tRelease|Any CPU = Release|Any CPU")
	_, err = w.WriteString("\t\tRelease|x64 = Release|x64")
	_, err = w.WriteString("\t\tRelease|x86 = Release|x86")
	_, err = w.WriteString("\tEndGlobalSection")

	_, err = w.WriteString("\tGlobalSection(ProjectConfigurationPlatforms) = postSolution")
	_, err = w.WriteString("\tEndGlobalSection")

	fmt.Printf("Solution file [%s] created with [%d] projects\n", sln.Filename, len(sln.Projects))

	w.Flush()
}
