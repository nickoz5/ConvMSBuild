package msbuild

import (
	"bufio"
	"fmt"
	"os"
)

func loadSolution(slnFilename string) {
	fmt.Println("Loading solution: " + slnFilename)

	inFile, _ := os.Open(slnFilename)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
