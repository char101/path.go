package path

import (
	"os"
)

// ChdirErr changes the working directory to this path
func (p Path) ChdirErr() (Path, error) {
	return p, WrapError(os.Chdir(string(p)))
}

// Chdir changes the working directory to this path
func (p Path) Chdir() Path {
	p, err := p.ChdirErr()
	HandleError(err)
	return p
}

// MkdirErr creates the directory
func (p Path) MkdirErr() (Path, error) {
	return p, os.Mkdir(string(p), DirPerm)
}

// Mkdir creates the directory
func (p Path) Mkdir() Path {
	p, err := p.MkdirErr()
	HandleError(err)
	return p
}

// MkdirAllErr creates the directories as necessary
func (p Path) MkdirAllErr() (Path, error) {
	return p, WrapError(os.MkdirAll(string(p), DirPerm))
}

// MkdirAll creates the directory as necessary
func (p Path) MkdirAll() Path {
	p, err := p.MkdirAllErr()
	HandleError(err)
	return p
}

// ChmodErr changes the file or directory permission
func (p Path) ChmodErr(perm os.FileMode) (Path, error) {
	return p, WrapError(os.Chmod(string(p), perm))
}

// Chmod changes the file or directory permission
func (p Path) Chmod(perm os.FileMode) Path {
	p, err := p.ChmodErr(perm)
	HandleError(err)
	return p
}

// SymlinkErr creates new symlink
func (p Path) SymlinkErr(dst interface{}) (Path, error) {
	dp := New(dst)
	if dp.IsDir() {
		dp = dp.Join(p.Basename())
	}
	return p, WrapError(os.Symlink(string(p), string(dp)))
}

// Symlink creates a new symlink from dst to p
func (p Path) Symlink(dst interface{}) Path {
	p, err := p.SymlinkErr(dst)
	HandleError(err)
	return p
}
