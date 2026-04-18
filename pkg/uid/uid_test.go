package uid_test

import (
	"testing"

	"momento/pkg/uid"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should generate a valid ID", func(t *testing.T) {
		got := uid.New()

		assert.NotEmpty(t, got)
		assert.Len(t, got, 24)
		assert.Regexp(t, `^[0-9a-f]{24}$`, got)
	})

	t.Run("should generate unique IDs on each call", func(t *testing.T) {
		a := uid.New()
		b := uid.New()

		assert.NotEqual(t, a, b)
	})
}
