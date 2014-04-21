package gexec

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var tmpDir string

func init() {
	var err error
	tmpDir, err = ioutil.TempDir("", "test_cmd_main")
	if err != nil {
		tmpDir = ""
	}
}

func Build(packagePath string, args ...string) (string, error) {
	return BuildIn(os.Getenv("GOPATH"), packagePath, args...)
}

func BuildIn(gopath string, packagePath string, args ...string) (string, error) {
	if tmpDir == "" {
		return "", errors.New("Failed to generate temporary directory!")
	}

	if len(gopath) == 0 {
		return "", errors.New("$GOPATH not provided when building " + packagePath)
	}

	executable := filepath.Join(tmpDir, filepath.Base(packagePath))

	cmdArgs := append([]string{"build"}, args...)
	cmdArgs = append(cmdArgs, "-o", executable, packagePath)

	build := exec.Command("go", cmdArgs...)
	build.Env = append([]string{"GOPATH=" + gopath}, os.Environ()...)

	output, err := build.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Failed to build %s:\n\nError:\n%s\n\nOutput:\n%s", packagePath, err, string(output))
	}

	return executable, nil
}

func CleanupBuildArtifacts() {
	if tmpDir != "" {
		os.RemoveAll(tmpDir)
	}
}
