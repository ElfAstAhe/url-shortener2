package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemCheckDBConn_ShouldReturnNull(t *testing.T) {
	// prepare
	repo, err := NewDBConnCheckImMemRepo()
	require.NoError(t, err)
	// act
	t.Run("with context", func(t *testing.T) {
		actual := repo.CheckDBConn(context.Background())
		assert.Nil(t, actual)
	})
	// act
	t.Run("without context", func(t *testing.T) {
		actual := repo.CheckDBConn(nil)
		assert.Nil(t, actual)
	})
}
