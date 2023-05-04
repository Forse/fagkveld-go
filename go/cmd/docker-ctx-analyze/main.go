package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/codeskyblue/dockerignore"
)

// This CLI should figure out how large the build context is for a given Dockerfile and context dir.
// It should take 2 arguments - the Dockerfile filepath and the context path.
// Similar to the `docker build` command, it should default to the current working directory if no context path is given.
// It should print the size of the context in bytes to stdout.

func main() {
	errorLogger := log.New(os.Stderr, "ERROR: ", 0)
	infoLogger := log.New(os.Stdout, "INFO: ", 0)

	config, err := getConfig(errorLogger)
	if err != nil {
		errorLogger.Fatalf("could not read configuration: %s", err)
	}

	infoLogger.Printf("analyzing context based on Dockerfile='%s', context directory='%s'", config.dockerfile, config.context)

	_, err = os.Stat(config.context)
	if err != nil {
		errorLogger.Fatalf("could not read context directory: %s", err)
	}

	_, err = os.Stat(config.dockerfile)
	if err != nil {
		errorLogger.Fatalf("could not read Dockerfile: %s", err)
	}

	dockerIgnorePath := filepath.Join(config.context, ".dockerignore")

	_, err = os.Stat(dockerIgnorePath)
	if err != nil {
		errorLogger.Fatalf("could not read .dockerignore: %s", err)
	}

	ignorePatterns, err := dockerignore.ReadIgnoreFile(dockerIgnorePath)
	if err != nil {
		errorLogger.Fatalf("could not read .dockerignore: %s", err)
	}

	infoLogger.Printf("ignore patterns: %v", ignorePatterns)

	var ignoredSize int64 = 0
	var includedSize int64 = 0
	var files = 0

	err = filepath.Walk(config.context, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errorLogger.Fatalf("could not read file: %s", err)
		}

		if info.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(config.context, path)
		if err != nil {
			return err
		}

		isMatch, err := dockerignore.Matches(relativePath, ignorePatterns)
		if err != nil {
			return err
		}

		infoLogger.Printf("file: %s", relativePath)

		files += 1
		if isMatch {
			ignoredSize += info.Size()
		} else {
			includedSize += info.Size()
		}

		return nil
	})

	if err != nil {
		errorLogger.Fatalf("could not read context directory: %s", err)
	}

	infoLogger.Printf("files: %d", files)
	infoLogger.Printf("context size: %d bytes", includedSize)
	infoLogger.Printf("ignored size: %d bytes", ignoredSize)
}

type config struct {
	dockerfile string
	context    string
}

func getConfig(errorLogger *log.Logger) (*config, error) {
	dockerfileFlag := flag.String("f", "./Dockerfile", "Dockerfile filepath - defaults to './Dockerfile'")

	flag.Parse()

	dockerContextPathFlag := flag.Arg(0)

	if dockerfileFlag == nil || *dockerfileFlag == "" {
		return nil, errors.New("dockerfile filepath not specified")
	}

	if dockerContextPathFlag == "" {
		return nil, errors.New("docker context directory not specified")
	}

	return &config{
		dockerfile: *dockerfileFlag,
		context:    dockerContextPathFlag,
	}, nil
}
