package path

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"
)

// Path contains the path
type Path string

var (
	// EmptyPath is the path returned with errors
	EmptyPath = New("")

	// FilePerm is the default file permission
	FilePerm os.FileMode = 0600

	// DirPerm is the default directory permission
	DirPerm os.FileMode = 0700

	// HandleError is the default error handler for Must* functions
	HandleError = func(err error) {
		if err != nil {
			if e, ok := err.(*errors.Error); ok {
				fmt.Println(e.ErrorStack())
			}
			os.Exit(1)
		}
	}

	// WrapError wraps the error data with additional data
	WrapError = func(err error) error {
		if err == nil {
			return nil
		}
		return errors.Wrap(err, 1)
	}
)

// NewErr creates a new Path
func NewErr(p interface{}) (Path, error) {
	switch v := p.(type) {
	case Path:
		return v, nil
	case string:
		if v == "" {
			var err error
			v, err = os.Getwd()
			if err != nil {
				return Path(""), WrapError(err)
			}
		}
		return Path(v), nil
	case fmt.Stringer:
		return Path(v.String()), nil
	default:
		return Path(""), WrapError(fmt.Errorf("Unsupported type: %T", v))
	}
}

// New creates a new Path
func New(p interface{}) Path {
	pp, err := NewErr(p)
	HandleError(err)
	return pp
}

// String returns the path as a string
func (p Path) String() string {
	return string(p)
}

// Empty returns true if the path is empty
func (p Path) Empty() bool {
	return len(string(p)) == 0
}

// Ext returns the extension part of the Path
func (p Path) Ext() string {
	return filepath.Ext(string(p))
}

// Basename returns the basename part of the Path
func (p Path) Basename() Path {
	return Path(filepath.Base(string(p)))
}

// Name is an alias for Basename
func (p Path) Name() Path {
	return p.Basename()
}

// Namebase returns the basename of the Path sans the extension
func (p Path) Namebase() string {
	b := p.Basename()
	e := b.Ext()
	if e != "" {
		return string(p)[:len(b)-len(e)-1]
	}
	return string(b)
}

// Dir returns the directory part of the path
func (p Path) Dir() Path {
	return Path(filepath.Dir(string(p)))
}

// AbsErr returns an absolute representation of the P
func (p Path) AbsErr() (Path, error) {
	a, err := filepath.Abs(string(p))
	if err != nil {
		return p, WrapError(err)
	}
	return Path(a), nil
}

// Abs returns an absolute representation of the P
func (p Path) Abs() Path {
	a, err := p.AbsErr()
	HandleError(err)
	return Path(a)
}

// Parts splits a path into its components
func (p Path) Parts() []string {
	if runtime.GOOS == "windows" {
		return regexp.MustCompile(`/\`).Split(string(p), -1)
	}
	return strings.Split(string(p), "/")
}

// Volume returns the volume of the path
func (p Path) Volume() string {
	return filepath.VolumeName(string(p))
}

// SplitVolume splits a path into its volume and subsequent part
func (p Path) SplitVolume() (string, string) {
	v := p.Volume()
	return v, string(p)[len(v):]
}

// RelToErr ...
func (p Path) RelToErr(base interface{}) (Path, error) {
	bp, err := NewErr(base)
	if err != nil {
		return EmptyPath, err
	}
	r, err := filepath.Rel(string(bp), string(p))
	if err != nil {
		return EmptyPath, WrapError(err)
	}
	return New(r), nil
}

// RelTo ...
func (p Path) RelTo(base interface{}) Path {
	r, err := p.RelToErr(base)
	HandleError(err)
	return r
}

// RelOfErr ...
func (p Path) RelOfErr(target interface{}) (Path, error) {
	t, err := NewErr(target)
	if err != nil {
		return EmptyPath, err
	}
	return t.RelToErr(p)
}

// RelOf ...
func (p Path) RelOf(target interface{}) Path {
	r, err := p.RelOfErr(target)
	HandleError(err)
	return r
}

// StatErr returns the os.Stat() value for the filesystem object
func (p Path) StatErr() (os.FileInfo, error) {
	i, err := os.Stat(string(p))
	return i, WrapError(err)
}

// Stat returns the os.Stat() value for the filesystem object
func (p Path) Stat() os.FileInfo {
	st, err := p.StatErr()
	HandleError(err)
	return st
}

// SizeErr returns the file size
func (p Path) SizeErr() (int64, error) {
	st, err := p.StatErr()
	if err != nil {
		return 0, WrapError(err)
	}
	return st.Size(), nil
}

// Size return the file size
func (p Path) Size() int64 {
	size, err := p.SizeErr()
	HandleError(err)
	return size
}

// ModTimeErr returns the file/directory modification time
func (p Path) ModTimeErr() (time.Time, error) {
	st, err := os.Stat(string(p))
	if err != nil {
		return time.Time{}, WrapError(err)
	}
	return st.ModTime(), nil
}

// ModTime returns the file/directory modification time
func (p Path) ModTime() time.Time {
	mod, err := p.ModTimeErr()
	HandleError(err)
	return mod
}

// Join joins this path with another Path
func (p Path) Join(o interface{}) Path {
	switch v := o.(type) {
	case string:
		return Path(filepath.Join(string(p), v))
	case Path:
		return Path(filepath.Join(string(p), string(v)))
	case fmt.Stringer:
		return Path(filepath.Join(string(p), v.String()))
	case int:
		return Path(filepath.Join(string(p), strconv.Itoa(v)))
	default:
		log.Fatalf("Unhandled type: %T", v)
	}
	return p
}
