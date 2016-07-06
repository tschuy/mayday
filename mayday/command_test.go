package mayday_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/coreos/mayday/mayday"
	"github.com/stretchr/testify/assert"
)

var workspace string

func TestMain(m *testing.M) {
	// test setup
	// current_dir lives in /tmp/go-build
	current_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	workspace = current_dir + "/mayday_test"
	os.MkdirAll(workspace+"/mayday_commands", os.ModePerm)

	retCode := m.Run()

	os.Exit(retCode)
}

func TestNonexistentCommand(t *testing.T) {
	cmd := mayday.Command{Args: []string{"nonexistent"}}
	err := cmd.Run(workspace)
	assert.Equal(t, err.Error(), `could not find "nonexistent" in PATH`, "A nonexistent command should fail")
}

func TestCommand(t *testing.T) {
	cmd := mayday.Command{Args: []string{"df"}}
	cmdres := cmd.Run(workspace)
	assert.Nil(t, cmdres)

	f, err := ioutil.ReadFile(workspace + "/mayday_commands/df")
	assert.Nil(t, err)
	assert.Contains(t, string(f), "Filesystem")
}

func TestCommandWithArgs(t *testing.T) {
	cmd := mayday.Command{Args: []string{"echo", "hello"}}
	cmdres := cmd.Run(workspace)
	assert.Nil(t, cmdres)

	f, err := ioutil.ReadFile(workspace + "/mayday_commands/echo_hello")
	assert.Nil(t, err)
	assert.Equal(t, string(f), "hello\n")
}

func TestCommandWithLink(t *testing.T) {
	cmd := mayday.Command{Args: []string{"df", "-h"}, Link: "files"}
	cmdres := cmd.Run(workspace)
	assert.Nil(t, cmdres)

	f, err := ioutil.ReadFile(workspace + "/files")
	assert.Nil(t, err)
	assert.Contains(t, string(f), "Filesystem")
}
