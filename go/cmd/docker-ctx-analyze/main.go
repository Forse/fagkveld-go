package main

import (
	"fmt"
	"os"
	"path"
)

// This CLI should figure out how large the build context is for a given Dockerfile and context dir.
// It should take 2 arguments - the Dockerfile filepath and the context path.
// Similar to the `docker build` command, it should default to the current working directory if no context path is given.
// It should print the size of the context in bytes to stdout.

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not get current working directory: '%s'\n", err)
		os.Exit(1)
	}

	// Check if a file named 'Dockerfile' exists in the pwd directory
	// The filepath for Dockerfile should come from a command line argument
	filePath := path.Join(pwd, "Dockerfile")
	_, err = os.Stat(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Dockerfile not found: '%s'\n", err)
		os.Exit(1)
	}

	// 1. Check that the context path exists
	// 2. Check that the context path is a directory
	// 3. Find the .dockerignore file in the context path
	// 4. Traverse the context path and sum the size of all files that are not ignored by .dockerignore

	os.Exit(0)
}
