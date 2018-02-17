package path

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

// Write writes the data into the file
func (p Path) Write(data interface{}) error {
	switch v := data.(type) {
	case string:
		return ioutil.WriteFile(string(p), []byte(v), FilePerm)
	case []byte:
		return ioutil.WriteFile(string(p), v, FilePerm)
	case io.Reader:
		fh, err := os.Create(string(p))
		if err != nil {
			return err
		}
		io.Copy(fh, v)
	default:
		return fmt.Errorf("Unsupported type: %T", v)
	}
	return nil
}

func (p Path) ReadAll() []byte {
	b, err := ioutil.ReadFile(string(p))
	HandleError(err)
	return b
}

// Open opens the file
func (p Path) Open() (*os.File, error) {
	return os.Open(string(p))
}

// Create creates the file
func (p Path) Create() (*os.File, error) {
	return os.Create(string(p))
}

// CopyErr copies the file
func (p Path) CopyErr(dst interface{}) (Path, error) {
	dp, err := NewErr(dst)
	if err != nil {
		return p, WrapError(err)
	}
	if dp.IsDir() {
		dp = dp.Join(p.Basename())
	}

	in, err := p.Open()
	if err != nil {
		return p, WrapError(err)
	}
	defer in.Close()

	out, err := dp.Create()
	if err != nil {
		return p, WrapError(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)

	return p, WrapError(err)
}

// Copy copies the path
func (p Path) Copy(dst interface{}) Path {
	p, err := p.CopyErr(dst)
	HandleError(err)
	return p
}

// RegexpReplace replaces the file content
func (p Path) RegexpReplace(pat string, repl string) {
	p.Write(regexp.MustCompile(pat).ReplaceAllString(string(p.ReadAll()), repl))
}
