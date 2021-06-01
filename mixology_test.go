package mixology_test

import (
	"testing"

	"github.com/codegangsta/mixology"
)

func ExampleMix(t *testing.T) {
	// create a new mixology
	m := mixology.New()

	equals(t, 1, 1)

	m.Run()
}
