package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

// func dumpStringsSliceToFile overwrites existing content in file_path with given contant
func dumpStringsSliceToFile(final_repos []string, file_path string) {
	content := strings.Join(final_repos, "\n")
	ioutil.WriteFile(file_path, []byte(content), 0755)
}

// func sliceContains checks if a string is present in given slice
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}

// func joinSlices add elements of new slice to olde slice,
// only if it does not already exist in old
func joinSlices(new_slice []string, old_slice []string) []string {
	for _, i := range new_slice {
		if !sliceContains(old_slice, i) {
			old_slice = append(old_slice, i)
		}
	}

	return old_slice
}

// func openFile returns a files object given a file path
// creates the file if it does not exist already
func openFile(file_path string) *os.File {
	f, err := os.OpenFile(file_path, os.O_RDWR|os.O_APPEND, 0660)

	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			_, err = os.Create(file_path)
			if err != nil {
				panic(err)
			}
		} else {
			// some other error
			panic(err)
		}
	}

	return f
}

// func parseFileLinesToSlice given a file path string, get each line
// and stores it in a slice
func parseFileLinesToSlice(file_path string) []string {
	f := openFile(file_path)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err := scanner.Err()
	if err != nil && err != io.EOF {
		panic(err)
	}

	return lines
}

// func addNewSliceElementsToFile given a slice of strings representing paths, 
// stores them to the filesystem
func addNewSliceElementsToFile(file_path string, new_repos []string){
	old_repos := parseFileLinesToSlice(file_path)
	final_repos := joinSlices(old_repos, new_repos)
	dumpStringsSliceToFile(final_repos, file_path)
}

// func scanGitFolders returns list of subfolders of `folder` ending with `.git`
// Returns base folder of the repo, the `.git` folder parent
// Recursively searches in the subfolders by passing an existing `folders` slice.
func scanGitFolders(folders []string, folder string) []string {
	// trim the last `/`
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "\\" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, ".git")
				fmt.Printf(path)
				fmt.Printf("\n\n")
				folders = append(folders, path)
				continue
			}

			// skip vendor and node_modules folders because they contain lots of files 
			// and generally aren't git repositories
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}

			folders = scanGitFolders(folders, path)
		}
	}

	return folders
}

// func recursiveScanFolder starts the recursive search for git repositories
// living in the `folder` subtree
func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder);
}

// func getDotFilePath returns dot file fot the repos list.
// Creates it and the enclosing folder if it does not exist.
func getDotFilePath() string {
	current_user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dot_file := current_user.HomeDir + "\\.gogitlocalstats"

	return dot_file
}

// func scan: scans a new folder for Git repositories
func scan(folder string){
	fmt.Printf("Found folders:\n\n")
	repositories_list := recursiveScanFolder(folder)
	dot_file_path := getDotFilePath()
	addNewSliceElementsToFile(dot_file_path, repositories_list)
	fmt.Printf("\n\nSuccessfully added\n\n")
}