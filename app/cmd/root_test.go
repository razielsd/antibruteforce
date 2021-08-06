package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtractFirstArg_ArgExists(t *testing.T) {
	exp := "my value"
	v, err := extractFirstArg([]string{exp})
	require.NoError(t, err)
	require.Equal(t, exp, v)
}

func TestExtractFirstArg_ArgNotExists(t *testing.T) {
	_, err := extractFirstArg([]string{})
	require.Error(t, err)
}
