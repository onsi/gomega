package internal_test

import (
	"errors"
	"runtime"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Asynchronous Assertions", func() {
	var ig *InstrumentedGomega
	BeforeEach(func() {
		ig = NewInstrumentedGomega()
	})

	Describe("Basic Eventually support", func() {
		Context("the positive case", func() {
			It("polls the function and matcher until a match occurs", func() {
				counter := 0
				ig.G.Eventually(func() string {
					counter++
					if counter > 5 {
						return MATCH
					}
					return NO_MATCH
				}).Should(SpecMatch())
				Ω(counter).Should(Equal(6))
				Ω(ig.FailureMessage).Should(BeZero())
			})

			It("continues polling even if the matcher errors", func() {
				counter := 0
				ig.G.Eventually(func() string {
					counter++
					if counter > 5 {
						return MATCH
					}
					return ERR_MATCH
				}).Should(SpecMatch())
				Ω(counter).Should(Equal(6))
				Ω(ig.FailureMessage).Should(BeZero())
			})

			It("times out eventually if the assertion doesn't match in time", func() {
				counter := 0
				ig.G.Eventually(func() string {
					counter++
					if counter > 100 {
						return MATCH
					}
					return NO_MATCH
				}, "200ms", "20ms").Should(SpecMatch())
				Ω(counter).Should(BeNumerically(">", 2))
				Ω(counter).Should(BeNumerically("<", 20))
				Ω(ig.FailureMessage).Should(ContainSubstring("Timed out after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("positive: no match"))
				Ω(ig.FailureSkip).Should(Equal([]int{3}))
			})
		})

		Context("the negative case", func() {
			It("polls the function and matcher until a match does not occur", func() {
				counter := 0
				ig.G.Eventually(func() string {
					counter++
					if counter > 5 {
						return NO_MATCH
					}
					return MATCH
				}).ShouldNot(SpecMatch())
				Ω(counter).Should(Equal(6))
				Ω(ig.FailureMessage).Should(BeZero())
			})

			It("continues polling when the matcher errors - an error does not count as a successful non-match", func() {
				counter := 0
				ig.G.Eventually(func() string {
					counter++
					if counter > 5 {
						return NO_MATCH
					}
					return ERR_MATCH
				}).ShouldNot(SpecMatch())
				Ω(counter).Should(Equal(6))
				Ω(ig.FailureMessage).Should(BeZero())
			})

			It("times out eventually if the assertion doesn't match in time", func() {
				counter := 0
				ig.G.Eventually(func() string {
					counter++
					if counter > 100 {
						return NO_MATCH
					}
					return MATCH
				}, "200ms", "20ms").ShouldNot(SpecMatch())
				Ω(counter).Should(BeNumerically(">", 2))
				Ω(counter).Should(BeNumerically("<", 20))
				Ω(ig.FailureMessage).Should(ContainSubstring("Timed out after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("negative: match"))
				Ω(ig.FailureSkip).Should(Equal([]int{3}))
			})
		})

		Context("when a failure occurs", func() {
			It("registers the appropriate helper functions", func() {
				ig.G.Eventually(NO_MATCH, "50ms", "10ms").Should(SpecMatch())
				Ω(ig.FailureMessage).Should(ContainSubstring("Timed out after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("positive: no match"))
				Ω(ig.FailureSkip).Should(Equal([]int{3}))
				Ω(ig.RegisteredHelpers).Should(ContainElement("(*AsyncAssertion).Should"))
				Ω(ig.RegisteredHelpers).Should(ContainElement("(*AsyncAssertion).match"))
			})

			It("renders the matcher's error if an error occured", func() {
				ig.G.Eventually(ERR_MATCH, "50ms", "10ms").Should(SpecMatch())
				Ω(ig.FailureMessage).Should(ContainSubstring("Timed out after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("Error: spec matcher error"))
			})

			It("renders the optional description", func() {
				ig.G.Eventually(NO_MATCH, "50ms", "10ms").Should(SpecMatch(), "boop")
				Ω(ig.FailureMessage).Should(ContainSubstring("boop"))
			})

			It("formats and renders the optional description when there are multiple arguments", func() {
				ig.G.Eventually(NO_MATCH, "50ms", "10ms").Should(SpecMatch(), "boop %d", 17)
				Ω(ig.FailureMessage).Should(ContainSubstring("boop 17"))
			})

			It("calls the optional description if it is a function", func() {
				ig.G.Eventually(NO_MATCH, "50ms", "10ms").Should(SpecMatch(), func() string { return "boop" })
				Ω(ig.FailureMessage).Should(ContainSubstring("boop"))
			})
		})
	})

	Describe("Basic Consistently support", func() {
		Context("the positive case", func() {
			It("polls the function and matcher ensuring a match occurs consistently", func() {
				counter := 0
				ig.G.Consistently(func() string {
					counter++
					return MATCH
				}, "50ms", "10ms").Should(SpecMatch())
				Ω(counter).Should(BeNumerically(">", 1))
				Ω(counter).Should(BeNumerically("<", 7))
				Ω(ig.FailureMessage).Should(BeZero())
			})

			It("fails if the matcher ever errors", func() {
				counter := 0
				ig.G.Consistently(func() string {
					counter++
					if counter == 3 {
						return ERR_MATCH
					}
					return MATCH
				}, "50ms", "10ms").Should(SpecMatch())
				Ω(counter).Should(Equal(3))
				Ω(ig.FailureMessage).Should(ContainSubstring("Failed after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("Error: spec matcher error"))
			})

			It("fails if the matcher doesn't match at any point", func() {
				counter := 0
				ig.G.Consistently(func() string {
					counter++
					if counter == 3 {
						return NO_MATCH
					}
					return MATCH
				}, "50ms", "10ms").Should(SpecMatch())
				Ω(counter).Should(Equal(3))
				Ω(ig.FailureMessage).Should(ContainSubstring("Failed after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("positive: no match"))
			})
		})

		Context("the negative case", func() {
			It("polls the function and matcher ensuring a match never occurs", func() {
				counter := 0
				ig.G.Consistently(func() string {
					counter++
					return NO_MATCH
				}, "50ms", "10ms").ShouldNot(SpecMatch())
				Ω(counter).Should(BeNumerically(">", 1))
				Ω(counter).Should(BeNumerically("<", 7))
				Ω(ig.FailureMessage).Should(BeZero())
			})

			It("fails if the matcher ever errors", func() {
				counter := 0
				ig.G.Consistently(func() string {
					counter++
					if counter == 3 {
						return ERR_MATCH
					}
					return NO_MATCH
				}, "50ms", "10ms").ShouldNot(SpecMatch())
				Ω(counter).Should(Equal(3))
				Ω(ig.FailureMessage).Should(ContainSubstring("Failed after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("Error: spec matcher error"))
			})

			It("fails if the matcher matches at any point", func() {
				counter := 0
				ig.G.Consistently(func() string {
					counter++
					if counter == 3 {
						return MATCH
					}
					return NO_MATCH
				}, "50ms", "10ms").ShouldNot(SpecMatch())
				Ω(counter).Should(Equal(3))
				Ω(ig.FailureMessage).Should(ContainSubstring("Failed after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("negative: match"))
			})
		})

		Context("when a failure occurs", func() {
			It("registers the appropriate helper functions", func() {
				ig.G.Consistently(NO_MATCH).Should(SpecMatch())
				Ω(ig.FailureMessage).Should(ContainSubstring("Failed after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("positive: no match"))
				Ω(ig.FailureSkip).Should(Equal([]int{3}))
				Ω(ig.RegisteredHelpers).Should(ContainElement("(*AsyncAssertion).Should"))
				Ω(ig.RegisteredHelpers).Should(ContainElement("(*AsyncAssertion).match"))
			})

			It("renders the matcher's error if an error occured", func() {
				ig.G.Consistently(ERR_MATCH).Should(SpecMatch())
				Ω(ig.FailureMessage).Should(ContainSubstring("Failed after"))
				Ω(ig.FailureMessage).Should(ContainSubstring("Error: spec matcher error"))
			})

			It("renders the optional description", func() {
				ig.G.Consistently(NO_MATCH).Should(SpecMatch(), "boop")
				Ω(ig.FailureMessage).Should(ContainSubstring("boop"))
			})

			It("formats and renders the optional description when there are multiple arguments", func() {
				ig.G.Consistently(NO_MATCH).Should(SpecMatch(), "boop %d", 17)
				Ω(ig.FailureMessage).Should(ContainSubstring("boop 17"))
			})

			It("calls the optional description if it is a function", func() {
				ig.G.Consistently(NO_MATCH).Should(SpecMatch(), func() string { return "boop" })
				Ω(ig.FailureMessage).Should(ContainSubstring("boop"))
			})
		})
	})

	Describe("the passed-in actual", func() {
		type Foo struct{ Bar string }

		Context("when passed a value", func() {
			It("(eventually) continuously checks on the value until a match occurs", func() {
				c := make(chan bool)
				go func() {
					time.Sleep(100 * time.Millisecond)
					close(c)
				}()
				ig.G.Eventually(c, "1s", "10ms").Should(BeClosed())
				Ω(ig.FailureMessage).Should(BeZero())
			})

			It("(consistently) continuously checks on the value ensuring a match always occurs", func() {
				c := make(chan bool)
				close(c)
				ig.G.Consistently(c, "50ms", "10ms").Should(BeClosed())
				Ω(ig.FailureMessage).Should(BeZero())
			})
		})

		Context("when passed a function that takes no arguments and returns one value", func() {
			It("(eventually) polls the function until the returned value satisfies the matcher", func() {
				counter := 0
				ig.G.Eventually(func() int {
					counter += 1
					return counter
				}, "1s", "10ms").Should(BeNumerically(">", 5))
				Ω(ig.FailureMessage).Should(BeZero())
			})

			It("(consistently) polls the function ensuring the returned value satisfies the matcher", func() {
				counter := 0
				ig.G.Consistently(func() int {
					counter += 1
					return counter
				}, "50ms", "10ms").Should(BeNumerically("<", 20))
				Ω(counter).Should(BeNumerically(">", 2))
				Ω(ig.FailureMessage).Should(BeZero())
			})

			It("works when the function returns nil", func() {
				counter := 0
				ig.G.Eventually(func() error {
					counter += 1
					if counter > 5 {
						return nil
					}
					return errors.New("oops")
				}, "1s", "10ms").Should(BeNil())
				Ω(ig.FailureMessage).Should(BeZero())
			})
		})

		Context("when passed a function that takes no arguments and returns mutliple values", func() {
			Context("with Eventually", func() {
				It("polls the function until the first returned value satisfies the matcher _and_ all additional values are zero", func() {
					counter, s, f, err := 0, "hi", Foo{Bar: "hi"}, errors.New("hi")
					ig.G.Eventually(func() (int, string, Foo, error) {
						switch counter += 1; counter {
						case 2:
							s = ""
						case 3:
							f = Foo{}
						case 4:
							err = nil
						}
						return counter, s, f, err
					}, "1s", "10ms").Should(BeNumerically("<", 100))
					Ω(ig.FailureMessage).Should(BeZero())
					Ω(counter).Should(Equal(4))
				})

				It("reports on the non-zero value if it times out", func() {
					ig.G.Eventually(func() (int, string, Foo, error) {
						return 1, "", Foo{Bar: "hi"}, nil
					}, "30ms", "10ms").Should(BeNumerically("<", 100))
					Ω(ig.FailureMessage).Should(ContainSubstring("Error: Unexpected non-nil/non-zero extra argument at index 2:"))
					Ω(ig.FailureMessage).Should(ContainSubstring(`Foo{Bar:"hi"}`))
				})

				Context("when making a ShouldNot assertion", func() {
					It("doesn't succeed until the matcher is (not) satisfied with the first returned value _and_ all additional values are zero", func() {
						counter, s, f, err := 0, "hi", Foo{Bar: "hi"}, errors.New("hi")
						ig.G.Eventually(func() (int, string, Foo, error) {
							switch counter += 1; counter {
							case 2:
								s = ""
							case 3:
								f = Foo{}
							case 4:
								err = nil
							}
							return counter, s, f, err
						}, "1s", "10ms").ShouldNot(BeNumerically("<", 0))
						Ω(ig.FailureMessage).Should(BeZero())
						Ω(counter).Should(Equal(4))
					})
				})
			})

			Context("with Consistently", func() {
				It("polls the function and succeeds if all the values are zero and the matcher is consistently satisfied", func() {
					var err error
					counter, s, f := 0, "", Foo{}
					ig.G.Consistently(func() (int, string, Foo, error) {
						counter += 1
						return counter, s, f, err
					}, "50ms", "10ms").Should(BeNumerically("<", 100))
					Ω(ig.FailureMessage).Should(BeZero())
					Ω(counter).Should(BeNumerically(">", 2))
				})

				It("polls the function and fails any of the values are non-zero", func() {
					var err error
					counter, s, f := 0, "", Foo{}
					ig.G.Consistently(func() (int, string, Foo, error) {
						counter += 1
						if counter == 3 {
							f = Foo{Bar: "welp"}
						}
						return counter, s, f, err
					}, "50ms", "10ms").Should(BeNumerically("<", 100))
					Ω(ig.FailureMessage).Should(ContainSubstring("Error: Unexpected non-nil/non-zero extra argument at index 2:"))
					Ω(ig.FailureMessage).Should(ContainSubstring(`Foo{Bar:"welp"}`))
					Ω(counter).Should(Equal(3))
				})

				Context("when making a ShouldNot assertion", func() {
					It("succeeds if all additional values are zero", func() {
						var err error
						counter, s, f := 0, "", Foo{}
						ig.G.Consistently(func() (int, string, Foo, error) {
							counter += 1
							return counter, s, f, err
						}, "50ms", "10ms").ShouldNot(BeNumerically(">", 100))
						Ω(ig.FailureMessage).Should(BeZero())
						Ω(counter).Should(BeNumerically(">", 2))
					})

					It("fails if any additional values are ever non-zero", func() {
						var err error
						counter, s, f := 0, "", Foo{}
						ig.G.Consistently(func() (int, string, Foo, error) {
							counter += 1
							if counter == 3 {
								s = "welp"
							}
							return counter, s, f, err
						}, "50ms", "10ms").ShouldNot(BeNumerically(">", 100))
						Ω(ig.FailureMessage).Should(ContainSubstring("Error: Unexpected non-nil/non-zero extra argument at index 1:"))
						Ω(ig.FailureMessage).Should(ContainSubstring(`<string>: "welp"`))
						Ω(counter).Should(Equal(3))
					})
				})
			})
		})

		Context("when passed a function that takes a Gomega argument and returns values", func() {
			Context("with Eventually", func() {
				It("passes in a Gomega and passes if the matcher matches, all extra values are zero, and there are no failed assertions", func() {
					counter, s, f, err := 0, "hi", Foo{Bar: "hi"}, errors.New("hi")
					ig.G.Eventually(func(g Gomega) (int, string, Foo, error) {
						switch counter += 1; counter {
						case 2:
							s = ""
						case 3:
							f = Foo{}
						case 4:
							err = nil
						}
						if counter == 5 {
							g.Expect(true).To(BeTrue())
						} else {
							g.Expect(false).To(BeTrue())
							panic("boom") //never see since the expectation stops execution
						}
						return counter, s, f, err
					}, "1s", "10ms").Should(BeNumerically("<", 100))
					Ω(ig.FailureMessage).Should(BeZero())
					Ω(counter).Should(Equal(5))
				})

				It("times out if assertions in the function never succeed and reports on the error", func() {
					_, file, line, _ := runtime.Caller(0)
					ig.G.Eventually(func(g Gomega) int {
						g.Expect(false).To(BeTrue())
						return 10
					}, "30ms", "10ms").Should(Equal(10))
					Ω(ig.FailureMessage).Should(ContainSubstring("Error: Assertion in callback at %s:%d failed:", file, line+2))
					Ω(ig.FailureMessage).Should(ContainSubstring("Expected\n    <bool>: false\nto be true"))
				})

				It("forwards panics", func() {
					Ω(func() {
						ig.G.Eventually(func(g Gomega) int {
							g.Expect(true).To(BeTrue())
							panic("boom")
						}, "30ms", "10ms").Should(Equal(10))
					}).Should(PanicWith("boom"))
					Ω(ig.FailureMessage).Should(BeEmpty())
				})

				Context("when making a ShouldNot assertion", func() {
					It("doesn't succeed until all extra values are zero, there are no failed assertions, and the matcher is (not) satisfied", func() {
						counter, s, f, err := 0, "hi", Foo{Bar: "hi"}, errors.New("hi")
						ig.G.Eventually(func(g Gomega) (int, string, Foo, error) {
							switch counter += 1; counter {
							case 2:
								s = ""
							case 3:
								f = Foo{}
							case 4:
								err = nil
							}
							if counter == 5 {
								g.Expect(true).To(BeTrue())
							} else {
								g.Expect(false).To(BeTrue())
								panic("boom") //never see since the expectation stops execution
							}
							return counter, s, f, err
						}, "1s", "10ms").ShouldNot(BeNumerically("<", 0))
						Ω(ig.FailureMessage).Should(BeZero())
						Ω(counter).Should(Equal(5))
					})
				})

				It("fails if an assertion is never satisfied", func() {
					_, file, line, _ := runtime.Caller(0)
					ig.G.Eventually(func(g Gomega) int {
						g.Expect(false).To(BeTrue())
						return 9
					}, "30ms", "10ms").ShouldNot(Equal(10))
					Ω(ig.FailureMessage).Should(ContainSubstring("Error: Assertion in callback at %s:%d failed:", file, line+2))
					Ω(ig.FailureMessage).Should(ContainSubstring("Expected\n    <bool>: false\nto be true"))
				})
			})

			Context("with Consistently", func() {
				It("passes in a Gomega and passes if the matcher matches, all extra values are zero, and there are no failed assertions", func() {
					var err error
					counter, s, f := 0, "", Foo{}
					ig.G.Consistently(func(g Gomega) (int, string, Foo, error) {
						counter += 1
						g.Expect(true).To(BeTrue())
						return counter, s, f, err
					}, "50ms", "10ms").Should(BeNumerically("<", 100))
					Ω(ig.FailureMessage).Should(BeZero())
					Ω(counter).Should(BeNumerically(">", 2))
				})

				It("fails if the passed-in gomega ever hits a failure", func() {
					var err error
					counter, s, f := 0, "", Foo{}
					_, file, line, _ := runtime.Caller(0)
					ig.G.Consistently(func(g Gomega) (int, string, Foo, error) {
						counter += 1
						g.Expect(true).To(BeTrue())
						if counter == 3 {
							g.Expect(false).To(BeTrue())
							panic("boom") //never see this
						}
						return counter, s, f, err
					}, "50ms", "10ms").Should(BeNumerically("<", 100))
					Ω(ig.FailureMessage).Should(ContainSubstring("Error: Assertion in callback at %s:%d failed:", file, line+5))
					Ω(ig.FailureMessage).Should(ContainSubstring("Expected\n    <bool>: false\nto be true"))
					Ω(counter).Should(Equal(3))
				})

				It("forwards panics", func() {
					Ω(func() {
						ig.G.Consistently(func(g Gomega) int {
							g.Expect(true).To(BeTrue())
							panic("boom")
						}, "50ms", "10ms").Should(Equal(10))
					}).Should(PanicWith("boom"))
					Ω(ig.FailureMessage).Should(BeEmpty())
				})

				Context("when making a ShouldNot assertion", func() {
					It("succeeds if any interior assertions always pass", func() {
						ig.G.Consistently(func(g Gomega) int {
							g.Expect(true).To(BeTrue())
							return 9
						}, "50ms", "10ms").ShouldNot(Equal(10))
						Ω(ig.FailureMessage).Should(BeEmpty())
					})

					It("fails if any interior assertions ever fail", func() {
						counter := 0
						_, file, line, _ := runtime.Caller(0)
						ig.G.Consistently(func(g Gomega) int {
							g.Expect(true).To(BeTrue())
							counter += 1
							if counter == 3 {
								g.Expect(false).To(BeTrue())
								panic("boom") //never see this
							}
							return 9
						}, "50ms", "10ms").ShouldNot(Equal(10))
						Ω(ig.FailureMessage).Should(ContainSubstring("Error: Assertion in callback at %s:%d failed:", file, line+5))
						Ω(ig.FailureMessage).Should(ContainSubstring("Expected\n    <bool>: false\nto be true"))
					})
				})
			})
		})

		Context("when passed a function that takes a Gomega argument and returns nothing", func() {
			Context("with Eventually", func() {
				It("returns the first failed assertion as an error and so should Succeed() if the callback ever runs without issue", func() {
					counter := 0
					ig.G.Eventually(func(g Gomega) {
						counter += 1
						if counter < 5 {
							g.Expect(false).To(BeTrue())
							g.Expect("bloop").To(Equal("blarp"))
						}
					}, "1s", "10ms").Should(Succeed())
					Ω(counter).Should(Equal(5))
					Ω(ig.FailureMessage).Should(BeZero())
				})

				It("returns the first failed assertion as an error and so should timeout if the callback always fails", func() {
					counter := 0
					ig.G.Eventually(func(g Gomega) {
						counter += 1
						if counter < 5000 {
							g.Expect(false).To(BeTrue())
							g.Expect("bloop").To(Equal("blarp"))
						}
					}, "100ms", "10ms").Should(Succeed())
					Ω(counter).Should(BeNumerically(">", 1))
					Ω(ig.FailureMessage).Should(ContainSubstring("Expected success, but got an error"))
					Ω(ig.FailureMessage).Should(ContainSubstring("<bool>: false"))
					Ω(ig.FailureMessage).Should(ContainSubstring("to be true"))
					Ω(ig.FailureMessage).ShouldNot(ContainSubstring("bloop"))
				})

				It("returns the first failed assertion as an error and should satisy ShouldNot(Succeed) eventually", func() {
					counter := 0
					ig.G.Eventually(func(g Gomega) {
						counter += 1
						if counter > 5 {
							g.Expect(false).To(BeTrue())
							g.Expect("bloop").To(Equal("blarp"))
						}
					}, "100ms", "10ms").ShouldNot(Succeed())
					Ω(counter).Should(Equal(6))
					Ω(ig.FailureMessage).Should(BeZero())
				})

				It("should fail to ShouldNot(Succeed) eventually if an error never occurs", func() {
					ig.G.Eventually(func(g Gomega) {
						g.Expect(true).To(BeTrue())
					}, "50ms", "10ms").ShouldNot(Succeed())
					Ω(ig.FailureMessage).Should(ContainSubstring("Timed out after"))
					Ω(ig.FailureMessage).Should(ContainSubstring("Expected failure, but got no error."))
				})
			})

			Context("with Consistently", func() {
				It("returns the first failed assertion as an error and so should Succeed() if the callback always runs without issue", func() {
					counter := 0
					ig.G.Consistently(func(g Gomega) {
						counter += 1
						g.Expect(true).To(BeTrue())
					}, "50ms", "10ms").Should(Succeed())
					Ω(counter).Should(BeNumerically(">", 2))
					Ω(ig.FailureMessage).Should(BeZero())
				})

				It("returns the first failed assertion as an error and so should fail if the callback ever fails", func() {
					counter := 0
					ig.G.Consistently(func(g Gomega) {
						counter += 1
						g.Expect(true).To(BeTrue())
						if counter == 3 {
							g.Expect(false).To(BeTrue())
							g.Expect("bloop").To(Equal("blarp"))
						}
					}, "50ms", "10ms").Should(Succeed())
					Ω(ig.FailureMessage).Should(ContainSubstring("Expected success, but got an error"))
					Ω(ig.FailureMessage).Should(ContainSubstring("<bool>: false"))
					Ω(ig.FailureMessage).Should(ContainSubstring("to be true"))
					Ω(ig.FailureMessage).ShouldNot(ContainSubstring("bloop"))
					Ω(counter).Should(Equal(3))
				})

				It("returns the first failed assertion as an error and should satisy ShouldNot(Succeed) consistently if an error always occur", func() {
					counter := 0
					ig.G.Consistently(func(g Gomega) {
						counter += 1
						g.Expect(true).To(BeFalse())
					}, "50ms", "10ms").ShouldNot(Succeed())
					Ω(counter).Should(BeNumerically(">", 2))
					Ω(ig.FailureMessage).Should(BeZero())
				})

				It("should fail to satisfy ShouldNot(Succeed) consistently if an error ever does not occur", func() {
					counter := 0
					ig.G.Consistently(func(g Gomega) {
						counter += 1
						if counter == 3 {
							g.Expect(true).To(BeTrue())
						} else {
							g.Expect(false).To(BeTrue())
						}
					}, "50ms", "10ms").ShouldNot(Succeed())
					Ω(ig.FailureMessage).Should(ContainSubstring("Failed after"))
					Ω(ig.FailureMessage).Should(ContainSubstring("Expected failure, but got no error."))
					Ω(counter).Should(Equal(3))
				})
			})
		})

		Describe("when passed an invalid function", func() {
			It("errors immediately", func() {
				ig.G.Eventually(func() {})
				Ω(ig.FailureMessage).Should(Equal("The function passed to Gomega's async assertions should either take no arguments and return values, or take a single Gomega interface that it can use to make assertions within the body of the function.  When taking a Gomega interface the function can optionally return values or return nothing.  The function you passed takes 0 arguments and returns 0 values."))
				Ω(ig.FailureSkip).Should(Equal([]int{4}))

				ig = NewInstrumentedGomega()
				ig.G.Eventually(func(g Gomega, foo string) {})
				Ω(ig.FailureMessage).Should(Equal("The function passed to Gomega's async assertions should either take no arguments and return values, or take a single Gomega interface that it can use to make assertions within the body of the function.  When taking a Gomega interface the function can optionally return values or return nothing.  The function you passed takes 2 arguments and returns 0 values."))
				Ω(ig.FailureSkip).Should(Equal([]int{4}))

				ig = NewInstrumentedGomega()
				ig.G.Eventually(func(foo string) {})
				Ω(ig.FailureMessage).Should(Equal("The function passed to Gomega's async assertions should either take no arguments and return values, or take a single Gomega interface that it can use to make assertions within the body of the function.  When taking a Gomega interface the function can optionally return values or return nothing.  The function you passed takes 1 arguments and returns 0 values."))
				Ω(ig.FailureSkip).Should(Equal([]int{4}))
			})
		})
	})

	Describe("when using OracleMatchers", func() {
		It("stops and gives up with an appropriate failure message if the OracleMatcher says things can't change", func() {
			c := make(chan bool)
			close(c)

			t := time.Now()
			ig.G.Eventually(c, "100ms", "10ms").Should(Receive(), "Receive is an OracleMatcher that gives up if the channel is closed")
			Ω(time.Since(t)).Should(BeNumerically("<", 90*time.Millisecond))
			Ω(ig.FailureMessage).Should(ContainSubstring("No future change is possible."))
			Ω(ig.FailureMessage).Should(ContainSubstring("The channel is closed."))
		})

		It("never gives up if actual is a function", func() {
			c := make(chan bool)
			close(c)

			t := time.Now()
			ig.G.Eventually(func() chan bool { return c }, "100ms", "10ms").Should(Receive(), "Receive is an OracleMatcher that gives up if the channel is closed")
			Ω(time.Since(t)).Should(BeNumerically(">=", 90*time.Millisecond))
			Ω(ig.FailureMessage).ShouldNot(ContainSubstring("No future change is possible."))
			Ω(ig.FailureMessage).Should(ContainSubstring("Timed out after"))
		})
	})
})
