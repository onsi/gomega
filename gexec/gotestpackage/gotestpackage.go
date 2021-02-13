package gotestpackage

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type Package struct {
	destination string
	packagePath string
}

func Get(destination, gopath, packagePath string, envs []string) (*Package, error) {
	if err := getForTest(gopath, packagePath, envs); err != nil {
		return nil, err
	}

	return &Package{
		destination: destination,
		packagePath: packagePath,
	}, nil
}

func isLocalPackage(packagePath string) bool {
	return strings.HasPrefix(packagePath, ".")
}

func getForTest(gopath, packagePath string, env []string) error {
	if isLocalPackage(packagePath) {
		return nil
	}

	return doGet(gopath, packagePath, env, "-t")
}

func doGet(gopath, packagePath string, env []string, args ...string) error {
	args = append(args, packagePath)
	args = append([]string{"get"}, args...)

	goGet := exec.Command("go", args...)
	goGet.Dir = gopath
	goGet.Env = replaceGoPath(os.Environ(), gopath)
	goGet.Env = append(goGet.Env, env...)

	output, err := goGet.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to get %s:\n\nError:\n%s\n\nOutput:\n%s", packagePath, err, string(output))
	}

	return nil
}

func replaceGoPath(environ []string, newGoPath string) []string {
	newEnviron := []string{}
	for _, v := range environ {
		if !strings.HasPrefix(v, "GOPATH=") {
			newEnviron = append(newEnviron, v)
		}
	}
	return append(newEnviron, "GOPATH="+newGoPath)
}

/*
Build uses go test to compile the test package.  The resulting binary is saved off in a temporary directory.
A path pointing to this binary is returned.

CompileTest uses the $GOPATH set in your environment. If $GOPATH is not set and you are using Go 1.8+,
it will use the default GOPATH instead.  It passes the variadic args on to `go test`.
*/
func (p *Package) Build(args ...string) (compiledPath string, err error) {
	return p.doCompileTest(build.Default.GOPATH, nil, args...)
}

/*
BuildWithEnvironment is identical to Build but allows you to specify env vars to be set at build time.
*/
func (p *Package) BuildWithEnvironment(env []string, args ...string) (compiledPath string, err error) {
	return p.doCompileTest(build.Default.GOPATH, env, args...)
}

/*
BuildIn is identical to Build but allows you to specify a custom $GOPATH (the first argument).
*/
func (p *Package) BuildIn(gopath string, args ...string) (compiledPath string, err error) {
	return p.doCompileTest(gopath, nil, args...)
}

func (p *Package) doCompileTest(gopath string, env []string, args ...string) (compiledPath string, err error) {
	executable, err := p.newExecutablePath(gopath, ".test")
	if err != nil {
		return "", err
	}

	cmdArgs := append([]string{"test", "-c"}, args...)
	cmdArgs = append(cmdArgs, "-o", executable, p.packagePath)

	build := exec.Command("go", cmdArgs...)
	build.Env = replaceGoPath(os.Environ(), gopath)
	build.Env = append(build.Env, env...)

	output, err := build.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Failed to build %s:\n\nError:\n%s\n\nOutput:\n%s", p.packagePath, err, string(output))
	}

	return executable, nil
}

func (p *Package) newExecutablePath(gopath string, suffixes ...string) (string, error) {
	if len(gopath) == 0 {
		return "", errors.New("$GOPATH not provided when building " + p.packagePath)
	}

	hash := md5.Sum([]byte(p.packagePath))
	filename := fmt.Sprintf("%s-%x%s", path.Base(p.packagePath), hex.EncodeToString(hash[:]), strings.Join(suffixes, ""))
	executable := filepath.Join(p.destination, filename)

	if runtime.GOOS == "windows" {
		executable += ".exe"
	}

	return executable, nil
}
