package path

import (
	"io/ioutil"
	"path/filepath"
)

// Files returns the list of files inside the Path
func (p Path) Files() ([]Path, error) {
	files := make([]Path, 0)

	items, err := ioutil.ReadDir(string(p))
	if err != nil {
		return nil, err
	}

	for _, info := range items {
		// TODO: should it handle symlinks to regular file too?
		if info.Mode().IsRegular() {
			file := info.Name()
			files = append(files, p.Join(file))
		}
	}

	return files, nil
}

// Dirs returns the list of directories inside the Path
func (p Path) Dirs() ([]Path, error) {
	dirs := make([]Path, 0)

	// TODO

	return dirs, nil
}

// WalkFiles runs the function through each files in the Path and its subdirectories
func (p Path) WalkFiles(fn filepath.WalkFunc) {
	// TODO
}

// WalkDirs runs the function through each directories in the Path and its directories
func (p Path) WalkDirs(fn filepath.WalkFunc) {
	// TODO
}

// Glob returns the list of paths matching the pattern
func (p Path) Glob(pat string) []Path {
	items := make([]Path, 0)

	// TODO

	return items
}
