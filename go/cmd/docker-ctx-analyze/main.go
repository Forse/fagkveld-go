package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/moby/buildkit/frontend/dockerfile/dockerignore"
	"github.com/moby/patternmatcher"
)

// This CLI should figure out how large the build context is for a given Dockerfile and context dir.
// It should take 2 arguments - the Dockerfile filepath and the context path.
// Similar to the `docker build` command, it should default to the current working directory if no context path is given.
// It should print the size of the context in bytes to stdout.

var errorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", 0)
var infoLogger *log.Logger = log.New(os.Stdout, "INFO: ", 0)

func main() {
	config, err := getConfig()
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

	err = readDockerignore(config)
	if err != nil {
		errorLogger.Fatalf("could not read .dockerignore: %s", err)
	}

	infoLogger.Printf("ignore patterns: %v", config.ignorePatterns)

	if err != nil {
		errorLogger.Fatalf("could not read context directory: %s", err)
	}

	analysis, err := analyzeDockerContext(config)
	if err != nil {
		errorLogger.Fatalf("could not complete analysis: %s", err)
	}

	infoLogger.Printf("files: %d", analysis.files)
	infoLogger.Printf("context size: %d bytes", analysis.includedSize)
	infoLogger.Printf("ignored size: %d bytes", analysis.ignoredSize)
}

type analysis struct {
	ignoredSize  int64
	includedSize int64
	files        int64
}

func analyzeDockerContext(config *config) (*analysis, error) {
	analysis := &analysis{}

	err := filepath.Walk(config.context, func(path string, info os.FileInfo, err error) error {
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

		relativePath = filepath.ToSlash(relativePath)
		isMatch, err := patternmatcher.Matches(relativePath, config.ignorePatterns)
		if err != nil {
			return err
		}

		infoLogger.Printf("file: %s", relativePath)

		analysis.files += 1
		if isMatch {
			analysis.ignoredSize += info.Size()
		} else {
			analysis.includedSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return analysis, nil
}

// From https://github.com/docker/cli/blob/f7600fb5390973c29315024ac2a9c0777735e7ee/cli/command/image/build/dockerignore.go#L13-L26
func readDockerignore(config *config) error {
	f, err := os.Open(filepath.Join(config.context, ".dockerignore"))
	switch {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	}
	defer f.Close()

	config.ignorePatterns, err = dockerignore.ReadAll(f)
	if err != nil {
		return err
	}

	if keep, _ := patternmatcher.Matches(".dockerignore", config.ignorePatterns); keep {
		config.ignorePatterns = append(config.ignorePatterns, "!.dockerignore")
	}

	dockerfile := filepath.ToSlash(config.dockerfile)
	if keep, _ := patternmatcher.Matches(dockerfile, config.ignorePatterns); keep {
		config.ignorePatterns = append(config.ignorePatterns, "!"+dockerfile)
	}

	return nil
}

type config struct {
	dockerfile     string
	context        string
	ignorePatterns []string
}

func getConfig() (*config, error) {
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
		dockerfile:     *dockerfileFlag,
		context:        dockerContextPathFlag,
		ignorePatterns: make([]string, 0),
	}, nil
}
