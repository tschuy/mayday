package mayday

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

// File encapsulates a file to be collected directly from the system
type File struct {
	Source io.ReadCloser // the source data to be copied. ex: os.Open(f.Path)
	Path   string        // full path to the file on the host system. Used for output name
	Link   string        // short name to link to the output (optional), e.g. "free"
}

// Collect copies the contents of File into a file of the same path in
// the given workspace
func (f *File) Collect(workspace string) error {
	var err error
	defer f.Source.Close()
	// Set everything up
	fn := path.Join(workspace, f.Path)
	dir := path.Dir(fn)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}
	dst, err := os.OpenFile(fn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error opening output file: %v", err)
	}

	// Actually copy the file
	log.Printf("Collecting file %q", f.Path)
	_, err = io.Copy(dst, f.Source)

	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	// If necessary, create a symlink
	return maybeCreateLink(f.Link, fn, workspace)
}
