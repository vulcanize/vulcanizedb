package helpers

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mitchellh/go-homedir"
)

func CleanPath(str string) (string, error) {
	path, err := homedir.Expand(filepath.Clean(str))
	if err != nil {
		return "", err
	}
	if strings.Contains(path, "$GOPATH") {
		env := os.Getenv("GOPATH")
		spl := strings.Split(path, "$GOPATH")[1]
		path = filepath.Join(env, spl)
	}

	return path, nil
}

func ClearFiles(files ...string) error {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		} else if os.IsNotExist(err) {
			// fall through
		} else {
			return err
		}
	}

	return nil
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, syscall.O_CREAT|syscall.O_EXCL|os.O_WRONLY, os.FileMode(0666)) // Doesn't overwrite files
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
