package defaults_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	d "github.com/onsi/gomega/internal/defaults"
)

var _ = Describe("Durations", func() {
	var (
		duration       *time.Duration
		envVarGot      string
		envVarToReturn string

		getDurationFromEnv = func(name string) string {
			envVarGot = name
			return envVarToReturn
		}

		setDuration = func(t time.Duration) {
			duration = &t
		}

		setDurationCalled = func() bool {
			return duration != nil
		}

		resetDuration = func() {
			duration = nil
		}
	)

	BeforeEach(func() {
		resetDuration()
	})

	Context("When the environment has a duration", func() {
		Context("When the duration is valid", func() {
			BeforeEach(func() {
				envVarToReturn = "10m"

				d.SetDurationFromEnv(getDurationFromEnv, setDuration, "MY_ENV_VAR")
			})

			It("sets the duration", func() {
				Expect(envVarGot).To(Equal("MY_ENV_VAR"))
				Expect(setDurationCalled()).To(Equal(true))
				Expect(*duration).To(Equal(10 * time.Minute))
			})
		})

		Context("When the duration is not valid", func() {
			BeforeEach(func() {
				envVarToReturn = "10"
			})

			It("panics with a helpful error message", func() {
				Expect(func() {
					d.SetDurationFromEnv(getDurationFromEnv, setDuration, "MY_ENV_VAR")
				}).To(PanicWith(MatchRegexp("Expected a duration when using MY_ENV_VAR")))
			})
		})
	})

	Context("When the environment does not have a duration", func() {
		BeforeEach(func() {
			envVarToReturn = ""

			d.SetDurationFromEnv(getDurationFromEnv, setDuration, "MY_ENV_VAR")
		})

		It("does not set the duration", func() {
			Expect(envVarGot).To(Equal("MY_ENV_VAR"))
			Expect(setDurationCalled()).To(Equal(false))
			Expect(duration).To(BeNil())
		})
	})
})
