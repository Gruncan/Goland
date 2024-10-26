package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

const large int64 = 8000 //8000 bytes is considered large here

func crawler(fsys fs.FS, fpath string) []string {
	dir, err := fs.ReadDir(fsys, fpath)
	if err != nil {
		// Simple and not very good error handling
		fmt.Println("Error while crawling: ", fpath, err)
		return nil
	}

	var largeFiles []string

	// First loop collects large files in the current directory
	for _, entry := range dir {
		info, _ := entry.Info()
		if info.Size() > large {
			//fmt.Println("Large file: ", info.Name(), info.Size())
			largeFiles = append(largeFiles, info.Name())
		}
	}

	// Second loop descends as required (also fine to do in one loop: goal is to
	// explore parallelism, not to write the best crawler)
	for _, entry := range dir {
		if entry.IsDir() {
			dirFiles := crawler(fsys, path.Join(fpath, entry.Name()))
			largeFiles = append(largeFiles, dirFiles...) // Slice concat syntax
		}
	}

	return largeFiles
}

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Sprintf("Usage: must provide starting directory as an argument"))
	}

	fmt.Println("Crawing from: ", os.Args[1])

	// Starting directory
	sDir := os.Args[1]
	fileSystem := os.DirFS(sDir)

	// Note: there is a WalkDir function that walks over directories, we don't
	// want to use this since we want to write a concurrent version
	largeFiles := crawler(fileSystem, sDir)

	fmt.Println("Large files: ")
	for _, f := range largeFiles {
		fmt.Println(f)
	}
}
