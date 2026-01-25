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
		var ctx = context.Background()
		actual := repo.CheckDBConn(ctx)
		assert.Nil(t, actual)
	})
	// act
	t.Run("without context", func(t *testing.T) {
		var ctx context.Context = nil
		actual := repo.CheckDBConn(ctx)
		assert.Nil(t, actual)
	})
}
