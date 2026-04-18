package tools_test

import (
	"momento/pkg/tools"
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

func TestValueOrDefault(t *testing.T) {
	t.Run("should return value when pointer is non-nil", func(t *testing.T) {
		v := "hello"
		got := tools.ValueOrDefault(&v, "default")
		assert.Equal(t, "hello", got)
	})

	t.Run("should return default when pointer is nil", func(t *testing.T) {
		got := tools.ValueOrDefault(nil, "default")
		assert.Equal(t, "default", got)
	})

	t.Run("should return zero value when pointer is nil and default is zero value", func(t *testing.T) {
		got := tools.ValueOrDefault(nil, 0)
		assert.Equal(t, 0, got)
	})

	t.Run("should work with int type", func(t *testing.T) {
		v := 42
		got := tools.ValueOrDefault(&v, 0)
		assert.Equal(t, 42, got)
	})

	t.Run("should return default when pointer to zero value is nil", func(t *testing.T) {
		got := tools.ValueOrDefault(nil, true)
		assert.Equal(t, true, got)
	})

	t.Run("should return pointed value even when it equals zero value", func(t *testing.T) {
		v := ""
		got := tools.ValueOrDefault(&v, "default")
		assert.Equal(t, "", got)
	})
}
