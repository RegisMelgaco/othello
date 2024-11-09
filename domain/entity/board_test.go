package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewBoard(t *testing.T) {
	t.Parallel()

	t.Run("expect new board", func(t *testing.T) {
		t.Parallel()

		b := NewBoard()

		assert.Equal(t, [][]PlayerName{
			{"none", "none", "none", "none", "none", "none", "none", "none"},
			{"none", "none", "none", "none", "none", "none", "none", "none"},
			{"none", "none", "none", "none", "none", "none", "none", "none"},
			{"none", "none", "none", "none", "none", "none", "none", "none"},
			{"none", "none", "none", "none", "none", "none", "none", "none"},
			{"none", "none", "none", "none", "none", "none", "none", "none"},
			{"none", "none", "none", "none", "none", "none", "none", "none"},
			{"none", "none", "none", "none", "none", "none", "none", "none"},
		}, b.Grid)
	})
}
