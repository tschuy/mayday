package tar

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/coreos/mayday/mayday/tarable"
	"io"
	"log"
	"strings"
	"time"
)

type Tar struct {
	gzw    *gzip.Writer
	tw     *tar.Writer
	subdir string // subdirectory to put files in to prevent polluting current directory
}

func (t *Tar) Init(w io.Writer, subdir string) error {
	t.gzw = gzip.NewWriter(w)
	t.tw = tar.NewWriter(t.gzw)
	t.subdir = subdir
	return nil
}

func (t *Tar) Add(tb tarable.Tarable) error {
	var err error

	// virtual files, like those in /proc, report a size of 0 to stat().
	// this means the header in the tarfile reports a size of 0 for the file.
	// to avoid this, we copy the file into a buffer, and use that to get the
	// number of bytes to copy.

	buf := new(bytes.Buffer)
	buf.ReadFrom(tb.Content())

	header := tb.Header()
	header.Size = int64(buf.Len())
	header.Name = t.subdir + "/" + strings.TrimPrefix(header.Name, "/")

	if err = t.tw.WriteHeader(header); err != nil {
		log.Printf("error writing header: %s", err)
		return err
	}

	_, err = io.Copy(t.tw, buf)

	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}
	return nil
}

func (t *Tar) MaybeMakeLink(src string, dst string) error {
	if src == "" {
		return nil
	}

	var header tar.Header
	header.Name = t.subdir + "/" + src
	// relative path from location of link, already inside t.subdir
	header.Linkname = strings.TrimPrefix(dst, "/")
	header.Typeflag = tar.TypeSymlink
	header.ModTime = time.Now()

	log.Printf("Creating link: %q -> %q", src, dst)
	if err := t.tw.WriteHeader(&header); err != nil {
		return err
	}

	return nil
}

func (t *Tar) Close() error {
	t.tw.Flush()
	t.gzw.Flush()

	if err := t.tw.Close(); err != nil {
		log.Fatalf("error closing zipfile: %v", err)
	}
	if err := t.gzw.Close(); err != nil {
		log.Fatalf("error closing zipfile: %v", err)
	}
	return nil
}