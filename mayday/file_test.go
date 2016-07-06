package mayday_test

import (
	"io/ioutil"
	"testing"

	"github.com/coreos/mayday/mayday"
	"github.com/stretchr/testify/assert"
)

func TestNonexistentFile(t *testing.T) {
	file := mayday.File{Path: "/etc/nonexistent"}
	err := file.Collect(workspace)
	assert.Equal(t, err.Error(), `error opening source file: open /etc/nonexistent: no such file or directory`)
}

func TestFile(t *testing.T) {
	file := mayday.File{Path: "/proc/meminfo"}
	res := file.Collect(workspace)
	assert.Nil(t, res)

	f, err := ioutil.ReadFile(workspace + "/proc/meminfo")
	assert.Nil(t, err)
	assert.Contains(t, string(f), "MemTotal")
}

func TestFileWithLink(t *testing.T) {
	file := mayday.File{Path: "/etc/crontab", Link: "cronfig"}
	res := file.Collect(workspace)
	assert.Nil(t, res)

	f, err := ioutil.ReadFile(workspace + "/cronfig")
	assert.Nil(t, err)
	assert.Contains(t, string(f), "crontab")
}
