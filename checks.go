package path

import (
	"os"

	"github.com/go-errors/errors"
)

// ExistsErr checks if the path exists
func (p Path) ExistsErr() (bool, error) {
	_, err := p.StatErr()
	if err != nil {
		uerr := err
		if e, ok := err.(*errors.Error); ok {
			uerr = e.Err
		}
		if os.IsNotExist(uerr) {
			return false, nil
		}
		return false, WrapError(err)
	}
	return true, nil
}

// Exists checks if the path exists
func (p Path) Exists() bool {
	ex, err := p.ExistsErr()
	HandleError(err)
	return ex
}

// IsFileErr checks if the path is a file
func (p Path) IsFileErr() (bool, error) {
	if !p.Exists() {
		return false, nil
	}
	st, err := p.StatErr()
	if err != nil {
		return false, WrapError(err)
	}
	return st.Mode().IsRegular(), nil
}

// IsFile checks if the path is a regular file
func (p Path) IsFile() bool {
	b, err := p.IsFileErr()
	HandleError(err)
	return b
}

// IsDirErr checks if the path is a directory
func (p Path) IsDirErr() (bool, error) {
	if !p.Exists() {
		return false, nil
	}
	st, err := p.StatErr()
	if err != nil {
		return false, WrapError(err)
	}
	return st.IsDir(), nil
}

// IsDir checks if the path is a directory
func (p Path) IsDir() bool {
	b, err := p.IsDirErr()
	HandleError(err)
	return b
}

// IsSymlinkErr checks if the path is a symlink
func (p Path) IsSymlinkErr() (bool, error) {
	if !p.Exists() {
		return false, nil
	}
	stat, err := p.StatErr()
	if err != nil {
		return false, WrapError(err)
	}
	return stat.Mode()&os.ModeSymlink > 0, nil
}

// IsSymlink checks if the path is a symlink
func (p Path) IsSymlink() bool {
	b, err := p.IsSymlinkErr()
	HandleError(err)
	return b
}
