package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func replaceGoPackage(filePath, oldPackage, newPackage string) error {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(filePath + ".tmp")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, oldPackage) {
			line = strings.Replace(line, oldPackage, newPackage, 1)
		}
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	writer.Flush()
	inputFile.Close()
	outputFile.Close()

	return os.Rename(filePath+".tmp", filePath)
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: replace_package <proto_dir> <old_package> <new_package>")
		os.Exit(1)
	}

	protoDir := os.Args[1]
	oldPackage := os.Args[2]
	newPackage := os.Args[3]

	err := filepath.Walk(protoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			err := replaceGoPackage(path, oldPackage, newPackage)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Package replacement complete.")
}
