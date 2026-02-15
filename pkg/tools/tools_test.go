package tools_test

import (
	"pinnado/pkg/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtoiOrDefault(t *testing.T) {
	t.Run("should return value when value is a number", func(t *testing.T) {
		got := tools.AtoiOrDefault("10", 1)
		assert.Equal(t, 10, got)
	})
	t.Run("should return default value when value is empty", func(t *testing.T) {
		got := tools.AtoiOrDefault("", 1)
		assert.Equal(t, 1, got)
	})

	t.Run("should return default value when value is not a number", func(t *testing.T) {
		got := tools.AtoiOrDefault("not a number", 1)
		assert.Equal(t, 1, got)
	})
}
