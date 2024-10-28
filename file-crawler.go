package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

const large int64 = 320000 //8000 bytes is considered large here

func crawlerCon(fsys fs.FS, fpath string, largeFileC chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	dir, err := fs.ReadDir(fsys, fpath)
	if err != nil {
		// Simple and not very good error handling
		fmt.Println("Error while crawling: ", fpath, err)
		return
	}

	for _, entry := range dir {
		if entry.IsDir() {
			wg.Add(1)
			go crawlerCon(fsys, filepath.Join(fpath, entry.Name()), largeFileC, wg)
		} else {
			info, _ := entry.Info()
			if info.Size() > large {
				largeFileC <- info.Name()
			}
		}
	}

}

func crawler(fsys fs.FS, fpath string) []string {
	var wg sync.WaitGroup

	largeFileC := make(chan string, 100)

	wg.Add(1)
	go crawlerCon(fsys, fpath, largeFileC, &wg)

	var largeFiles []string

	go func() {
		wg.Wait()
		close(largeFileC)
	}()

	for file := range largeFileC {
		largeFiles = append(largeFiles, file)
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
	largeFiles := crawler(fileSystem, ".")

	fmt.Println("Large files: ")
	for _, f := range largeFiles {
		fmt.Println(f)
	}
}
