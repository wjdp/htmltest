package output

import (
	"github.com/daviddengcn/go-assert"
	"testing"
)

func TestPluralise(t *testing.T) {
	assert.Equals(t, "singular", Pluralise(1, "1 apple", "2 apples"), "1 apple")
	assert.Equals(t, "plural", Pluralise(2, "1 apple", "2 apples"), "2 apples")
}
