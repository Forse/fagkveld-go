package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/moby/buildkit/frontend/dockerfile/dockerignore"
	"github.com/moby/patternmatcher"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// This CLI should figure out how large the build context is for a given Dockerfile and context dir.
// It should take 2 arguments - the Dockerfile filepath and the context path.
// Similar to the `docker build` command, it should default to the current working directory if no context path is given.
// It should print the size of the context in bytes to stdout.

var errorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", 0)

type errMsg error
type argsParsed struct{ config *config }
type configurationRead struct{ config *config }
type analysisCompleted struct{ analysis *analysis }

func main() {
	app := tea.NewProgram(initialModel())

	go backgroundWorker(app.Send)

	if _, err := app.Run(); err != nil {
		errorLogger.Fatalf("could not run analysis: %s", err)
	}
}

func backgroundWorker(sendMessage func(msg tea.Msg)) {
	config, err := getConfig()
	if err != nil {
		sendMessage(fmt.Errorf("could not read configuration: %s", err))
	}

	sendMessage(argsParsed{config: config})

	_, err = os.Stat(config.context)
	if err != nil {
		sendMessage(fmt.Errorf("could not read context directory: %s", err))
	}

	_, err = os.Stat(config.dockerfile)
	if err != nil {
		sendMessage(fmt.Errorf("could not read Dockerfile: %s", err))
	}

	err = readDockerignore(config)
	if err != nil {
		sendMessage(fmt.Errorf("could not read .dockerignore: %s", err))
	}

	sendMessage(configurationRead{config: config})

	analysis, err := analyzeDockerContext(config)
	if err != nil {
		sendMessage(fmt.Errorf("could not complete analysis: %s", err))
	}

	sendMessage(analysisCompleted{analysis: analysis})
}

type model struct {
	stopwatch stopwatch.Model
	spinner   spinner.Model
	config    *config
	analysis  *analysis
	err       error
	quitting  bool
}

func initialModel() model {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	timer := stopwatch.NewWithInterval(time.Millisecond)

	return model{spinner: spin, stopwatch: timer}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.stopwatch.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	updateSpinner := func() {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	updateTimer := func() {
		var cmd tea.Cmd
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			break
		}

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case argsParsed:
		m.config = msg.config
		updateTimer()
		updateSpinner()

	case configurationRead:
		m.config = msg.config
		updateTimer()
		updateSpinner()

	case analysisCompleted:
		m.analysis = msg.analysis
		m.stopwatch.Stop()
		updateTimer()
		updateSpinner()
		cmds = append(cmds, tea.Quit)

	default:
		updateTimer()
		updateSpinner()
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var str string
	if m.analysis == nil {
		str += fmt.Sprintf("\n   %s Analyzing docker context...press q to quit", m.spinner.View())
		str += fmt.Sprintf("\n      Elapsed time: %s\n\n", m.stopwatch.View())
	} else {
		//lint:ignore S1039 For alignment
		str += fmt.Sprintf("\n   ✅ Analysis done!")
		str += fmt.Sprintf("\n      Elapsed time: %s\n\n", m.stopwatch.View())

		str += fmt.Sprintf("   Files: %d\n", m.analysis.includedFiles)
		str += fmt.Sprintf("   Context size: %s\n", humanize.Bytes(m.analysis.includedSize))
		str += fmt.Sprintf("   Ignored size: %s\n", humanize.Bytes(m.analysis.ignoredSize))
	}

	if m.quitting {
		return str + "\n"
	}
	return str
}

type analysis struct {
	ignoredSize   uint64
	includedSize  uint64
	includedFiles uint64
	ignoredFiles  uint64
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

		if isMatch {
			analysis.ignoredSize += uint64(info.Size())
			analysis.ignoredFiles += 1
		} else {
			analysis.includedSize += uint64(info.Size())
			analysis.includedFiles += 1
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
