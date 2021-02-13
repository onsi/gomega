package gexec_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const currentPackage = "github.com/onsi/gomega"

var _ = Context("a local package", func() {
	suiteTest("./_fixture/firefly")

	Context("with remote url", func() {
		suiteTest("github.com/onsi/gomega/gexec/_fixture/firefly")
	})
})

var _ = Context("a remote package", func() {
	suiteTest("github.com/onsi/ginkgo/types")
})

func suiteTest(packagePath string) {
	Describe(".Get", func() {
		It("get the specified package", func() {
			_, err := gexec.Get(packagePath)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe(".GetTests", func() {
		It("get the specified package", func() {
			_, err := gexec.GetTests(packagePath)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe(".Build", func() {
		var goPackage gexec.GoPackage

		BeforeEach(func() {
			var err error
			goPackage, err = gexec.Get(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		When("there have been previous calls to CompileTest", func() {
			BeforeEach(func() {
				p, err := gexec.GetTests(packagePath)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = p.Build()
				Expect(err).NotTo(HaveOccurred())
			})

			It("compiles the specified package", func() {
				compiledPath, err := goPackage.Build()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(compiledPath).Should(BeAnExistingFile())
			})

			Context("and CleanupBuildArtifacts has been called", func() {
				BeforeEach(func() {
					gexec.CleanupBuildArtifacts()
				})

				It("compiles the specified package", func() {
					fireflyPath, err := goPackage.Build()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(fireflyPath).Should(BeAnExistingFile())
				})
			})
		})

		When("there have been previous calls to Build", func() {
			BeforeEach(func() {
				p, err := gexec.Get(packagePath)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = p.Build()
				Expect(err).NotTo(HaveOccurred())
			})

			It("compiles the specified package", func() {
				compiledPath, err := goPackage.Build()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(compiledPath).Should(BeAnExistingFile())
			})

			Context("and CleanupBuildArtifacts has been called", func() {
				BeforeEach(func() {
					gexec.CleanupBuildArtifacts()
				})

				It("compiles the specified package", func() {
					fireflyPath, err := goPackage.Build()
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

		var goPackage gexec.GoPackage

		BeforeEach(func() {
			var err error
			goPackage, err = gexec.Get(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("compiles the specified package with the specified env vars", func() {
			compiledPath, err := goPackage.BuildWithEnvironment(env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})

		It("returns the environment to a good state", func() {
			_, err = goPackage.BuildWithEnvironment(env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(os.Environ()).ShouldNot(ContainElement("GOOS=linux"))
		})
	})

	var _ = Describe(".BuildIn", func() {
		var (
			original  string
			gopath    string
			goPackage gexec.GoPackage
		)

		BeforeEach(func() {
			var err error
			original = os.Getenv("GOPATH")
			gopath, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			wd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			destination := filepath.Join(gopath, "src", currentPackage)
			copy(path.Join(wd, ".."), destination)

			Expect(os.Setenv("GOPATH", filepath.Join(os.TempDir(), "emptyFakeGopath"))).To(Succeed())
			Expect(os.Environ()).To(ContainElement(fmt.Sprintf("GOPATH=%s", filepath.Join(os.TempDir(), "emptyFakeGopath"))))

			goPackage, err = gexec.Get(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
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
			compiledPath, err := goPackage.BuildIn(gopath)
			Expect(err).NotTo(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})

		It("resets GOPATH to its original value", func() {
			_, err := goPackage.BuildIn(gopath)
			Expect(err).NotTo(HaveOccurred())
			Expect(os.Getenv("GOPATH")).To(Equal(filepath.Join(os.TempDir(), "emptyFakeGopath")))
		})
	})

	var _ = Describe(".CompileTest", func() {
		var goPackage gexec.GoPackage

		BeforeEach(func() {
			var err error
			goPackage, err = gexec.GetTests(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		When("there have been previous calls to Build", func() {
			BeforeEach(func() {
				_, err := goPackage.Build()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("compiles the specified test package", func() {
				compiledPath, err := goPackage.Build()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(compiledPath).Should(BeAnExistingFile())
			})

			Context("and CleanupBuildArtifacts has been called", func() {
				BeforeEach(func() {
					gexec.CleanupBuildArtifacts()
				})

				It("compiles the specified test package", func() {
					fireflyTestPath, err := goPackage.Build()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(fireflyTestPath).Should(BeAnExistingFile())
				})
			})
		})

		When("there have been previous calls to Build", func() {
			BeforeEach(func() {
				p, err := gexec.Get(packagePath)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = p.Build()
				Expect(err).NotTo(HaveOccurred())
			})

			It("compiles the specified test package", func() {
				compiledPath, err := goPackage.Build()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(compiledPath).Should(BeAnExistingFile())
			})

			Context("and CleanupBuildArtifacts has been called", func() {
				BeforeEach(func() {
					gexec.CleanupBuildArtifacts()
				})

				It("compiles the specified test package", func() {
					fireflyTestPath, err := goPackage.Build()
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

		var goPackage gexec.GoPackage

		BeforeEach(func() {
			var err error
			goPackage, err = gexec.GetTests(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("compiles the specified test package with the specified env vars", func() {
			compiledPath, err := goPackage.BuildWithEnvironment(env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})

		It("returns the environment to a good state", func() {
			_, err = goPackage.BuildWithEnvironment(env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(os.Environ()).ShouldNot(ContainElement("GOOS=linux"))
		})
	})

	Describe(".CompiledTestIn", func() {
		var (
			original  string
			gopath    string
			goPackage gexec.GoPackage
		)

		BeforeEach(func() {
			var err error
			original = os.Getenv("GOPATH")
			gopath, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			wd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			destination := filepath.Join(gopath, "src", currentPackage)
			copy(path.Join(wd, ".."), destination)

			Expect(os.Setenv("GOPATH", filepath.Join(os.TempDir(), "emptyFakeGopath"))).To(Succeed())
			Expect(os.Environ()).To(ContainElement(fmt.Sprintf("GOPATH=%s", filepath.Join(os.TempDir(), "emptyFakeGopath"))))

			goPackage, err = gexec.Get(packagePath)
			Expect(err).ShouldNot(HaveOccurred())
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
			compiledPath, err := goPackage.BuildIn(gopath)
			Expect(err).NotTo(HaveOccurred())
			Expect(compiledPath).Should(BeAnExistingFile())
		})

		It("resets GOPATH to its original value", func() {
			_, err := goPackage.BuildIn(gopath)
			Expect(err).NotTo(HaveOccurred())
			Expect(os.Getenv("GOPATH")).To(Equal(filepath.Join(os.TempDir(), "emptyFakeGopath")))
		})
	})
}

func copy(source, destination string) {
	Expect(os.MkdirAll(destination, 0755)).To(Succeed())

	err := filepath.Walk(source, func(filePath string, info os.FileInfo, err error) error {
		relPath := strings.Replace(filePath, source, "", 1)
		if relPath == "" {
			return nil
		}

		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), 0755)
		} else {
			data, err := ioutil.ReadFile(filepath.Join(source, relPath))
			if err != nil {
				return err
			}

			return ioutil.WriteFile(filepath.Join(destination, relPath), data, info.Mode())
		}
	})
	Expect(err).NotTo(HaveOccurred())
}
