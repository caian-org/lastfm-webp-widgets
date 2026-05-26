package lastfm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("empty API key returns error", func(t *testing.T) {
		c, err := New(LastFmClientOptions{Username: "someone", APIKey: ""})
		require.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("whitespace-only API key returns error", func(t *testing.T) {
		c, err := New(LastFmClientOptions{Username: "someone", APIKey: "   \t\n"})
		require.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("valid API key returns a client with username copied", func(t *testing.T) {
		c, err := New(LastFmClientOptions{Username: "rj", APIKey: "secret"})
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, "rj", c.Username)
	})
}
