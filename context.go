package path

import "os"

// With runs the function within the context of the Path
func (p Path) With(f func()) error {
	wd, err := os.Getwd()
	f()
	if err == nil {
		err = os.Chdir(wd)
	}
	return err
}
