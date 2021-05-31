package gexec_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var packagePath = "./_fixture/firefly"

var _ = Describe(".Build", func() {
	When("there have been previous calls to CompileTest", func() {
		BeforeEach(func() {
			_, err := gexec.CompileTest(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("compiles the specified package", func() {
			compiledPath, err := gexec.Build(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
			Expect(filepath.Base(compiledPath)).Should(MatchRegexp(`firefly(|.exe)$`))
		})

		Context("and CleanupBuildArtifacts has been called", func() {
			BeforeEach(func() {
				gexec.CleanupBuildArtifacts()
			})

			It("compiles the specified package", func() {
				fireflyPath, err := gexec.Build(packagePath)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fireflyPath).Should(BeAnExistingFile())
			})
		})
	})

	When("there have been previous calls to Build", func() {
		BeforeEach(func() {
			_, err := gexec.Build(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("compiles the specified package", func() {
			compiledPath, err := gexec.Build(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})

		Context("and CleanupBuildArtifacts has been called", func() {
			BeforeEach(func() {
				gexec.CleanupBuildArtifacts()
			})

			It("compiles the specified package", func() {
				fireflyPath, err := gexec.Build(packagePath)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fireflyPath).Should(BeAnExistingFile())
			})
		})
	})
})

var _ = Describe(".BuildWithEnvironment", func() {
	var err error
	env := []string{
		"GOOS=linux",
		"GOARCH=amd64",
	}

	It("compiles the specified package with the specified env vars", func() {
		compiledPath, err := gexec.BuildWithEnvironment(packagePath, env)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(compiledPath).Should(BeAnExistingFile())
	})

	It("returns the environment to a good state", func() {
		_, err = gexec.BuildWithEnvironment(packagePath, env)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(os.Environ()).ShouldNot(ContainElement("GOOS=linux"))
	})
})

var _ = Describe(".BuildIn", func() {
	const (
		target = "github.com/onsi/gomega/gexec/_fixture/firefly/"
	)

	var (
		original string
		gopath   string
	)

	BeforeEach(func() {
		var err error
		original = os.Getenv("GOPATH")
		gopath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		copyFile(filepath.Join("_fixture", "firefly", "main.go"), filepath.Join(gopath, "src", target), "main.go")
		Expect(os.Setenv("GOPATH", filepath.Join(os.TempDir(), "emptyFakeGopath"))).To(Succeed())
		Expect(os.Environ()).To(ContainElement(fmt.Sprintf("GOPATH=%s", filepath.Join(os.TempDir(), "emptyFakeGopath"))))
	})

	AfterEach(func() {
		if original == "" {
			Expect(os.Unsetenv("GOPATH")).To(Succeed())
		} else {
			Expect(os.Setenv("GOPATH", original)).To(Succeed())
		}
		if gopath != "" {
			os.RemoveAll(gopath)
		}
	})

	It("appends the gopath env var", func() {
		compiledPath, err := gexec.BuildIn(gopath, target)
		Expect(err).NotTo(HaveOccurred())
		Expect(compiledPath).Should(BeAnExistingFile())
	})

	It("resets GOPATH to its original value", func() {
		_, err := gexec.BuildIn(gopath, target)
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Getenv("GOPATH")).To(Equal(filepath.Join(os.TempDir(), "emptyFakeGopath")))
	})
})

var _ = Describe(".CompileTest", func() {
	Context("a remote package", func() {
		const remotePackage = "github.com/onsi/ginkgo/types"

		It("compiles the specified test package", func() {
			compiledPath, err := gexec.GetAndCompileTest(remotePackage)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})
	})

	When("there have been previous calls to CompileTest", func() {
		BeforeEach(func() {
			_, err := gexec.CompileTest(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("compiles the specified test package", func() {
			compiledPath, err := gexec.CompileTest(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})

		Context("and CleanupBuildArtifacts has been called", func() {
			BeforeEach(func() {
				gexec.CleanupBuildArtifacts()
			})

			It("compiles the specified test package", func() {
				fireflyTestPath, err := gexec.CompileTest(packagePath)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fireflyTestPath).Should(BeAnExistingFile())
			})
		})
	})

	When("there have been previous calls to Build", func() {
		BeforeEach(func() {
			_, err := gexec.Build(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("compiles the specified test package", func() {
			compiledPath, err := gexec.CompileTest(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})

		Context("and CleanupBuildArtifacts has been called", func() {
			BeforeEach(func() {
				gexec.CleanupBuildArtifacts()
			})

			It("compiles the specified test package", func() {
				fireflyTestPath, err := gexec.CompileTest(packagePath)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fireflyTestPath).Should(BeAnExistingFile())
			})
		})
	})
})

var _ = Describe(".CompileTestWithEnvironment", func() {
	var err error
	env := []string{
		"GOOS=linux",
		"GOARCH=amd64",
	}

	Context("a remote package", func() {
		const remotePackage = "github.com/onsi/ginkgo/types"

		It("compiles the specified test package with the specified env vars", func() {
			compiledPath, err := gexec.GetAndCompileTestWithEnvironment(remotePackage, env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})
	})

	It("compiles the specified test package with the specified env vars", func() {
		compiledPath, err := gexec.CompileTestWithEnvironment(packagePath, env)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(compiledPath).Should(BeAnExistingFile())
	})

	It("returns the environment to a good state", func() {
		_, err = gexec.CompileTestWithEnvironment(packagePath, env)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(os.Environ()).ShouldNot(ContainElement("GOOS=linux"))
	})
})

var _ = Describe(".CompiledTestIn", func() {
	const (
		target = "github.com/onsi/gomega/gexec/_fixture/firefly"
	)

	var (
		original string
		gopath   string
	)

	BeforeEach(func() {
		var err error
		original = os.Getenv("GOPATH")
		gopath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		copyFile(filepath.Join("_fixture", "firefly", "main.go"), filepath.Join(gopath, "src", target), "main.go")
		Expect(os.Setenv("GOPATH", filepath.Join(os.TempDir(), "emptyFakeGopath"))).To(Succeed())
		Expect(os.Environ()).To(ContainElement(fmt.Sprintf("GOPATH=%s", filepath.Join(os.TempDir(), "emptyFakeGopath"))))
	})

	AfterEach(func() {
		if original == "" {
			Expect(os.Unsetenv("GOPATH")).To(Succeed())
		} else {
			Expect(os.Setenv("GOPATH", original)).To(Succeed())
		}
		if gopath != "" {
			os.RemoveAll(gopath)
		}
	})

	Context("a remote package", func() {
		const remotePackage = "github.com/onsi/ginkgo/types"

		It("compiles the specified test package", func() {
			compiledPath, err := gexec.GetAndCompileTestIn(gopath, remotePackage)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})
	})

	It("appends the gopath env var", func() {
		compiledPath, err := gexec.CompileTestIn(gopath, target)
		Expect(err).NotTo(HaveOccurred())
		Expect(compiledPath).Should(BeAnExistingFile())
	})

	It("resets GOPATH to its original value", func() {
		_, err := gexec.CompileTestIn(gopath, target)
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Getenv("GOPATH")).To(Equal(filepath.Join(os.TempDir(), "emptyFakeGopath")))
	})
})

func copyFile(source, directory, basename string) {
	Expect(os.MkdirAll(directory, 0755)).To(Succeed())
	content, err := ioutil.ReadFile(source)
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(filepath.Join(directory, basename), content, 0644)).To(Succeed())
}
