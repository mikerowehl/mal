package mal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenize(t *testing.T) {
	tokens := tokenize("(1 (2 3) 4)")
	require.Equal(t, len(tokens), 8)
}

func TestReader(t *testing.T) {
	val, err := Read_str("(1 (2 3) 4 5)")
	require.Nil(t, err)
	require.IsType(t, MalList{}, val)
	list := val.(MalList)
	require.Equal(t, len(list), 4)
}
