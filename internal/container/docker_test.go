package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCommand(t *testing.T) {
	command := "ls -la \"/path with space\""
	expectedArguments := []string{"ls", "-la", "/path with space"}

	actualArguments, _ := splitCommand(command)

	assert.Equal(t, expectedArguments, actualArguments)
}
