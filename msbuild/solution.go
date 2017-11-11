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
	Filename       string
	Projects       map[string]ProjectDefinition
	NestedProjects []string
}

type ProjectDefinition struct {
	TypeGUID     string
	Name         string
	Path         string
	ProjectGUID  string
	Dependencies []string
}

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
		line = strings.TrimSpace(line)

		found, _ := regexp.MatchString("Project\\(\\\"\\{[\\w]{8}\\-[\\w]{4}\\-[\\w]{4}\\-[\\w]{4}\\-[\\w]{12}\\}\\\"\\)", line)
		if found {
			project := scanProjectData(&solution, scanner, line)
			solution.Projects[project.Path] = project
		}

		if line == "GlobalSection(NestedProjects) = preSolution" {
			scanNestedProjects(&solution, scanner)
		}
	}

	return solution, 0
}

func scanProjectData(solution *SolutionFile, scanner *bufio.Scanner, line string) ProjectDefinition {

	// found a valid project definition
	project, _ := parseSolutionProject(solution.Filename, line)
	if project.Name == "netlib" {
		fmt.Println("test..")
	}

	for scanner.Scan() {

		line = scanner.Text()
		line = strings.TrimSpace(line)

		if line == "EndProject" {
			break
		}

		if line == "ProjectSection(ProjectDependencies) = postProject" {
			scanProjectDeps(scanner, &project)
		}
	}

	return project
}

func scanNestedProjects(solution *SolutionFile, scanner *bufio.Scanner) {

	for scanner.Scan() {

		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.Contains(line, "EndGlobalSection") {
			break
		}

		solution.NestedProjects = append(solution.NestedProjects, line)
	}
}

func scanProjectDeps(scanner *bufio.Scanner, proj *ProjectDefinition) {

	for scanner.Scan() {

		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.Contains(line, "EndProjectSection") {
			break
		}

		found := false
		for _, item := range proj.Dependencies {
			if item == line {
				found = true
				break
			}
		}

		if found == false {
			proj.Dependencies = append(proj.Dependencies, line)
		}
	}
}

func parseSolutionProject(solutionfilename string, line string) (ProjectDefinition, int) {
	var proj ProjectDefinition

	projectKeypair := strings.Split(line, "=")
	if len(projectKeypair) != 2 {
		return proj, -1
	}

	proj.TypeGUID = line[9:47]

	// Project type GUID values: https://www.codeproject.com/Reference/720512/List-of-Visual-Studio-Project-Type-GUIDs

	projectDefFull := strings.Replace(projectKeypair[1], "\"", "", -1)
	projectDef := strings.Split(projectDefFull, ",")

	proj.Name = strings.Trim(projectDef[0], " ")
	proj.Path = strings.Trim(projectDef[1], " ")
	proj.ProjectGUID = strings.Trim(projectDef[2], " ")

	proj.Path = filepath.Dir(solutionfilename) + "\\" + proj.Path

	return proj, 0
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

		// dump dependencies
		if len(proj.Dependencies) > 0 {
			_, err = w.WriteString("\tProjectSection(ProjectDependencies) = postProject\n")
			for i := 0; i < len(proj.Dependencies); i++ {
				_, err = w.WriteString("\t\t" + proj.Dependencies[i] + "\n")
			}
			_, err = w.WriteString("\tEndProjectSection\n")
		}

		_, err = w.WriteString("EndProject\n")
		checkError(err)
	}

	_, err = w.WriteString("Global\n")

	_, err = w.WriteString("\tGlobalSection(SolutionConfigurationPlatforms) = preSolution\n")
	_, err = w.WriteString("\t\tDebug|Any CPU = Debug|Any CPU\n")
	_, err = w.WriteString("\t\tDebug|x64 = Debug|x64\n")
	_, err = w.WriteString("\t\tDebug|x86 = Debug|x86\n")
	_, err = w.WriteString("\t\tRelease|Any CPU = Release|Any CPU\n")
	_, err = w.WriteString("\t\tRelease|x64 = Release|x64\n")
	_, err = w.WriteString("\t\tRelease|x86 = Release|x86\n")
	_, err = w.WriteString("\tEndGlobalSection\n")

	_, err = w.WriteString("\tGlobalSection(ProjectConfigurationPlatforms) = postSolution\n")
	_, err = w.WriteString("\tEndGlobalSection\n")

	if len(sln.NestedProjects) > 0 {
		_, err = w.WriteString("\tGlobalSection(NestedProjects) = preSolution\n")
		for i := 0; i < len(sln.NestedProjects); i++ {
			_, err = w.WriteString("\t\t" + sln.NestedProjects[i] + "\n")
		}
		_, err = w.WriteString("\tEndGlobalSection\n")
	}

	_, err = w.WriteString("EndGlobal\n")

	fmt.Printf("Solution file [%s] created with [%d] projects\n", sln.Filename, len(sln.Projects))

	w.Flush()
}
