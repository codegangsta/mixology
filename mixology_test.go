package mixology_test

import (
	"testing"

	"github.com/codegangsta/mixology"
)

func TestBasicMix(t *testing.T) {
	m := mixology.New()

	equals(t, 1, 1)
}
