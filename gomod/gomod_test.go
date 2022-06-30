package gomod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathToID(t *testing.T) {
	assert.Equal(t, "c.b.a.d.e", PathToID("a.b.c/d/e"))
}
