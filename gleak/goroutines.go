package gleak

import "github.com/onsi/gomega/gleak/goroutine"

// Goroutines returns information about all goroutines: their goroutine IDs, the
// names of the topmost functions in the backtraces, and finally the goroutine
// backtraces.
func Goroutines() []goroutine.Goroutine {
	return goroutine.Goroutines()
}
