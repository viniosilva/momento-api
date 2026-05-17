package uid_test

import (
	"testing"

	"momento/pkg/uid"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should generate a valid UUID v4", func(t *testing.T) {
		got := uid.New()

		assert.NotEmpty(t, got)
		_, err := uuid.Parse(got)
		assert.NoError(t, err)
		assert.Equal(t, uuid.Version(4), uuid.MustParse(got).Version())
	})

	t.Run("should generate unique IDs on each call", func(t *testing.T) {
		a := uid.New()
		b := uid.New()

		assert.NotEqual(t, a, b)
	})
}
