package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapGroups(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	m := map[string]string{
		"my_map": "group_my",
	}

	type data struct{
		groups []string
		expected []string
	}

	items := []data{
		{groups: []string{"other", "my_group"}, expected: []string{"other", "my_group"}},
		{groups: []string{"my_map", "my_group"}, expected: []string{"group_my", "my_group"}},
	}

	for _, item := range items {
		assert.Equal(item.expected, mapSystemGroups(item.groups, m))
	}
}
